package module

import (
	"encoding/json"
	"path/filepath"
	"strings"

	proto "github.com/golang/protobuf/proto"

	model "github.com/da-moon/protoc-gen-hcplugin/internal/helper/pgshelper/proto"
	utils "github.com/da-moon/protoc-gen-hcplugin/internal/helper/utils"
	pb "github.com/da-moon/protoc-gen-hcplugin/proto"

	pgshelper "github.com/da-moon/protoc-gen-hcplugin/internal/helper/pgshelper"
	pgs "github.com/lyft/protoc-gen-star"
)

type base struct {
	Package  string
	Imports  map[string]model.Package
	Services []service
}
type service struct {
	model.Service
	pb.Handshake
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
			"grpcx":   "google.golang.org/grpc",
			"plugin":  "github.com/hashicorp/go-plugin",
			"base":    basePath,
			"shared":  filepath.Join(basePath, dir, "shared"),
			"context": "context",
		}, false)
		b.Package = "grpc"

		bb := &base{}
		bb.Package = b.Package
		bb.Imports = b.Imports
		for _, srv := range file.Services() {
			s := model.Service{}
			s.ServiceName = utils.NewString(srv.Name().String())
			for _, method := range srv.Methods() {
				m, err := pgshelper.ExtractMethod(method, p.Context)
				p.CheckErr(err)
				s.Methods = append(s.Methods, *m)

			}
			ss := service{}
			ss.ServiceName = s.ServiceName
			ss.Methods = s.Methods
			opt := srv.Descriptor().GetOptions()
			handShakeopt, err := proto.GetExtension(opt, pb.E_Handshake)
			if err != nil {
				if err == proto.ErrMissingExtension {
					continue
				}
			}
			byteData, err := json.Marshal(handShakeopt)
			if err != nil {
				p.CheckErr(err)
			}
			handshake := pb.Handshake{}
			err = json.Unmarshal(byteData, &handshake)
			if err != nil {
				p.CheckErr(err)
			}
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
			bb.Services = append(bb.Services, ss)
			if len(bb.Services) == 0 {
				continue
			}
		}
		p.OverwriteGeneratorTemplateFile(
			filepath.Join("./"+dir+"/grpc/plugin.go"),
			template.Lookup("Base"),
			&bb,
		)
	}
	return p.Artifacts()
}
