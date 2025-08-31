package ginutil

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetInt32 gets a int32 value from path, query, or form.
func GetInt32(c *gin.Context, key string) int32 {
	valStr, ok := getParamValue(c, key)
	if !ok {
		return 0
	}

	val, err := strconv.ParseInt(valStr, 10, 32)
	if err != nil {
		return 0
	}
	return int32(val)
}

// GetInt64 gets a int64 value from path, query, or form.
func GetInt64(c *gin.Context, key string) int64 {
	valStr, ok := getParamValue(c, key)
	if !ok {
		return 0
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

// getParamValue tries to get a value from path, query, then form.
func getParamValue(c *gin.Context, key string) (string, bool) {
	// Try to get from path parameter first
	if val := c.Param(key); val != "" {
		return val, true
	}

	// Then try to get from query parameter
	if val := c.Query(key); val != "" {
		return val, true
	}

	// Finally try to get from form data
	if val := c.PostForm(key); val != "" {
		return val, true
	}

	return "", false
}
