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
		dir := strings.Split(strings.ReplaceAll(file.Name().String(), "_", "-"), ".")[0]
		b, err := pgshelper.ExtractBase(file, p.Context, map[string]string{
			"shared": filepath.Join(basePath, dir, "shared"),
			"cli":    "github.com/urfave/cli",
		}, false)

		p.CheckErr(err)
		b.Package = "cmd"
		for _, srv := range file.Services() {
			s := model.Service{}
			s.ServiceName = utils.NewString(srv.Name().String())
			for _, method := range srv.Methods() {
				m, err := pgshelper.ExtractMethod(method, p.Context)
				p.CheckErr(err)
				s.Methods = append(s.Methods, *m)
			}
			b.Services = append(b.Services, s)
			if len(b.Services) == 0 {
				continue
			}
		}
		p.OverwriteGeneratorTemplateFile(
			filepath.Join("./"+dir+"/handler/cmd/action.go"),
			template.Lookup("Base"),
			&b,
		)
	}
	return p.Artifacts()
}
