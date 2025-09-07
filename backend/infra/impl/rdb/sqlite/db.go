package sqlite

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/impl/rdb/internal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// provider 包装 gorm.db 并强制提供 ctx 以串联 trace
type provider struct {
	db *gorm.DB
}

var _ rdb.Provider = &provider{}

// NewDB 从配置创建一个 db 实例
func NewDB(cfg *Config, opts ...gorm.Option) (rdb.Provider, error) {
	opts = append(opts, &gorm.Config{
		TranslateError: true,
	})

	db, err := gorm.Open(sqlite.Open(cfg.buildDSN()), opts...)
	if err != nil {
		return nil, err
	}

	return &provider{db: db}, nil
}

func (p *provider) NewSession(ctx context.Context, opts ...rdb.Option) rdb.RDB {
	session := p.db

	opt := &internal.Option{}
	for _, fn := range opts {
		fn(opt)
	}
	if opt.Transaction() != nil {
		session = opt.Transaction().DB()
	}
	if opt.IsDebug() {
		session = session.Debug()
	}
	if opt.IsDeleted() {
		session = session.Unscoped()
	}
	if opt.IsSelectForUpdate() {
		// SQLite 对 FOR UPDATE 支持有限，这里保留但可能需要根据实际需求调整
		session = session.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	return &sqliteService{db: session.WithContext(ctx)}
}
func (p *provider) Transaction(ctx context.Context, fc func(tx rdb.RDB) error, opts ...rdb.Option) error {
	session := p.NewSession(ctx, opts...)
	return session.Transaction(ctx, fc)
}
