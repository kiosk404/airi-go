package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromConvert(t *testing.T) {
	// 测试用例1: int 转 int (相同类型)
	t.Run("int to int conversion", func(t *testing.T) {
		val := 42
		result := FromConvert[int, int](&val)
		assert.NotNil(t, result)
		assert.Equal(t, 42, *result)
	})

	// 测试用例2: string 转 string (相同类型)
	t.Run("string to string conversion", func(t *testing.T) {
		val := "hello"
		result := FromConvert[string, string](&val)
		assert.NotNil(t, result)
		assert.Equal(t, "hello", *result)
	})

	// 测试用例3: nil 指针输入
	t.Run("nil pointer input", func(t *testing.T) {
		result := FromConvert[int, int](nil)
		assert.Nil(t, result)
	})

	// 测试用例4: 结构体转换 (相同字段)
	t.Run("struct with same fields conversion", func(t *testing.T) {
		type StructA struct {
			Name string
			Age  int
		}
		type StructB struct {
			Name string
			Age  int
		}
		val := StructA{Name: "Alice", Age: 30}
		result := FromConvert[StructA, StructB](&val)
		assert.NotNil(t, result)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})

	// 测试用例5: 结构体别名转换
	t.Run("struct alias conversion", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		type User = Person // 别名
		val := Person{Name: "Bob", Age: 25}
		result := FromConvert[Person, User](&val)
		assert.NotNil(t, result)
		assert.Equal(t, "Bob", result.Name)
		assert.Equal(t, 25, result.Age)
	})

	// 测试用例6: 数组转换
	t.Run("array conversion", func(t *testing.T) {
		val := [3]int{1, 2, 3}
		result := FromConvert[[3]int, [3]int](&val)
		assert.NotNil(t, result)
		assert.Equal(t, [3]int{1, 2, 3}, *result)
	})

	// 测试用例7: 切片转换
	t.Run("slice conversion", func(t *testing.T) {
		val := []int{1, 2, 3}
		result := FromConvert[[]int, []int](&val)
		assert.NotNil(t, result)
		assert.Equal(t, []int{1, 2, 3}, *result)
	})

	// 测试用例8: map转换
	t.Run("map conversion", func(t *testing.T) {
		val := map[string]int{"a": 1, "b": 2}
		result := FromConvert[map[string]int, map[string]int](&val)
		assert.NotNil(t, result)
		assert.Equal(t, map[string]int{"a": 1, "b": 2}, *result)
	})

	// 测试用例9: 不兼容类型转换
	t.Run("incompatible types conversion", func(t *testing.T) {
		val := 42
		result := FromConvert[int, string](&val)
		assert.Nil(t, result) // 无法转换 int 到 string
	})

	// 测试用例10: 不兼容结构体转换
	t.Run("incompatible struct conversion", func(t *testing.T) {
		type StructA struct {
			Name string
		}
		type StructB struct {
			Age int
		}
		val := StructA{Name: "Alice"}
		result := FromConvert[StructA, StructB](&val)
		assert.Nil(t, result) // 字段不兼容
	})

	// 测试用例11: 指针类型转换
	t.Run("pointer conversion", func(t *testing.T) {
		value := 42
		val := &value
		result := FromConvert[*int, *int](&val)
		assert.NotNil(t, result)
		assert.Equal(t, 42, **result)
	})

	// 测试用例12: 复杂结构体转换
	t.Run("complex struct conversion", func(t *testing.T) {
		type Address struct {
			Street string
			City   string
		}
		type Person struct {
			Name    string
			Age     int
			Address Address
		}
		type User struct {
			Name    string
			Age     int
			Address Address
		}
		val := Person{
			Name: "Charlie",
			Age:  35,
			Address: Address{
				Street: "123 Main St",
				City:   "New York",
			},
		}
		result := FromConvert[Person, User](&val)
		assert.NotNil(t, result)
		assert.Equal(t, "Charlie", result.Name)
		assert.Equal(t, 35, result.Age)
		assert.Equal(t, "123 Main St", result.Address.Street)
		assert.Equal(t, "New York", result.Address.City)
	})
}