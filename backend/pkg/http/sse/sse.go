package sse

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Event SSE 事件结构
type Event struct {
	Event string      `json:"event,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	ID    string      `json:"id,omitempty"`
	Retry int         `json:"retry,omitempty"`
}

// SSESender SSE发送器接口
type SSESender interface {
	Send(ctx context.Context, event *Event) error
	Close() error
}

// SSenderImpl Gin SSE发送器实现
type SSenderImpl struct {
	c      *gin.Context
	writer gin.ResponseWriter
	closed bool
}

// NewSSESender 创建新的SSE发送器
func NewSSESender(c *gin.Context) *SSenderImpl {
	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")

	return &SSenderImpl{
		c:      c,
		writer: c.Writer,
		closed: false,
	}
}

// Send 发送SSE事件
func (s *SSenderImpl) Send(ctx context.Context, event *Event) error {
	if s.closed {
		return fmt.Errorf("SSE connection is closed")
	}

	// 检查上下文是否取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 构建SSE数据格式
	if event.ID != "" {
		s.writer.WriteString(fmt.Sprintf("id: %s\n", event.ID))
	}

	if event.Event != "" {
		s.writer.WriteString(fmt.Sprintf("event: %s\n", event.Event))
	}

	if event.Retry > 0 {
		s.writer.WriteString(fmt.Sprintf("retry: %d\n", event.Retry))
	}

	// 写入数据
	switch data := event.Data.(type) {
	case string:
		s.writer.WriteString(fmt.Sprintf("data: %s\n\n", data))
	case []byte:
		s.writer.WriteString(fmt.Sprintf("data: %s\n\n", string(data)))
	default:
		// 对于其他类型，使用JSON序列化
		s.c.SSEvent("data", data)
		s.writer.WriteString("\n")
	}

	// 刷新缓冲区
	s.writer.Flush()

	return nil
}

// SendString 发送字符串数据的便捷方法
func (s *SSenderImpl) SendString(ctx context.Context, eventType, data string) error {
	return s.Send(ctx, &Event{
		Event: eventType,
		Data:  data,
	})
}

// SendJSON 发送JSON数据的便捷方法
func (s *SSenderImpl) SendJSON(ctx context.Context, eventType string, data interface{}) error {
	return s.Send(ctx, &Event{
		Event: eventType,
		Data:  data,
	})
}

// SendWithID 发送带ID的事件
func (s *SSenderImpl) SendWithID(ctx context.Context, id, eventType string, data interface{}) error {
	return s.Send(ctx, &Event{
		ID:    id,
		Event: eventType,
		Data:  data,
	})
}

// Close 关闭SSE连接
func (s *SSenderImpl) Close() error {
	if s.closed {
		return nil
	}

	s.closed = true
	// 发送关闭事件
	s.writer.WriteString("event: close\ndata: \n\n")
	s.writer.Flush()

	return nil
}

// IsClosed 检查连接是否已关闭
func (s *SSenderImpl) IsClosed() bool {
	return s.closed
}

// Stream SSE流管理器
type Stream struct {
	clients map[string]*SSenderImpl

	// 用于添加/删除客户端的通道
	addClient    chan *Client
	removeClient chan *Client

	// 广播消息通道
	broadcast chan *Event

	// 停止信号
	stopCh chan struct{}
}

// Client 客户端信息
type Client struct {
	ID     string
	Sender *SSenderImpl
	Events chan *Event
}

// NewStream 创建新的SSE流
func NewStream() *Stream {
	return &Stream{
		clients:      make(map[string]*SSenderImpl),
		addClient:    make(chan *Client),
		removeClient: make(chan *Client),
		broadcast:    make(chan *Event, 100),
		stopCh:       make(chan struct{}),
	}
}

// Start 启动流管理器
func (s *Stream) Start() {
	go s.run()
}

// Stop 停止流管理器
func (s *Stream) Stop() {
	close(s.stopCh)
}

// AddClient 添加客户端
func (s *Stream) AddClient(client *Client) {
	s.addClient <- client
}

// RemoveClient 移除客户端
func (s *Stream) RemoveClient(client *Client) {
	s.removeClient <- client
}

// Publish 广播消息到所有客户端
func (s *Stream) Publish(event *Event) {
	select {
	case s.broadcast <- event:
	default:
		// 广播通道满了，丢弃消息
	}
}

// run 运行流管理器
func (s *Stream) run() {
	for {
		select {
		case client := <-s.addClient:
			s.clients[client.ID] = client.Sender

		case client := <-s.removeClient:
			if _, exists := s.clients[client.ID]; exists {
				delete(s.clients, client.ID)
				client.Sender.Close()
			}

		case event := <-s.broadcast:
			// 广播到所有客户端
			for id, sender := range s.clients {
				if err := sender.Send(context.Background(), event); err != nil {
					// 发送失败，移除客户端
					delete(s.clients, id)
					sender.Close()
				}
			}

		case <-s.stopCh:
			// 关闭所有客户端连接
			for _, sender := range s.clients {
				sender.Close()
			}
			return
		}
	}
}

// SSEHandler 创建SSE处理器的便捷函数
func SSEHandler(stream *Stream) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端ID（可以从查询参数或头部获取）
		clientID := c.Query("client_id")
		if clientID == "" {
			clientID = c.GetHeader("X-Client-ID")
		}
		if clientID == "" {
			// 生成默认ID
			clientID = fmt.Sprintf("client_%d", generateClientID())
		}

		// 创建SSE发送器
		sender := NewSSESender(c)

		// 创建客户端
		client := &Client{
			ID:     clientID,
			Sender: sender,
			Events: make(chan *Event, 10),
		}

		// 添加到流中
		stream.AddClient(client)

		// 发送连接成功消息
		sender.SendString(c.Request.Context(), "connected", "Connection established")

		// 保持连接活跃
		ctx := c.Request.Context()
		for {
			select {
			case <-ctx.Done():
				// 客户端断开连接
				stream.RemoveClient(client)
				return
			case event := <-client.Events:
				// 发送特定于客户端的事件
				if err := sender.Send(ctx, event); err != nil {
					stream.RemoveClient(client)
					return
				}
			}
		}
	}
}

// generateClientID 生成客户端ID的简单函数
func generateClientID() int64 {
	// 简单的时间戳ID，实际使用中可以使用UUID
	return 1000000 // 这里应该使用实际的ID生成逻辑
}

// 使用示例函数
func ExampleUsage() {
	r := gin.Default()

	// 创建SSE流
	sseStream := NewStream()
	sseStream.Start()

	// SSE端点
	r.GET("/events", SSEHandler(sseStream))

	// 触发广播的端点
	r.POST("/broadcast", func(c *gin.Context) {
		var req struct {
			Event string      `json:"event"`
			Data  interface{} `json:"data"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 广播事件
		sseStream.Publish(&Event{
			Event: req.Event,
			Data:  req.Data,
		})

		c.JSON(200, gin.H{"message": "Event broadcasted"})
	})

	// 启动服务器
	r.Run(":8080")
}
