package module

import (
	"path/filepath"
	"strings"

	pgshelper "github.com/da-moon/protoc-gen-hcplugin/internal/helper/pgshelper"
	model "github.com/da-moon/protoc-gen-hcplugin/internal/helper/pgshelper/proto"
	utils "github.com/da-moon/protoc-gen-hcplugin/internal/helper/utils"
	pgs "github.com/lyft/protoc-gen-star"
)

// Execute -
// TODO Add a list of built in methods (like size ) and handle collision cases
func (p *plugin) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, file := range targets {
		basePath := p.Context.ImportPath(file).String()
		basePath, err := utils.FromRelativePath(filepath.Join(filepath.Dir(basePath)))
		p.CheckErr(err)
		base, _ := pgshelper.ExtractBase(file, p.Context, map[string]string{
			"stacktrace": "github.com/palantir/stacktrace",
			"base":       basePath,
			"rpc":        "net/rpc",
		}, true)
		for _, srv := range file.Services() {
			s := model.Service{}
			s.ServiceName = utils.NewString(srv.Name().String())
			for _, method := range srv.Methods() {
				m, err := pgshelper.ExtractMethod(method, p.Context)
				p.CheckErr(err)
				s.Methods = append(s.Methods, *m)
			}
			base.Services = append(base.Services, s)
			if len(base.Services) == 0 {
				continue
			}
		}
		dir := strings.Split(strings.ReplaceAll(file.Name().String(), "_", "-"), ".")[0]
		base.Package = "netrpc"
		p.OverwriteGeneratorTemplateFile(
			filepath.Join("./"+dir+"/net-rpc/client.go"),
			template.Lookup("Base"),
			&base,
		)
	}
	return p.Artifacts()
}
