package option

import (
	"github.com/kiosk404/airi-go/backend/infra/contract/rdb"
)

var _ rdb.OptionService = &Option{}

type Option struct{}

func (o *Option) WithTransaction(tx rdb.RDB) rdb.Option {
	return func(option rdb.OptionService) {
		option.WithTransaction(tx)
	}
}

func (o *Option) Debug() rdb.Option {
	return func(option rdb.OptionService) {
		option.Debug()
	}
}

func (o *Option) WithDeleted() rdb.Option {
	return func(option rdb.OptionService) {
		option.WithDeleted()
	}
}

func (o *Option) WithMaster() rdb.Option {
	return func(option rdb.OptionService) {
		option.WithMaster()
	}
}

func (o *Option) WithSelectForUpdate() rdb.Option {
	return func(option rdb.OptionService) {
		option.WithSelectForUpdate()
	}
}
