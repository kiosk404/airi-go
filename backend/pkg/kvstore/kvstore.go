package kvstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kiosk404/airi-go/backend/pkg/json"
	"gorm.io/gorm"
)

var ErrKeyNotFound = errors.New("key not found")

type KVStore[T any] struct {
	repo *gorm.DB
}

var defaultDB *gorm.DB

func SetDefault(db *gorm.DB) {
	defaultDB = db
}

func New[T any](db *gorm.DB) *KVStore[T] {
	return &KVStore[T]{
		repo: db,
	}
}

func (g *KVStore[T]) db(ctx context.Context) *gorm.DB {
	if g.repo == nil {
		return defaultDB.WithContext(ctx)
	}

	return g.repo.WithContext(ctx)
}

func (g *KVStore[T]) Save(ctx context.Context, namespace, k string, v *T) error {
	if v == nil {
		return fmt.Errorf("cannot save nil value for key: %s", k)
	}

	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal failed for key %s for type %T: %w", k, *v, err)
	}

	res := g.db(ctx).Exec(
		"INSERT INTO `kv_entries` (`namespace`, `key_data`, `value_data`) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE `value_data` = ?",
		namespace, k, data, data,
	)

	if res.Error != nil {
		return fmt.Errorf("failed to save key %s: %w", k, res.Error)
	}

	return nil
}

func (g *KVStore[T]) Get(ctx context.Context, namespace, k string) (*T, error) {
	var obj T

	row := g.db(ctx).Raw(
		"SELECT `value_data` FROM `kv_entries` WHERE `namespace` = ? AND `key_data` = ? LIMIT 1",
		namespace, k,
	).Row()

	var value []byte
	if err := row.Scan(&value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrKeyNotFound
		}
		return nil, fmt.Errorf("failed to get key %s: %w", k, err)
	}

	if err := json.Unmarshal(value, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json for key %s: %w", k, err)
	}

	return &obj, nil
}

func (g *KVStore[T]) Delete(ctx context.Context, namespace, k string) error {
	res := g.db(ctx).Exec(
		"DELETE FROM `kv_entries` WHERE `namespace` = ? AND `key_data` = ?",
		namespace, k,
	)

	return res.Error
}
