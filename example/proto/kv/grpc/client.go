// Code generated by protoc-gen-hcplugin. DO NOT EDIT.
package grpc

import (
	base "github.com/da-moon/protoc-gen-hcplugin/example/proto"

	context "context"

	stacktrace "github.com/palantir/stacktrace"
)

// Client is an implementation of shared.KV that talks over gRPC.
type Client struct {
	client base.KVClient
}

func (c *Client) Get(key string) ([]byte, error) {
	_resp, err := c.client.Get(context.Background(), &base.GetRequest{
		Key: key,
	})
	if err != nil {
		err = stacktrace.Propagate(err, "Get call failed with request %#v", &base.GetRequest{
			Key: key,
		})
	}
	return _resp.Value, err
}

func (c *Client) Put(key string, value []byte) error {
	_, err := c.client.Put(context.Background(), &base.PutRequest{Key: key, Value: value})
	if err != nil {
		err = stacktrace.Propagate(err, "Put call failed with request %#v", &base.PutRequest{
			Key:   key,
			Value: value,
		})
	}
	return err
}
