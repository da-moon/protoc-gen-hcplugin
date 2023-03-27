package module

import (
	"encoding/json"
	"path/filepath"
	"strings"

	proto "github.com/golang/protobuf/proto"

	pgshelper "github.com/da-moon/protoc-gen-hcplugin/internal/helper/pgshelper"
	model "github.com/da-moon/protoc-gen-hcplugin/internal/helper/pgshelper/proto"
	utils "github.com/da-moon/protoc-gen-hcplugin/internal/helper/utils"
	utilspb "github.com/da-moon/protoc-gen-hcplugin/internal/helper/utils/proto"
	pb "github.com/da-moon/protoc-gen-hcplugin/proto"
	pgs "github.com/lyft/protoc-gen-star"
)

type base struct {
	Package  string
	Imports  map[string]model.Package
	Services []service
}
type service struct {
	ServiceName utilspb.String
	pb.Handshake
	Methods []rpc
}
type rpc struct {
	ServiceName utilspb.String
	model.Method
}

// Execute -
// TODO Add a list of built in methods (like size ) and handle collision cases
func (p *plugin) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, file := range targets {
		basePath := p.Context.ImportPath(file).String()
		basePath, err := utils.FromRelativePath(filepath.Join(filepath.Dir(basePath)))
		p.CheckErr(err)
		dir := strings.Split(strings.ReplaceAll(file.Name().String(), "_", "-"), ".")[0]
		b, _ := pgshelper.ExtractBase(file, p.Context, map[string]string{
			"plugin": "github.com/hashicorp/go-plugin",
		}, true)
		b.Package = "shared"
		bb := &base{}
		bb.Package = b.Package
		bb.Imports = b.Imports
		for _, srv := range file.Services() {
			ss := service{}
			ss.ServiceName = utils.NewString(srv.Name().UpperCamelCase().String())
			for _, method := range srv.Methods() {
				m, err := pgshelper.ExtractMethod(method, p.Context)
				p.CheckErr(err)
				mm := rpc{}
				mm.ServiceName = utils.NewString(srv.Name().UpperCamelCase().String())
				mm.InputType = m.InputType
				mm.OutputType = m.OutputType
				mm.Name = m.Name
				mm.InputFields = m.InputFields
				mm.OutputFields = m.OutputFields
				ss.Methods = append(ss.Methods, mm)
			}
			opt := srv.Descriptor().GetOptions()
			if opt != nil {
				handShakeopt, err := proto.GetExtension(opt, pb.E_Handshake)
				if err != nil {
					if err == proto.ErrMissingExtension {
						bb.Services = append(bb.Services, ss)
						continue
					}
				}
				byteData, err := json.Marshal(handShakeopt)
				p.CheckErr(err)
				handshake := pb.Handshake{}
				err = json.Unmarshal(byteData, &handshake)
				p.CheckErr(err)
				ss.ProtocolVersion = handshake.ProtocolVersion
				if handshake.ProtocolVersion == 0 {
					ss.ProtocolVersion = 1
				}
				ss.MagicCookieKey = handshake.MagicCookieKey
				if len(handshake.MagicCookieKey) == 0 {
					ss.MagicCookieKey = "BASIC_PLUGIN"
				}
				ss.MagicCookieValue = handshake.MagicCookieValue
				if len(handshake.MagicCookieValue) == 0 {
					ss.MagicCookieValue = "hello"
				}
			}
			bb.Services = append(bb.Services, ss)
			if len(bb.Services) == 0 {
				continue
			}
		}
		p.OverwriteGeneratorTemplateFile(
			filepath.Join("./"+dir+"/shared/shared.go"),
			template.Lookup("Base"),
			&bb,
		)
	}
	return p.Artifacts()
}
