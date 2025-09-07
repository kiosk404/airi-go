package mysql

import (
	"context"

	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
	"github.com/kiosk404/airi-go/backend/infra/impl/rdb/internal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/utils"
	"gorm.io/plugin/dbresolver"
)

// provider 包装 gorm.db 并强制提供 ctx 以串联 trace
type provider struct {
	db *gorm.DB
}

var _ rdb.Provider = &provider{}

// NewDB 从配置创建一个 db 实例
func NewDB(cfg *Config, opts ...gorm.Option) (rdb.Provider, error) {
	if !utils.Contains(mysql.UpdateClauses, "RETURNING") {
		mysql.UpdateClauses = append(mysql.UpdateClauses, "RETURNING")
	}
	// Known issue: this option will make the opts using gorm.Config not working.
	opts = append(opts, &gorm.Config{
		TranslateError: true,
	})

	db, err := gorm.Open(mysql.Open(cfg.buildDSN()), opts...)
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
	if opt.IsMaster() {
		session = session.Clauses(dbresolver.Write)
	}
	if opt.IsDeleted() {
		session = session.Unscoped()
	}
	if opt.IsSelectForUpdate() {
		session = session.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	return &mysqlService{db: session.WithContext(ctx)}
}

func (p *provider) Transaction(ctx context.Context, fc func(tx rdb.RDB) error, opts ...rdb.Option) error {
	session := p.NewSession(ctx, opts...)
	return session.Transaction(ctx, fc)
}

func ContainWithMasterOpt(opt []rdb.Option) bool {
	o := &internal.Option{}
	for _, fn := range opt {
		fn(o)
		if o.IsMaster() {
			return true
		}
	}
	return false
}
