package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	//"fmt"
	"strings"
	//"time"

	//"github.com/vpa/quanlynhahang-backend/models/AIChatBot"

	//"github.com/google/uuid"
	//"github.com/jackc/pgx/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	core "github.com/vpa/quanlynhahang-backend/models/AIChatBot"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func normalizeRestaurantID(restaurantID string) (string, error) {
	restaurantID = strings.TrimSpace(restaurantID)
	if restaurantID == "" {
		return "", errors.New("missing restaurant_id")
	}
	return restaurantID, nil
}

const ensureRestaurantSQL = `INSERT INTO restaurants (id, updated_at) VALUES ($1, now())
ON CONFLICT (id) DO UPDATE SET updated_at = EXCLUDED.updated_at;`

func (s *PostgresStore) ensureRestaurant(ctx context.Context, restaurantID string) {
	// Idempotent, best-effort: keep other fields untouched.
	_, _ = s.pool.Exec(ctx, ensureRestaurantSQL, restaurantID)
}

func (s *PostgresStore) EnsureThread(threadID string) (string, error) {
	threadID = strings.TrimSpace(threadID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Nếu threadID truyền vào và tồn tại thì dùng lại
	if threadID != "" {
		var ok bool
		err := s.pool.QueryRow(ctx, `
			SELECT true 
			FROM threads 
			WHERE id = $1
		`, threadID).Scan(&ok)

		if err == nil && ok {
			return threadID, nil
		}

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
	}

	// tạo thread mới
	id := uuid.NewString()

	_, err := s.pool.Exec(ctx, `
		INSERT INTO threads (id, created_at)
		VALUES ($1, now())
	`, id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *PostgresStore) AppendThreadMessage(threadID, role, content string) error {
	//restaurantID, err := normalizeRestaurantID(restaurantID)
	// if err != nil {
	// 	return err
	// }
	threadID = strings.TrimSpace(threadID)
	role = strings.TrimSpace(role)
	content = strings.TrimSpace(content)
	if threadID == "" {
		return errors.New("missing thread id")
	}
	if role == "" {
		return errors.New("missing role")
	}
	if content == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd, err := s.pool.Exec(ctx, `
INSERT INTO thread_messages (thread_id, role, content, created_at)
SELECT $1, $2, $3, now()
WHERE EXISTS (SELECT 1 FROM threads WHERE id = $1 )
`, threadID, role, content)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("thread not found")
	}
	return nil
}

func (s *PostgresStore) GetThreadMessages(threadID string) ([]core.Message, error) {
	threadID = strings.TrimSpace(threadID)

	if threadID == "" {
		return []core.Message{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx, `
		SELECT role, content
		FROM thread_messages
		WHERE thread_id = $1
		ORDER BY id ASC
	`, threadID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []core.Message{}
	for rows.Next() {
		var role, content string
		if err := rows.Scan(&role, &content); err != nil {
			return nil, err
		}
		out = append(out, core.Message{
			Role:    role,
			Content: content,
		})
	}

	return out, rows.Err()
}

// func (s *PostgresStore) UpsertMenuItem(restaurantID string, item core.MenuItem) (core.MenuItem, error) {
// 	restaurantID, err := normalizeRestaurantID(restaurantID)
// 	if err != nil {
// 		return core.MenuItem{}, err
// 	}
// 	item.Name = strings.TrimSpace(item.Name)
// 	if item.Name == "" {
// 		return core.MenuItem{}, errors.New("missing name")
// 	}
// 	if strings.TrimSpace(item.ID) == "" {
// 		item.ID = uuid.NewString()
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
// 	defer cancel()

// 	s.ensureRestaurant(ctx, restaurantID)

// 	// NOTE(single-restaurant mode): we intentionally do not validate tenant for menu_items.

// 	tagsJSON, _ := json.Marshal(item.Tags)
// 	allergensJSON, _ := json.Marshal(item.Allergens)
// 	ingredientsJSON, _ := json.Marshal(item.Ingredients)

// 	_, err = s.pool.Exec(ctx, `
// INSERT INTO menu_items (id, restaurant_id, name, description, price, tags, allergens, ingredients, updated_at)
// VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now())
// ON CONFLICT (id) DO UPDATE SET
// 	restaurant_id = EXCLUDED.restaurant_id,
//   name = EXCLUDED.name,
//   description = EXCLUDED.description,
//   price = EXCLUDED.price,
//   tags = EXCLUDED.tags,
//   allergens = EXCLUDED.allergens,
//   ingredients = EXCLUDED.ingredients,
//   updated_at = now()
// `, item.ID, restaurantID, item.Name, item.Description, item.Price, tagsJSON, allergensJSON, ingredientsJSON)
// 	if err != nil {
// 		return core.MenuItem{}, err
// 	}
// 	return item, nil
// }

// func (s *PostgresStore) ListMenuItems(restaurantID string, offset, limit int) ([]core.MenuItem, int, error) {
// 	// NOTE(single-restaurant mode): ignore restaurantID and return global menu_items.
// 	if offset < 0 {
// 		offset = 0
// 	}
// 	if limit <= 0 {
// 		limit = 50
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
// 	defer cancel()

// 	var total int
// 	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM menu_items`).Scan(&total); err != nil {
// 		return nil, 0, err
// 	}

// 	rows, err := s.pool.Query(ctx, `
// SELECT id, name, description, price, tags, allergens, ingredients
// FROM menu_items
// ORDER BY lower(name)
// OFFSET $1
// LIMIT $2
// `, offset, limit)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	defer rows.Close()

// 	items := []core.MenuItem{}
// 	for rows.Next() {
// 		var it core.MenuItem
// 		var tagsB, allergensB, ingredientsB []byte
// 		if err := rows.Scan(&it.ID, &it.Name, &it.Description, &it.Price, &tagsB, &allergensB, &ingredientsB); err != nil {
// 			return nil, 0, err
// 		}
// 		_ = json.Unmarshal(tagsB, &it.Tags)
// 		_ = json.Unmarshal(allergensB, &it.Allergens)
// 		_ = json.Unmarshal(ingredientsB, &it.Ingredients)
// 		items = append(items, it)
// 	}
// 	return items, total, rows.Err()
// }

// func (s *PostgresStore) DeleteMenuItem(id string) (bool, error) {
// 	// NOTE(single-restaurant mode): ignore restaurantID and delete by id only.
// 	id = strings.TrimSpace(id)
// 	if id == "" {
// 		return false, nil
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
// 	defer cancel()

// 	cmd, err := s.pool.Exec(ctx, `DELETE FROM menu_items WHERE id = $1`, id)
// 	if err != nil {
// 		return false, err
// 	}
// 	return cmd.RowsAffected() > 0, nil
// }

// func (s *PostgresStore) SetRestaurant(restaurantID string, info core.RestaurantInfo) (core.RestaurantInfo, error) {
// 	restaurantID, err := normalizeRestaurantID(restaurantID)
// 	if err != nil {
// 		return core.RestaurantInfo{}, err
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
// 	defer cancel()

// 	_, err = s.pool.Exec(ctx, `
// INSERT INTO restaurants (id, name, address, open_hours, phone, style, policies, updated_at)
// VALUES ($1, $2, $3, $4, $5, $6, $7, now())
// ON CONFLICT (id) DO UPDATE SET
//   name = EXCLUDED.name,
//   address = EXCLUDED.address,
//   open_hours = EXCLUDED.open_hours,
//   phone = EXCLUDED.phone,
//   style = EXCLUDED.style,
//   policies = EXCLUDED.policies,
//   updated_at = now()
// `, restaurantID, info.Name, info.Address, info.OpenHours, info.Phone, info.Style, info.Policies)
// 	if err != nil {
// 		return core.RestaurantInfo{}, err
// 	}
// 	return info, nil
// }

// func (s *PostgresStore) GetRestaurant(restaurantID string) (core.RestaurantInfo, error) {
// 	restaurantID, err := normalizeRestaurantID(restaurantID)
// 	if err != nil {
// 		return core.RestaurantInfo{}, err
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
// 	defer cancel()

// 	var out core.RestaurantInfo
// 	err = s.pool.QueryRow(ctx, `
// SELECT name, address, open_hours, phone, style, policies
// FROM restaurants
// WHERE id = $1
// `, restaurantID).Scan(&out.Name, &out.Address, &out.OpenHours, &out.Phone, &out.Style, &out.Policies)
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return core.RestaurantInfo{}, nil
// 		}
// 		return core.RestaurantInfo{}, err
// 	}
// 	return out, nil
// }
