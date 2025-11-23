package ptr

import (
	"fmt"
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

func PtrConvert[T1, T2 any](pA *T1) *T2 {
	if pA == nil {
		return nil
	}

	t1 := reflect.TypeOf(*pA)
	t2 := reflect.TypeOf(*new(T2))

	if !checkUnderlyingCompatibility(t1, t2) {
		fmt.Printf("Error: Types are structurally incompatible. T1: %s, T2: %s\n", t1, t2)
		return nil
	}

	vA := reflect.ValueOf(pA).Elem()
	vB := reflect.New(t2).Elem()

	if t1.Kind() == reflect.Struct {
		// 结构体：按字段安全赋值
		for i := 0; i < t1.NumField(); i++ {
			vB.Field(i).Set(vA.Field(i))
		}
	} else {
		if vA.Type().ConvertibleTo(t2) {
			vB.Set(vA.Convert(t2))
		} else {
			vB.Set(vA)
		}
	}

	return vB.Addr().Interface().(*T2)
}

func OfConvert[T1, T2 any](pA T1) T2 {
	// 获取 T1 和 T2 的 Type
	t1 := reflect.TypeOf(pA)
	t2 := reflect.TypeOf(*new(T2))

	// 【核心校验】检查底层结构是否严格兼容
	if !checkUnderlyingCompatibility(t1, t2) {
		fmt.Printf("Error: Types are structurally incompatible. T1: %s, T2: %s\n", t1, t2)
		// 返回 T2 的零值
		var zero T2
		return zero
	}

	// 获取 T1 的值
	vA := reflect.ValueOf(pA)

	// 创建 T2 的 Value 对象
	vB := reflect.New(t2).Elem()

	if t1.Kind() == reflect.Struct {
		// 结构体：按字段安全赋值
		for i := 0; i < t1.NumField(); i++ {
			vB.Field(i).Set(vA.Field(i))
		}
	} else {
		if vA.Type().ConvertibleTo(t2) {
			vB.Set(vA.Convert(t2))
		} else {
			vB.Set(vA)
		}
	}

	return vB.Interface().(T2)
}

func FromPtrConvert[T1, T2 any](pA *T1) T2 {
	if pA == nil {
		var zero T2
		return zero
	}

	// 获取 T1 和 T2 的 Type
	t1 := reflect.TypeOf(*pA)
	t2 := reflect.TypeOf(*new(T2))

	if !checkUnderlyingCompatibility(t1, t2) {
		fmt.Printf("Error: Types are structurally incompatible. T1: %s, T2: %s\n", t1, t2)
		// 返回 T2 的零值
		var zero T2
		return zero
	}

	// 结构兼容，执行转换（通过值中转）
	// 获取 T1 的值
	vA := reflect.ValueOf(*pA)

	// 创建 T2 的 Value 对象
	vB := reflect.New(t2).Elem()

	if t1.Kind() == reflect.Struct {
		// 结构体：按字段安全赋值
		for i := 0; i < t1.NumField(); i++ {
			vB.Field(i).Set(vA.Field(i))
		}
	} else {
		if vA.Type().ConvertibleTo(t2) {
			vB.Set(vA.Convert(t2))
		} else {
			vB.Set(vA)
		}
	}

	return vB.Interface().(T2)
}

func checkUnderlyingCompatibility(t1, t2 reflect.Type) bool {
	// 如果是相同类型，直接返回 true
	if t1 == t2 {
		return true
	}

	k1 := t1.Kind()
	k2 := t2.Kind()

	// 1. 基础类型（包括 int, string, bool 等）必须一致，别名也算
	if k1 != k2 {
		return false
	}

	switch k1 {
	case reflect.Ptr:
		// 如果是指针，递归检查其元素类型
		return checkUnderlyingCompatibility(t1.Elem(), t2.Elem())

	case reflect.Array:
		// 数组要求长度和元素类型都一致
		if t1.Len() != t2.Len() {
			return false
		}
		return checkUnderlyingCompatibility(t1.Elem(), t2.Elem())

	case reflect.Slice, reflect.Chan:
		// 切片和通道只要求元素类型一致
		return checkUnderlyingCompatibility(t1.Elem(), t2.Elem())

	case reflect.Map:
		// 映射要求 Key 和 Value 类型都一致
		if !checkUnderlyingCompatibility(t1.Key(), t2.Key()) {
			return false
		}
		return checkUnderlyingCompatibility(t1.Elem(), t2.Elem())

	case reflect.Struct:
		// 结构体要求字段数量一致
		if t1.NumField() != t2.NumField() {
			return false
		}
		// 递归检查所有字段：名称、类型和偏移量（保证顺序）
		for i := 0; i < t1.NumField(); i++ {
			f1 := t1.Field(i)
			f2 := t2.Field(i)

			// 字段名称和偏移量必须一致（确保布局和语义一致）
			if f1.Name != f2.Name || f1.Offset != f2.Offset {
				return false
			}

			// 递归检查字段类型
			if !checkUnderlyingCompatibility(f1.Type, f2.Type) {
				return false
			}
		}
		return true

	default:
		// 其他类型 (如 interface, func, unsafe.Pointer, 或原始类型)
		// 只要 Kind 一致，且非结构体/复合类型，就认为兼容。
		return true
	}
}

func FromOrDefault[T any](p *T, def T) T {
	if p != nil {
		return *p
	}
	return def
}

func PtrConvertMap[F any, T any](f *F, c func(f F) T) *T {
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
