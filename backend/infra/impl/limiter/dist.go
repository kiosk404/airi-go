package limiter

import (
	"github.com/kiosk404/airi-go/backend/infra/contract/cache"
	"github.com/kiosk404/airi-go/backend/infra/contract/limiter"
)

// NewRateLimiterFactory
// NewRateLimiter `Rule.KeyExpr/Rule.Match` are configured based on the `expr-lang` syntax.
// In the `AllowN` method, tags are matched against expressions using `Rule.Match`.
// When multiple rules are matched, the first rule will be used.
func NewRateLimiterFactory(cmdable cache.Cmdable, opts ...FactoryOpt) limiter.IRateLimiterFactory {
	opt := &factoryOpt{}
	for _, fn := range opts {
		fn(opt)
	}

	return &factory{
		cmdable:   cmdable,
		cacheSize: opt.exprCacheSize,
	}
}

type FactoryOpt func(opt *factoryOpt)

// WithExprCacheSize size is in bytes.
func WithExprCacheSize(size int) FactoryOpt {
	return func(c *factoryOpt) {
		c.exprCacheSize = size
	}
}

type factoryOpt struct {
	exprCacheSize int
}
