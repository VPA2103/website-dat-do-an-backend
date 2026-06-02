# Flow tạo `menu_items` (Create/Import) + tạo `document`, `embedding`, `json`

Tài liệu này mô tả **chi tiết end-to-end** luồng tạo dữ liệu menu trong `backend.golang`, bao gồm:

- Tạo/cập nhật record `menu_items` (dữ liệu gốc: name/description/price/tags/allergens/ingredients)
- Sinh `document` (text) dùng cho RAG
- Sinh `embedding` (vector) từ `document`
- Sinh các field JSON:
  - JSONB “gốc”: `tags`, `allergens`, `ingredients`
  - JSONB “dẫn xuất”: `metadata`

> Ghi chú: code hiện tại đang ở “single restaurant mode”. `RestaurantID(c)` luôn trả về `"default"`.

---

## 0) Router / Endpoint tạo menu

Các endpoint liên quan menu nằm trong router:

- File: `internal/httpserver/router.go`
  - `POST /api/knowledge/menu/items` → `KnowledgeHandler.AddMenuItem` (tạo 1 món)
  - `POST /api/knowledge/menu/import` → `KnowledgeHandler.ImportMenu` (import Excel)
  - (tham khảo) `GET /api/knowledge/menu/items` / `DELETE /api/knowledge/menu/items/:item_id`

Restaurant ID được lấy từ Gin context:

- File: `internal/httpserver/handlers/tenant.go`
  - `func RestaurantID(c *gin.Context) string { return "default" }`

---

## 1) Data model ở DB: cột nào là “gốc”, cột nào là “dẫn xuất”

Schema được tạo trong:

- File: `internal/db/schema.go` (`db.EnsureSchema`)
- Bảng `menu_items` có 2 nhóm cột quan trọng:

### 1.1. Cột dữ liệu gốc (được tạo bởi FileStore)

- `id`, `restaurant_id`, `name`, `description`, `price`
- `tags JSONB`
- `allergens JSONB`
- `ingredients JSONB`

### 1.2. Cột dữ liệu dẫn xuất cho RAG (được tạo bởi pipeline embed + VectorStore)

- `embedding vector(<dim>)`
- `document TEXT`
- `metadata JSONB`

**Ý nghĩa:**
- `document`: text “mô tả chuẩn hoá” của món (để đưa vào context và để embed)
- `embedding`: vector của `document` (để vector search)
- `metadata`: JSONB phụ trợ, dùng kèm kết quả vector search (hiện tại chủ yếu để debug/hiển thị)

---

## 2) Flow: `POST /api/knowledge/menu/items` (tạo 1 món)

### 2.1. Handler: parse request + normalize

- File: `internal/httpserver/handlers/knowledge.go`
- Hàm: `func (h *KnowledgeHandler) AddMenuItem(c *gin.Context)`

Các bước:

1) Lấy `rid := RestaurantID(c)` (hiện tại luôn là `"default"`)
2) `ShouldBindJSON(&in)` vào struct `menuItemIn`
3) Normalize input JSON về `core.MenuItem`:
   - File: `internal/ingest/normalize_json_types.go`
   - Hàm: `NormalizeMenuItemFromJSONIn(...)`
   - Gọi vào `NormalizeMenuItemFromJSON(...)`:
     - File: `internal/ingest/normalize_json.go`
     - Chuẩn hoá list field (`tags`, `allergens`, `ingredients`) qua `normalizeAnyList(...)`:
       - File: `internal/ingest/helpers.go`

4) Validate tối thiểu: nếu `item.Name == ""` thì trả 400.

### 2.2. Persist dữ liệu gốc (FileStore)

- File: `internal/store/postgres_store.go`
- Hàm: `func (s *PostgresStore) UpsertMenuItem(restaurantID string, item core.MenuItem) (core.MenuItem, error)`

Các bước chính:

1) Validate/normalize `restaurantID` (`normalizeRestaurantID`) và trim `item.Name`
2) Nếu `item.ID` rỗng → generate UUID
3) Ensure restaurant tồn tại (best-effort, idempotent):
   - `s.ensureRestaurant(ctx, restaurantID)`
   - `ensureRestaurantSQL`: `INSERT INTO restaurants (id, updated_at) ... ON CONFLICT DO UPDATE ...`
4) Tạo JSONB cho các list “gốc”:
   - `tagsJSON := json.Marshal(item.Tags)`
   - `allergensJSON := json.Marshal(item.Allergens)`
   - `ingredientsJSON := json.Marshal(item.Ingredients)`

5) `INSERT INTO menu_items (...) VALUES (...) ON CONFLICT (id) DO UPDATE ...`

**Kết quả:** record `menu_items` đã có các cột “gốc” (name/description/price/tags/allergens/ingredients).

> Lưu ý: ở bước này **chưa tạo** `embedding`, `document`, `metadata`.

### 2.3. Sinh `document` + `embedding` + `metadata` (pipeline embed)

Sau khi lưu món xong, handler gọi:

- File: `internal/httpserver/handlers/knowledge.go`
  - `ingest.EmbedMenuItem(ctx, rid, created, h.vector, h.llm)`

Chi tiết embed:

- File: `internal/ingest/embed.go`
- Hàm: `func EmbedMenuItem(ctx, restaurantID string, item core.MenuItem, vector core.VectorStore, llm core.Gemini) error`

Các bước:

1) **Tạo `document`** từ `core.MenuItem`:
   - `doc := menuItemToDocument(item)`
   - File: `internal/ingest/documents.go`
   - Format document (ví dụ):
     - `Món: <name>. Mô tả: ... Giá: ... Tags: ... Dị ứng: ... Nguyên liệu: ...`

2) **Tạo `embedding`** từ `document`:
   - `emb, err := llm.Embed(ctx, doc)`
   - Nếu lỗi → return error
   - Nếu `len(emb) == 0` → return nil (coi như không có embedding)

3) **Tạo `metadata` (map)**:
   - `meta := map[string]any{"name": item.Name}`
   - Nếu có `Price` → `meta["price"] = *item.Price`
   - Nếu có `Tags` → `meta["tags"] = item.Tags`

4) Gọi VectorStore để lưu dữ liệu dẫn xuất:
   - `vector.UpsertMenuItem(ctx, restaurantID, item, emb, doc, meta)`

### 2.4. Persist dữ liệu dẫn xuất (VectorStore → UPDATE cột embedding/document/metadata)

- File: `internal/vector/pgvector/store.go`
- Hàm: `func (s *Store) UpsertMenuItem(ctx, restaurantID string, item core.MenuItem, embedding []float32, document string, metadata map[string]any) error`

Các bước:

1) Validate embedding dim (nếu cấu hình dim > 0):
   - `validateEmbedding(embedding)`
2) Marshal `metadata` → JSON bytes:
   - `b, _ := json.Marshal(metadata)`
3) Update các cột dẫn xuất:

```sql
UPDATE menu_items
SET embedding = $2, document = $3, metadata = $4, updated_at = now()
WHERE id = $1
```

**Kết quả:** record `menu_items` đã có `embedding`, `document`, `metadata`.

### 2.5. Error handling / tính “best-effort”

Trong handler `AddMenuItem`:

- Nếu persist gốc (`fs.UpsertMenuItem`) fail → request fail (500)
- Nếu embed fail (`EmbedMenuItem`) → chỉ `log.Printf(...)`, **không fail request**

Lý do: hệ thống ưu tiên “tạo món” thành công; embedding có thể chạy lại sau.

---

## 3) Flow: `POST /api/knowledge/menu/import` (import Excel nhiều món)

### 3.1. Handler: đọc file + gọi ingest

- File: `internal/httpserver/handlers/knowledge.go`
- Hàm: `func (h *KnowledgeHandler) ImportMenu(c *gin.Context)`

Các bước:

1) `rid := RestaurantID(c)`
2) Lấy file từ multipart form field `file`
3) Validate extension `.xlsx` (`ingest.HasXLSXExtension`)
4) `io.ReadAll(f)` → bytes
5) Gọi ingest:
   - `ingest.IngestMenuXLSX(ctx, rid, bytes, h.fs, h.vector, h.llm)`

### 3.2. Ingest: parse Excel → normalize từng row

- File: `internal/ingest/import.go`
- Hàm: `func IngestMenuXLSX(...) (ImportResult, error)`

Các bước:

1) `excelize.OpenReader(bytes.NewReader(xlsxBytes))`
2) Lấy sheet đầu tiên, đọc rows
3) Map header (tolerant, nhiều biến thể tiếng Việt/không dấu):
   - File: `internal/ingest/xlsx_shared.go`
   - Hàm: `headerMap(headers)` + `canonicalHeaders`

4) Với mỗi row:
   - Build `raw map[string]string` từ cell
   - Normalize sang `core.MenuItem`:
     - File: `internal/ingest/normalize_from_map.go`
     - Hàm: `NormalizeMenuItemFromMap(raw)`
   - Validate có `Name`

### 3.3. Persist gốc + embed từng món (best-effort)

Trong loop mỗi row:

1) `created, err := fs.UpsertMenuItem(restaurantID, item)`
   - Persist các cột gốc + JSONB gốc (`tags/allergens/ingredients`)
2) `_ = EmbedMenuItem(ctx, restaurantID, created, vector, llm)`
   - Best-effort: lỗi embed không làm fail import

`ImportResult` trả về:
- `Imported`: số món đã upsert
- `Errors`: lỗi validate row (ví dụ thiếu name)
- `Items`: list item đã tạo

---

## 4) Tóm tắt: ai tạo field nào?

### 4.1. `menu_items.tags/allergens/ingredients` (JSONB gốc)

- Tạo tại FileStore khi upsert:
  - `internal/store/postgres_store.go` → `PostgresStore.UpsertMenuItem`
  - Marshal từ `[]string` → JSON bytes

### 4.2. `menu_items.document`

- Tạo ở ingest layer:
  - `internal/ingest/documents.go` → `menuItemToDocument`

### 4.3. `menu_items.embedding`

- Tạo ở ingest layer:
  - `internal/ingest/embed.go` → `llm.Embed(ctx, doc)`

### 4.4. `menu_items.metadata` (JSONB dẫn xuất)

- Build map ở ingest layer:
  - `internal/ingest/embed.go` (`meta := map[string]any{...}`)
- Marshal + persist ở vector layer:
  - `internal/vector/pgvector/store.go` (`json.Marshal(metadata)` rồi `UPDATE menu_items ... metadata = $4`)

---

## 5) Ghi chú vận hành

- Embedding là “derived data”. Nếu thay đổi dimension, `db.ensureEmbeddingDim(...)` sẽ clear embedding cũ (set `NULL`) để tránh mismatch.
- Xoá món ăn:
  - `FileStore.DeleteMenuItem` xoá row khỏi `menu_items`.
  - `VectorStore.DeleteMenuItem` chỉ clear `embedding/document/metadata` (set `NULL`) — trong code hiện tại handler gọi cả hai theo best-effort.
