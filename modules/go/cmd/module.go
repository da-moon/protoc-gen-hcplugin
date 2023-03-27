package module

import (
	"encoding/json"
	"path/filepath"
	"strings"

	proto "github.com/golang/protobuf/proto"

	model "github.com/da-moon/protoc-gen-hcplugin/internal/helper/pgshelper/proto"
	utils "github.com/da-moon/protoc-gen-hcplugin/internal/helper/utils"
	utilspb "github.com/da-moon/protoc-gen-hcplugin/internal/helper/utils/proto"
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
	ServiceName utilspb.String
	Methods     []rpc
	pb.CMD
}
type rpc struct {
	ServiceName utilspb.String
	Usage       string
	Aliases     []string
	model.Method
}

// Execute -
// TODO Add a list of built in methods (like size ) and handle collision cases
// remove extension for service
// merge name and module name fields
// merge srv cmd into method cmd
// fix srvice.name -> String
func (p *plugin) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, file := range targets {
		basePath := p.Context.ImportPath(file).String()
		basePath, err := utils.FromRelativePath(filepath.Join(filepath.Dir(basePath)))
		p.CheckErr(err)
		dir := strings.Split(strings.ReplaceAll(file.Name().String(), "_", "-"), ".")[0]
		b, _ := pgshelper.ExtractBase(file, p.Context, map[string]string{
			"grpc":       filepath.Join(basePath, dir, "grpc"),
			"plugin":     "github.com/hashicorp/go-plugin",
			"shared":     filepath.Join(basePath, dir, "shared"),
			"cli":        "github.com/urfave/cli",
			"stacktrace": "github.com/palantir/stacktrace",
			"handler":    filepath.Join(basePath, dir, "handler", "cmd"),
			"rpc":        filepath.Join(basePath, dir, "net-rpc"),
			"fmt":        "fmt",
			"exec":       "os/exec",
		}, false)
		b.Package = "cmd"

		bb := &base{}
		bb.Package = b.Package
		bb.Imports = b.Imports
		for _, srv := range file.Services() {
			ss := service{}
			ss.ServiceName = utils.NewString(srv.Name().UpperCamelCase().String())
			for _, method := range srv.Methods() {
				m, err := pgshelper.ExtractMethod(method, p.Context)
				opt := method.Descriptor().GetOptions()
				p.CheckErr(err)
				mm := rpc{}
				mm.ServiceName = utils.NewString(srv.Name().UpperCamelCase().String())
				mm.InputType = m.InputType
				mm.OutputType = m.OutputType
				mm.Name = m.Name
				mm.InputFields = m.InputFields
				mm.OutputFields = m.OutputFields
				if opt != nil {
					opts, err := proto.GetExtension(opt, pb.E_CmdMethodOptions)
					if err != nil {
						if err == proto.ErrMissingExtension {
							ss.Methods = append(ss.Methods, mm)
							continue
						}
					}
					byteData, err := json.Marshal(opts)
					p.CheckErr(err)
					p.CheckErr(err)
					cmd := pb.CMD{}
					err = json.Unmarshal(byteData, &cmd)
					p.CheckErr(err)

					mm.Name = m.Name
					if len(cmd.Name) != 0 {
						mm.Name = utils.NewString(cmd.Name)
					}
					mm.Usage = cmd.Usage
					if len(cmd.Usage) == 0 {
						mm.Usage = strings.ReplaceAll(mm.Name.LowerSnakeCase, "_", "-") + " Command"
					}
					mm.Aliases = cmd.Aliases
					if len(cmd.Aliases) == 0 {
						mm.Aliases = append(cmd.Aliases, strings.ReplaceAll(mm.Name.LowerSnakeCase, "_", "-"))
					}
				}
				ss.Methods = append(ss.Methods, mm)
			}
			opt := srv.Descriptor().GetOptions()
			if opt != nil {
				opts, err := proto.GetExtension(opt, pb.E_CmdServiceOptions)
				if err != nil {
					if err == proto.ErrMissingExtension {
						bb.Services = append(bb.Services, ss)
						continue
					}
				}
				byteData, err := json.Marshal(opts)
				p.CheckErr(err)
				cmd := pb.CMD{}
				err = json.Unmarshal(byteData, &cmd)
				p.CheckErr(err)
				ss.Name = ss.ServiceName.UpperCamelCase
				if len(cmd.Name) != 0 {
					ss.Name = cmd.Name
				}
				ss.Usage = cmd.Usage
				if len(cmd.Usage) == 0 {
					ss.Usage = strings.ReplaceAll(pgs.Name(ss.Name).LowerSnakeCase().String(), "_", "-") + " Engine"
				}
				ss.Aliases = cmd.Aliases
				if len(cmd.Aliases) == 0 {
					ss.Aliases = append(cmd.Aliases, strings.ReplaceAll(pgs.Name(ss.Name).LowerSnakeCase().String(), "_", "-"))
				}
			}
			bb.Services = append(bb.Services, ss)
			if len(bb.Services) == 0 {
				continue
			}
		}
		p.OverwriteGeneratorTemplateFile(
			filepath.Join("./"+dir+"/cmd/"+strings.Split(file.Name().String(), ".")[0]+".go"),
			template.Lookup("Base"),
			&bb,
		)
	}
	return p.Artifacts()
}
