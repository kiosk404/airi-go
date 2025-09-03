package local

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/kiosk404/airi-go/backend/infra/contract/cache"
)

var Nil error = fmt.Errorf("ristretto: nil")

func SetDefaultNilError(err error) {
	Nil = err
}

func New() (cache.Cmdable, error) {
	return NewRistrettoClient()
}

type CacheType string

var (
	CacheTypeString CacheType = "string"
	CacheTypeHash   CacheType = "hash"
	CacheTypeList   CacheType = "list"
)

// CacheValue 统一的缓存值包装器
type CacheValue struct {
	Data     interface{}
	Type     CacheType // "string", "hash", "list"
	CreateAt time.Time
	TTL      time.Duration
}

// RistrettoClient Ristretto v2缓存客户端
type RistrettoClient struct {
	cache    *ristretto.Cache[string, *CacheValue]
	hashData map[string]map[string]string // 模拟Redis Hash数据结构
	listData map[string][]string          // 模拟Redis List数据结构
	mu       sync.RWMutex
}

// NewRistrettoClient 创建新的Ristretto v2客户端
func NewRistrettoClient() (*RistrettoClient, error) {
	c, err := ristretto.NewCache(&ristretto.Config[string, *CacheValue]{
		NumCounters: 1e7,     // 跟踪频率的键数量
		MaxCost:     1 << 30, // 最大缓存大小(1GB)
		BufferItems: 64,      // 接受键的缓冲区项数量
	})
	if err != nil {
		return nil, err
	}

	client := &RistrettoClient{
		cache:    c,
		hashData: make(map[string]map[string]string),
		listData: make(map[string][]string),
	}

	// 启动过期清理goroutine
	go client.expireCleanup()

	return client, nil
}

// expireCleanup 定期清理过期数据
func (c *RistrettoClient) expireCleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()

		// 清理Hash数据
		for key := range c.hashData {
			if value, found := c.cache.Get(key); found {
				if !value.CreateAt.IsZero() && now.After(value.CreateAt.Add(value.TTL)) {
					c.cache.Del(key)
					delete(c.hashData, key)
				}
			} else {
				delete(c.hashData, key)
			}
		}

		// 清理List数据
		for key := range c.listData {
			if value, found := c.cache.Get(key); found {
				if !value.CreateAt.IsZero() && now.After(value.CreateAt.Add(value.TTL)) {
					c.cache.Del(key)
					delete(c.listData, key)
				}
			} else {
				delete(c.listData, key)
			}
		}

		c.mu.Unlock()
	}
}

// isExpired 检查键是否过期
func (c *RistrettoClient) isExpired(key string) bool {
	value, found := c.cache.Get(key)
	if !found {
		return true
	}

	now := time.Now()
	if !value.CreateAt.IsZero() && now.After(value.CreateAt.Add(value.TTL)) {
		c.cache.Del(key)
		c.mu.Lock()
		delete(c.hashData, key)
		delete(c.listData, key)
		c.mu.Unlock()
		return true
	}

	return false
}

// Pipeline 创建管道
func (c *RistrettoClient) Pipeline() cache.Pipeliner {
	return &RistrettoPipeline{
		client:   c,
		commands: make([]PipelineCmd, 0),
	}
}

// 命令结果实现
type cmdResult struct {
	err error
}

func (c *cmdResult) Err() error {
	return c.err
}

// cache.IntCmd 实现
type intCmd struct {
	cmdResult
	val int64
}

func (c *intCmd) Result() (int64, error) {
	return c.val, c.err
}

func newIntCmd(val int64, err error) cache.IntCmd {
	return &intCmd{
		cmdResult: cmdResult{err: err},
		val:       val,
	}
}

// cache.StringCmd 实现
type stringCmd struct {
	cmdResult
	val string
}

func (c *stringCmd) Result() (string, error) {
	return c.val, c.err
}

func (c *stringCmd) Val() string {
	return c.val
}

func (c *stringCmd) Int64() (int64, error) {
	if c.err != nil {
		return 0, c.err
	}
	return strconv.ParseInt(c.val, 10, 64)
}

func (c *stringCmd) Bytes() ([]byte, error) {
	if c.err != nil {
		return nil, c.err
	}
	return []byte(c.val), nil
}

func newStringCmd(val string, err error) cache.StringCmd {
	return &stringCmd{
		cmdResult: cmdResult{err: err},
		val:       val,
	}
}

// cache.StatusCmd 实现
type statusCmd struct {
	cmdResult
	val string
}

func (c *statusCmd) Result() (string, error) {
	return c.val, c.err
}

func newStatusCmd(val string, err error) cache.StatusCmd {
	return &statusCmd{
		cmdResult: cmdResult{err: err},
		val:       val,
	}
}

// cache.BoolCmd 实现
type boolCmd struct {
	cmdResult
	val bool
}

func (c *boolCmd) Result() (bool, error) {
	return c.val, c.err
}

func newBoolCmd(val bool, err error) cache.BoolCmd {
	return &boolCmd{
		cmdResult: cmdResult{err: err},
		val:       val,
	}
}

// cache.MapStringStringCmd 实现
type mapStringStringCmd struct {
	cmdResult
	val map[string]string
}

func (c *mapStringStringCmd) Result() (map[string]string, error) {
	return c.val, c.err
}

func newMapStringStringCmd(val map[string]string, err error) cache.MapStringStringCmd {
	return &mapStringStringCmd{
		cmdResult: cmdResult{err: err},
		val:       val,
	}
}

// cache.StringSliceCmd 实现
type stringSliceCmd struct {
	cmdResult
	val []string
}

func (c *stringSliceCmd) Result() ([]string, error) {
	return c.val, c.err
}

func newStringSliceCmd(val []string, err error) cache.StringSliceCmd {
	return &stringSliceCmd{
		cmdResult: cmdResult{err: err},
		val:       val,
	}
}

// Set 设置键值对
func (c *RistrettoClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) cache.StatusCmd {
	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	case int, int64, int32:
		strValue = fmt.Sprintf("%v", v)
	case float32, float64:
		strValue = fmt.Sprintf("%v", v)
	default:
		// 对于复杂类型，使用JSON序列化
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return newStatusCmd("", err)
		}
		strValue = string(jsonBytes)
	}

	// 创建缓存值
	cacheValue := &CacheValue{
		Data: strValue,
		Type: CacheTypeString,
	}

	if ttl > 0 {
		cacheValue.TTL = ttl
	}
	cacheValue.CreateAt = time.Now()

	// 使用泛型Cache设置值
	success := c.cache.SetWithTTL(key, cacheValue, cost(cacheValue), cacheValue.TTL)
	if !success {
		return newStatusCmd("", fmt.Errorf("failed to set key: %s", key))
	}

	return newStatusCmd("OK", nil)
}

// Get 获取键值
func (c *RistrettoClient) Get(ctx context.Context, key string) cache.StringCmd {
	if c.isExpired(key) {
		return newStringCmd("", Nil)
	}

	value, found := c.cache.Get(key)
	if !found {
		return newStringCmd("", Nil)
	}

	if value.Type != "string" {
		return newStringCmd("", fmt.Errorf("value is not a string"))
	}

	strValue, ok := value.Data.(string)
	if !ok {
		return newStringCmd("", fmt.Errorf("value is not a string"))
	}

	return newStringCmd(strValue, nil)
}

// IncrBy 按指定值递增
func (c *RistrettoClient) IncrBy(ctx context.Context, key string, increment int64) cache.IntCmd {
	c.mu.Lock()
	defer c.mu.Unlock()

	current := int64(0)
	if value, found := c.cache.Get(key); found && value.Type == "string" {
		if strVal, ok := value.Data.(string); ok {
			if parsed, err := strconv.ParseInt(strVal, 10, 64); err == nil {
				current = parsed
			}
		}
	}

	newValue := current + increment
	newValueStr := strconv.FormatInt(newValue, 10)

	// 创建新的缓存值
	cacheValue := &CacheValue{
		Data: newValueStr,
		Type: "string",
	}

	c.cache.Set(key, cacheValue, cost(cacheValue))
	return newIntCmd(newValue, nil)
}

// Incr 递增1
func (c *RistrettoClient) Incr(ctx context.Context, key string) cache.IntCmd {
	return c.IncrBy(ctx, key, 1)
}

// HashCmdable 实现

// HSet 设置Hash字段
func (c *RistrettoClient) HSet(ctx context.Context, key string, values ...interface{}) cache.IntCmd {
	if len(values)%2 != 0 {
		return newIntCmd(0, fmt.Errorf("wrong number of arguments"))
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.hashData[key] == nil {
		c.hashData[key] = make(map[string]string)
	}

	setCount := int64(0)
	for i := 0; i < len(values); i += 2 {
		field := fmt.Sprintf("%v", values[i])
		value := fmt.Sprintf("%v", values[i+1])

		if _, exists := c.hashData[key][field]; !exists {
			setCount++
		}
		c.hashData[key][field] = value
	}

	// 在缓存中标记为hash类型
	cacheValue := &CacheValue{
		Data: c.hashData[key],
		Type: CacheTypeHash,
	}
	c.cache.Set(key, cacheValue, cost(cacheValue))

	return newIntCmd(setCount, nil)
}

// HGetAll 获取Hash所有字段
func (c *RistrettoClient) HGetAll(ctx context.Context, key string) cache.MapStringStringCmd {
	if c.isExpired(key) {
		return newMapStringStringCmd(make(map[string]string), nil)
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if hashMap, exists := c.hashData[key]; exists {
		// 创建副本避免并发问题
		result := make(map[string]string)
		for k, v := range hashMap {
			result[k] = v
		}
		return newMapStringStringCmd(result, nil)
	}

	return newMapStringStringCmd(make(map[string]string), nil)
}

// GenericCmdable 实现

// Del 删除键
func (c *RistrettoClient) Del(ctx context.Context, keys ...string) cache.IntCmd {
	c.mu.Lock()
	defer c.mu.Unlock()

	deletedCount := int64(0)
	for _, key := range keys {
		if _, found := c.cache.Get(key); found {
			c.cache.Del(key)
			deletedCount++
		}

		// 同时删除Hash和List数据
		if _, exists := c.hashData[key]; exists {
			delete(c.hashData, key)
		}
		if _, exists := c.listData[key]; exists {
			delete(c.listData, key)
		}
	}

	return newIntCmd(deletedCount, nil)
}

// Exists 检查键是否存在
func (c *RistrettoClient) Exists(ctx context.Context, keys ...string) cache.IntCmd {
	existsCount := int64(0)
	for _, key := range keys {
		if !c.isExpired(key) {
			if _, found := c.cache.Get(key); found {
				existsCount++
			}
		}
	}

	return newIntCmd(existsCount, nil)
}

// Expire 设置键过期时间
func (c *RistrettoClient) Expire(ctx context.Context, key string, ttl time.Duration) cache.BoolCmd {
	value, found := c.cache.Get(key)
	if !found {
		return newBoolCmd(false, nil)
	}

	value.TTL = ttl
	c.cache.SetWithTTL(key, value, cost(value), value.TTL)

	return newBoolCmd(true, nil)
}

// ListCmdable 实现

// LIndex 获取列表指定索引的元素
func (c *RistrettoClient) LIndex(ctx context.Context, key string, index int64) cache.StringCmd {
	if c.isExpired(key) {
		return newStringCmd("", Nil)
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	list, exists := c.listData[key]
	if !exists {
		return newStringCmd("", Nil)
	}

	// 处理负索引
	if index < 0 {
		index = int64(len(list)) + index
	}

	if index < 0 || index >= int64(len(list)) {
		return newStringCmd("", Nil)
	}

	return newStringCmd(list[index], nil)
}

// LPush 从列表左侧推入元素
func (c *RistrettoClient) LPush(ctx context.Context, key string, values ...interface{}) cache.IntCmd {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.listData[key] == nil {
		c.listData[key] = make([]string, 0)
	}

	// 转换值为字符串并添加到列表开头
	newElements := make([]string, 0, len(values))
	for _, value := range values {
		newElements = append(newElements, fmt.Sprintf("%v", value))
	}

	// 在开头插入新元素
	c.listData[key] = append(newElements, c.listData[key]...)

	// 在缓存中标记为list类型
	cacheValue := &CacheValue{
		Data:     c.listData[key],
		Type:     CacheTypeList,
		CreateAt: time.Now(),
	}
	c.cache.Set(key, cacheValue, cost(cacheValue))

	return newIntCmd(int64(len(c.listData[key])), nil)
}

// RPush 从列表右侧推入元素
func (c *RistrettoClient) RPush(ctx context.Context, key string, values ...interface{}) cache.IntCmd {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.listData[key] == nil {
		c.listData[key] = make([]string, 0)
	}

	// 转换值为字符串并添加到列表末尾
	for _, value := range values {
		c.listData[key] = append(c.listData[key], fmt.Sprintf("%v", value))
	}

	// 在缓存中标记为list类型
	cacheValue := &CacheValue{
		Data:     c.listData[key],
		Type:     CacheTypeList,
		CreateAt: time.Now(),
	}
	c.cache.Set(key, cacheValue, cost(cacheValue))

	return newIntCmd(int64(len(c.listData[key])), nil)
}

// LSet 设置列表指定索引的值
func (c *RistrettoClient) LSet(ctx context.Context, key string, index int64, value interface{}) cache.StatusCmd {
	if c.isExpired(key) {
		return newStatusCmd("", fmt.Errorf("no such key"))
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	list, exists := c.listData[key]
	if !exists {
		return newStatusCmd("", fmt.Errorf("no such key"))
	}

	// 处理负索引
	if index < 0 {
		index = int64(len(list)) + index
	}

	if index < 0 || index >= int64(len(list)) {
		return newStatusCmd("", fmt.Errorf("index out of range"))
	}

	c.listData[key][index] = fmt.Sprintf("%v", value)

	// 更新缓存
	cacheValue := &CacheValue{
		Data:     c.listData[key],
		Type:     CacheTypeList,
		CreateAt: time.Now(),
	}
	c.cache.Set(key, cacheValue, cost(cacheValue))

	return newStatusCmd("OK", nil)
}

// LPop 从列表左侧弹出元素
func (c *RistrettoClient) LPop(ctx context.Context, key string) cache.StringCmd {
	if c.isExpired(key) {
		return newStringCmd("", Nil)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	list, exists := c.listData[key]
	if !exists || len(list) == 0 {
		return newStringCmd("", Nil)
	}

	// 弹出第一个元素
	value := list[0]
	c.listData[key] = list[1:]

	// 如果列表为空，删除该键
	if len(c.listData[key]) == 0 {
		delete(c.listData, key)
		c.cache.Del(key)
	} else {
		// 更新缓存
		cacheValue := &CacheValue{
			Data:     c.listData[key],
			Type:     CacheTypeList,
			CreateAt: time.Now(),
		}
		c.cache.Set(key, cacheValue, cost(cacheValue))
	}

	return newStringCmd(value, nil)
}

// LRange 获取列表指定范围的元素
func (c *RistrettoClient) LRange(ctx context.Context, key string, start, stop int64) cache.StringSliceCmd {
	if c.isExpired(key) {
		return newStringSliceCmd([]string{}, nil)
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	list, exists := c.listData[key]
	if !exists {
		return newStringSliceCmd([]string{}, nil)
	}

	listLen := int64(len(list))

	// 处理负索引
	if start < 0 {
		start = listLen + start
	}
	if stop < 0 {
		stop = listLen + stop
	}

	// 边界检查
	if start < 0 {
		start = 0
	}
	if stop >= listLen {
		stop = listLen - 1
	}

	if start > stop || start >= listLen {
		return newStringSliceCmd([]string{}, nil)
	}

	// 创建结果切片
	result := make([]string, stop-start+1)
	copy(result, list[start:stop+1])

	return newStringSliceCmd(result, nil)
}

// RistrettoPipeline Pipeline 实现
type RistrettoPipeline struct {
	client   *RistrettoClient
	commands []PipelineCmd
	mu       sync.Mutex
}

type PipelineCmd struct {
	cmd    func() cache.Cmder
	result cache.Cmder
}

// Set Pipeline methods - 这些方法将命令添加到pipeline而不是立即执行
func (p *RistrettoPipeline) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) cache.StatusCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &statusCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.Set(ctx, key, value, expiration)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) Get(ctx context.Context, key string) cache.StringCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &stringCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.Get(ctx, key)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) IncrBy(ctx context.Context, key string, value int64) cache.IntCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &intCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.IncrBy(ctx, key, value)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) Incr(ctx context.Context, key string) cache.IntCmd {
	return p.IncrBy(ctx, key, 1)
}

func (p *RistrettoPipeline) HSet(ctx context.Context, key string, values ...interface{}) cache.IntCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &intCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.HSet(ctx, key, values...)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) HGetAll(ctx context.Context, key string) cache.MapStringStringCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &mapStringStringCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.HGetAll(ctx, key)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) Del(ctx context.Context, keys ...string) cache.IntCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &intCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.Del(ctx, keys...)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) Exists(ctx context.Context, keys ...string) cache.IntCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &intCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.Exists(ctx, keys...)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) Expire(ctx context.Context, key string, expiration time.Duration) cache.BoolCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &boolCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.Expire(ctx, key, expiration)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) LIndex(ctx context.Context, key string, index int64) cache.StringCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &stringCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.LIndex(ctx, key, index)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) LPush(ctx context.Context, key string, values ...interface{}) cache.IntCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &intCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.LPush(ctx, key, values...)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) RPush(ctx context.Context, key string, values ...interface{}) cache.IntCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &intCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.RPush(ctx, key, values...)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) LSet(ctx context.Context, key string, index int64, value interface{}) cache.StatusCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &statusCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.LSet(ctx, key, index, value)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) LPop(ctx context.Context, key string) cache.StringCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &stringCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.LPop(ctx, key)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) LRange(ctx context.Context, key string, start, stop int64) cache.StringSliceCmd {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := &stringSliceCmd{}
	p.commands = append(p.commands, PipelineCmd{
		cmd: func() cache.Cmder {
			return p.client.LRange(ctx, key, start, stop)
		},
		result: cmd,
	})

	return cmd
}

func (p *RistrettoPipeline) Pipeline() cache.Pipeliner {
	return p // 返回自己，因为pipeline可以嵌套
}

// Exec 执行pipeline中的所有命令
func (p *RistrettoPipeline) Exec(ctx context.Context) ([]cache.Cmder, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	results := make([]cache.Cmder, len(p.commands))

	// 按顺序执行所有命令
	for i, pipeCmd := range p.commands {
		result := pipeCmd.cmd()
		results[i] = result

		// 更新原始命令的结果
		switch cmd := pipeCmd.result.(type) {
		case *statusCmd:
			if statusResult, ok := result.(*statusCmd); ok {
				cmd.val = statusResult.val
				cmd.err = statusResult.err
			}
		case *stringCmd:
			if stringResult, ok := result.(*stringCmd); ok {
				cmd.val = stringResult.val
				cmd.err = stringResult.err
			}
		case *intCmd:
			if intResult, ok := result.(*intCmd); ok {
				cmd.val = intResult.val
				cmd.err = intResult.err
			}
		case *boolCmd:
			if boolResult, ok := result.(*boolCmd); ok {
				cmd.val = boolResult.val
				cmd.err = boolResult.err
			}
		case *mapStringStringCmd:
			if mapResult, ok := result.(*mapStringStringCmd); ok {
				cmd.val = mapResult.val
				cmd.err = mapResult.err
			}
		case *stringSliceCmd:
			if sliceResult, ok := result.(*stringSliceCmd); ok {
				cmd.val = sliceResult.val
				cmd.err = sliceResult.err
			}
		}
	}

	// 清空命令列表
	p.commands = p.commands[:0]

	return results, nil
}

// Close 关闭缓存
func (c *RistrettoClient) Close() {
	c.cache.Close()
}

// cost helps encode any value to a byte buffer
func cost(v interface{}) int64 {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(v)
	if err != nil {
		log.Printf("Failed to encode value: %v", err)
		return 0
	}
	// Return the size of the encoded buffer
	return int64(buf.Len())
}
