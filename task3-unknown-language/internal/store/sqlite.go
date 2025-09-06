package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	_ "modernc.org/sqlite" // pure Go sqlite driver

	"github.com/yourname/inventory-go/internal/model"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type Store struct {
	DB *sql.DB
}

func Open(dsn string) (*sql.DB, error) {
	// Examples:
	//   file:inventory.db?cache=shared&_pragma=foreign_keys(1)
	//   file::memory:?cache=shared&_pragma=foreign_keys(1)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	// sanity pragma (in addition to DSN)
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

const schema0001 = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS items (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    name     TEXT NOT NULL UNIQUE,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    price    REAL NOT NULL CHECK (price > 0),
    tags     TEXT NOT NULL,
    status   TEXT NOT NULL CHECK (status IN ('active','archived'))
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_items_name ON items(name);
`

func (s *Store) Migrate(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, schema0001)
	return err
}

func New(db *sql.DB) *Store { return &Store{DB: db} }

func (s *Store) CreateItem(ctx context.Context, in model.ItemCreate) (model.Item, error) {
	tagsJSON, err := json.Marshal(in.Tags)
	if err != nil {
		return model.Item{}, err
	}
	res, err := s.DB.ExecContext(ctx, `
INSERT INTO items (name, quantity, price, tags, status)
VALUES (?, ?, ?, ?, ?)
`, in.Name, in.Quantity, in.Price, string(tagsJSON), in.Status)
	if err != nil {
		if isUniqueErr(err) {
			return model.Item{}, ErrConflict
		}
		return model.Item{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return model.Item{}, err
	}
	return s.GetItemByID(ctx, int(id))
}

func (s *Store) GetItemByID(ctx context.Context, id int) (model.Item, error) {
	var it model.Item
	var tagsText string
	row := s.DB.QueryRowContext(ctx, `
SELECT id, name, quantity, price, tags, status
FROM items WHERE id = ?
`, id)
	if err := row.Scan(&it.ID, &it.Name, &it.Quantity, &it.Price, &tagsText, &it.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Item{}, ErrNotFound
		}
		return model.Item{}, err
	}
	if err := json.Unmarshal([]byte(tagsText), &it.Tags); err != nil {
		// Если по какой-то причине испортили JSON — считаем это серверной ошибкой
		return model.Item{}, fmt.Errorf("decode tags: %w", err)
	}
	return it, nil
}

func (s *Store) ListItems(ctx context.Context, limit, offset int) ([]model.Item, error) {
	rows, err := s.DB.QueryContext(ctx, `
SELECT id, name, quantity, price, tags, status
FROM items
ORDER BY id ASC
LIMIT ? OFFSET ?
`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Item
	for rows.Next() {
		var it model.Item
		var tagsText string
		if err := rows.Scan(&it.ID, &it.Name, &it.Quantity, &it.Price, &tagsText, &it.Status); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(tagsText), &it.Tags); err != nil {
			return nil, fmt.Errorf("decode tags: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (s *Store) UpdateItemPartial(ctx context.Context, id int, in model.ItemUpdate) (model.Item, error) {
	// Предварительная проверка на конфликт имени (если пришло name)
	if in.Name != nil {
		var existingID int
		err := s.DB.QueryRowContext(ctx, `SELECT id FROM items WHERE name = ? AND id <> ?`, *in.Name, id).Scan(&existingID)
		if err == nil {
			return model.Item{}, ErrConflict
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return model.Item{}, err
		}
	}

	sets := make([]string, 0, 5)
	args := make([]any, 0, 6)

	if in.Name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *in.Name)
	}
	if in.Quantity != nil {
		sets = append(sets, "quantity = ?")
		args = append(args, *in.Quantity)
	}
	if in.Price != nil {
		sets = append(sets, "price = ?")
		args = append(args, *in.Price)
	}
	if in.Tags != nil {
		tagsJSON, err := json.Marshal(*in.Tags)
		if err != nil {
			return model.Item{}, err
		}
		sets = append(sets, "tags = ?")
		args = append(args, string(tagsJSON))
	}
	if in.Status != nil {
		sets = append(sets, "status = ?")
		args = append(args, *in.Status)
	}

	if len(sets) == 0 {
		// Нечего обновлять — возвращаем текущее состояние
		return s.GetItemByID(ctx, id)
	}

	args = append(args, id)
	q := fmt.Sprintf("UPDATE items SET %s WHERE id = ?", strings.Join(sets, ", "))
	_, err := s.DB.ExecContext(ctx, q, args...)
	if err != nil {
		if isUniqueErr(err) {
			return model.Item{}, ErrConflict
		}
		return model.Item{}, err
	}
	return s.GetItemByID(ctx, id)
}

func (s *Store) DeleteItem(ctx context.Context, id int) error {
	res, err := s.DB.ExecContext(ctx, `DELETE FROM items WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func isUniqueErr(err error) bool {
	msg := err.Error()
	// modernc.org/sqlite возвращает текст вида:
	// "UNIQUE constraint failed: items.name (1555)"
	return strings.Contains(strings.ToLower(msg), "unique constraint failed")
}
