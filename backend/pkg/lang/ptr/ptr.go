package ptr

import (
	"reflect"

	"github.com/samber/lo"
)

func Of[T any](t T) *T {
	return lo.ToPtr(t)
}

func From[T any](p *T) T {
	if p != nil {
		return *p
	}
	var t T
	return t
}

func FromOrDefault[T any](p *T, def T) T {
	if p != nil {
		return *p
	}
	return def
}

func PtrConvert[F any, T any](f *F, c func(f F) T) *T {
	if f == nil {
		return nil
	}
	return Of(c(*f))
}

type Integer interface {
	~int64 | ~int32 | ~int16 | ~int8 | ~int
}

func ConvIntPtr[T, K Integer](val *T) *K {
	if val == nil {
		return nil
	}
	return Of((K)(*val))
}

func IsNull[T any](v T) bool {
	return reflect.ValueOf(v).IsZero()
}
