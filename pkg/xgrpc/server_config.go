// Package xgrpc ..
// grpc服务端配置信息
// 提供Build函数，构建grpc服务端
package xgrpc

import (
	"strings"

	"google.golang.org/grpc"
)

// ConfigHandler ..
type ConfigHandler interface {
	GetString(key string) string
	GetInt(key string) int
}

// Register ..
type Register interface {
	RegistryService(key, value string)
	UnRegistryService(key string)
}

// DefaultServerConfig ..
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Network: "tcp",
		Addr:    "127.0.0.1:3000",
		Scheme:  "schema",
		Name:    "application",
	}
}

// RawServerConfig ..
func RawServerConfig(confPrefix string, confHandler ConfigHandler) *ServerConfig {
	if strings.HasSuffix(confPrefix, ".") {
		confPrefix = confPrefix[:len(confPrefix)-1]
	}
	return &ServerConfig{
		Network: confHandler.GetString(confPrefix + ".network"),
		Addr:    confHandler.GetString(confPrefix + ".addr"),
		Scheme:  confHandler.GetString(confPrefix + ".scheme"),
		Name:    confHandler.GetString(confPrefix + ".name"),
	}
}

// ServerConfig ..
type ServerConfig struct {
	Addr    string
	Network string
	Scheme  string
	Name    string

	fns      []func(*grpc.Server)
	register Register
}

// Setup ..
func (t *ServerConfig) Setup(fns ...func(*grpc.Server)) *ServerConfig {
	t.fns = fns
	return t
}

// SetRegistry ..
func (t *ServerConfig) WithRegister(register Register) *ServerConfig {
	t.register = register
	return t
}

// Build ..
func (t *ServerConfig) Build() *GrpcServer {
	return newServer(t)
}