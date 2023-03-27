package module

import (
	pgs "github.com/lyft/protoc-gen-star"

	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

// template -
type plugin struct {
	pgs.ModuleBase
	pgsgo.Context
}

// New -
func New() *plugin {
	return &plugin{
		ModuleBase: pgs.ModuleBase{},
	}
}

// Name -
func (*plugin) Name() string {
	return "hcplugin_gRPC_client_module"
}

func (p *plugin) InitContext(c pgs.BuildContext) {
	p.ModuleBase.InitContext(c)
	p.Context = pgsgo.InitContext(c.Parameters())
}
