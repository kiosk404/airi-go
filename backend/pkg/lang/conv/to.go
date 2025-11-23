package conv

import (
	"encoding/json"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/bytedance/gg/gconv"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/spf13/cast"
)

// StrToInt64E returns strconv.ParseInt(v, 10, 64)
func StrToInt64(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}

// Int64ToStr returns strconv.FormatInt(v, 10) result
func Int64ToStr(v int64) string {
	return strconv.FormatInt(v, 10)
}

func StrToFloat64(v string) (float64, error) {
	return strconv.ParseFloat(v, 64)
}

// StrToInt64 returns strconv.ParseInt(v, 10, 64)'s value.
// if error occurs, returns defaultValue as result.
func StrToInt64D(v string, defaultValue int64) int64 {
	toV, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultValue
	}
	return toV
}

// DebugJsonToStr
func DebugJsonToStr(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func BoolToInt(p bool) int {
	if p == true {
		return 1
	}

	return 0
}

// BoolToIntPointer returns 1 or 0 as pointer
func BoolToIntPointer(p *bool) *int {
	if p == nil {
		return nil
	}

	if *p == true {
		return ptr.Of(int(1))
	}

	return ptr.Of(int(0))
}

func UnsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// UnsafeStringToBytes
//
//nolint:staticcheck
func UnsafeStringToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

func ToBool(v any) bool {
	return cast.ToBool(v)
}

func ToString(v any) string {
	return cast.ToString(v)
}

// Int64 will convert the given value to a int64, returns the default value of 0
// if a conversion can not be made.
func Int64(from interface{}) (int64, error) {
	return gconv.ToE[int64, any](from)
}
