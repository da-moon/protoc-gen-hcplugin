package main

import (
	cmdmodule "github.com/da-moon/protoc-gen-hcplugin/modules/go/cmd"
	grpcclient "github.com/da-moon/protoc-gen-hcplugin/modules/go/grpc/client"
	grpcplugin "github.com/da-moon/protoc-gen-hcplugin/modules/go/grpc/plugin"
	grpcserver "github.com/da-moon/protoc-gen-hcplugin/modules/go/grpc/server"
	cmdhandlermodule "github.com/da-moon/protoc-gen-hcplugin/modules/go/handler/cmd"
	enginehandlermodule "github.com/da-moon/protoc-gen-hcplugin/modules/go/handler/engine"
	netrpcclient "github.com/da-moon/protoc-gen-hcplugin/modules/go/net-rpc/client"
	netrpcplugin "github.com/da-moon/protoc-gen-hcplugin/modules/go/net-rpc/plugin"
	netrpcserver "github.com/da-moon/protoc-gen-hcplugin/modules/go/net-rpc/server"
	sharedmodule "github.com/da-moon/protoc-gen-hcplugin/modules/go/shared"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func main() {
	plugin := pgs.Init(
		pgs.DebugEnv("DEBUG"),
	)
	// Registering shared/grpc module
	plugin.RegisterModule(
		grpcclient.New(),
		grpcserver.New(),
		grpcplugin.New(),
		netrpcplugin.New(),
		netrpcclient.New(),
		netrpcserver.New(),
		sharedmodule.New(),
		cmdmodule.New(),
		enginehandlermodule.New(),
		cmdhandlermodule.New(),
	)
	// TODO : This may cause issues when other language targets are added
	plugin.RegisterPostProcessor(pgsgo.GoFmt()).Render()
}
