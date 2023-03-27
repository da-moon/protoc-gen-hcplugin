package module

import (
	"path/filepath"
	"strings"

	model "github.com/da-moon/protoc-gen-hcplugin/internal/helper/pgshelper/proto"
	utils "github.com/da-moon/protoc-gen-hcplugin/internal/helper/utils"

	pgshelper "github.com/da-moon/protoc-gen-hcplugin/internal/helper/pgshelper"
	pgs "github.com/lyft/protoc-gen-star"
)

// Execute -
// TODO Add a list of built in methods (like size ) and handle collision cases
func (p *plugin) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, file := range targets {
		basePath := p.Context.ImportPath(file).String()
		basePath, err := utils.FromRelativePath(filepath.Join(filepath.Dir(basePath)))
		p.CheckErr(err)
		dir := strings.Split(strings.ReplaceAll(file.Name().String(), "_", "-"), ".")[0]

		base, _ := pgshelper.ExtractBase(file, p.Context, map[string]string{
			"stacktrace": "github.com/palantir/stacktrace",
			"base":       basePath,
			"shared":     filepath.Join(basePath, dir, "shared"),
			"context":    "context",
		}, false)
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
		base.Package = "grpc"

		p.OverwriteGeneratorTemplateFile(
			filepath.Join("./"+dir+"/grpc/server.go"),
			template.Lookup("Base"),
			&base,
		)
	}
	return p.Artifacts()
}
