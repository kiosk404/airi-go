package mcache

import (
	"time"
)

type IByteCache interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte, expiration time.Duration) error
}
