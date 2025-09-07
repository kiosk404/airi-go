package internal

import (
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
)

type Option struct {
	tx          rdb.RDB
	debug       bool
	withMaster  bool
	withDeleted bool
	forUpdate   bool
}

func (o *Option) WithTransaction(tx rdb.RDB) rdb.Option {
	o.tx = tx
	return nil
}

func (o *Option) Transaction() rdb.RDB {
	return o.tx
}

func (o *Option) Debug() rdb.Option {
	o.debug = true
	return nil
}

func (o *Option) IsDebug() bool {
	return o.debug
}

func (o *Option) WithDeleted() rdb.Option {
	o.withDeleted = true
	return nil
}

func (o *Option) IsDeleted() bool {
	return o.withDeleted
}

func (o *Option) WithMaster() rdb.Option {
	o.withMaster = true
	return nil
}

func (o *Option) IsMaster() bool {
	return o.withMaster
}

func (o *Option) WithSelectForUpdate() rdb.Option {
	o.forUpdate = true
	return nil
}

func (o *Option) IsSelectForUpdate() bool {
	return o.forUpdate
}
