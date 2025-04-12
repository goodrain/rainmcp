package transport

import (
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

// NewSSEServerTransport 创建一个新的SSE服务器传输
func NewSSEServerTransport(addr string) (transport.ServerTransport, error) {
	return transport.NewSSEServerTransport(addr)
}
