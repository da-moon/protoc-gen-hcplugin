// Code generated by protoc-gen-hcplugin. DO NOT EDIT.
package grpc

import (
	base "github.com/da-moon/protoc-gen-hcplugin/example/proto"

	context "context"

	shared "github.com/da-moon/protoc-gen-hcplugin/example/proto/kv/shared"

	stacktrace "github.com/palantir/stacktrace"
)

// Server - Here is the gRPC server that Client talks to.
type Server struct {
	Impl shared.KVInterface
}

func (s *Server) Get(ctx context.Context, _req *base.GetRequest) (*base.GetResponse, error) {
	value, err := s.Impl.Get(
		_req.Key,
	)
	if err != nil {
		err = stacktrace.Propagate(err, "Get call failed with request %#v", &base.GetRequest{
			Key: _req.Key,
		})
	}
	return &base.GetResponse{
		Value: value,
	}, nil
}

func (s *Server) Put(ctx context.Context, _req *base.PutRequest) (*base.Empty, error) {
	err := s.Impl.Put(
		_req.Key,
		_req.Value,
	)
	if err != nil {
		err = stacktrace.Propagate(err, "Put call failed with request %#v", &base.PutRequest{
			Key:   _req.Key,
			Value: _req.Value,
		})
	}
	return &base.Empty{}, nil
}
