package ingest

import (
	"bytes"
	"context"
	"strconv"
	"strings"

	"github.com/vpa/quanlynhahang-backend/ai/core"
	"github.com/xuri/excelize/v2"
)

type ImportResult struct {
	Imported int             `json:"imported"`
	Errors   []string        `json:"errors"`
	Items    []core.MenuItem `json:"items"`
}

func IngestMenuXLSX(ctx context.Context, restaurantID string, xlsxBytes []byte, fs core.FileStore, vector core.VectorStore, llm core.Gemini) (ImportResult, error) {
	f, err := excelize.OpenReader(bytes.NewReader(xlsxBytes))
	if err != nil {
		return ImportResult{Imported: 0, Errors: []string{"Empty workbook"}, Items: []core.MenuItem{}}, nil
	}
	defer func() { _ = f.Close() }()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return ImportResult{Imported: 0, Errors: []string{"Empty workbook"}, Items: []core.MenuItem{}}, nil
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return ImportResult{}, err
	}
	if len(rows) == 0 {
		return ImportResult{Imported: 0, Errors: []string{"Empty workbook"}, Items: []core.MenuItem{}}, nil
	}

	headers := rows[0]
	mapping := headerMap(headers)
	if _, ok := mapping["name"]; !ok {
		return ImportResult{Imported: 0, Errors: []string{"Missing required column: name (tên món). ", "Found headers: " + strings.Join(headers, ", ")}, Items: []core.MenuItem{}}, nil
	}

	res := ImportResult{Imported: 0, Errors: []string{}, Items: []core.MenuItem{}}

	for ridx := 1; ridx < len(rows); ridx++ {
		row := rows[ridx]
		if isBlankRow(row) {
			continue
		}

		raw := map[string]string{
			"id":          cell(row, mapping["id"], has(mapping, "id")),
			"name":        cell(row, mapping["name"], true),
			"description": cell(row, mapping["description"], has(mapping, "description")),
			"price":       cell(row, mapping["price"], has(mapping, "price")),
			"tags":        cell(row, mapping["tags"], has(mapping, "tags")),
			"allergens":   cell(row, mapping["allergens"], has(mapping, "allergens")),
			"ingredients": cell(row, mapping["ingredients"], has(mapping, "ingredients")),
		}

		item := NormalizeMenuItemFromMap(raw)
		if item.Name == "" {
			res.Errors = append(res.Errors, "Row "+strconv.Itoa(ridx+1)+": missing name")
			continue
		}

		created, err := fs.UpsertMenuItem(restaurantID, item)
		if err != nil {
			return ImportResult{}, err
		}
		_ = EmbedMenuItem(ctx, restaurantID, created, vector, llm)
		res.Items = append(res.Items, created)
		res.Imported++
	}

	return res, nil
}
