// Code generated by protoc-gen-hcplugin. DO NOT EDIT.
package shared

import (
	plugin "github.com/hashicorp/go-plugin"
)

// KVInterface - this is the interface that we're exposing as a plugin.
type KVInterface interface {
	Get(key string) ([]byte, error)

	Put(key string, value []byte) error
}

// HandshakeConfig - engine-interface handshake configuration
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}