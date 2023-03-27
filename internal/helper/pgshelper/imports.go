package pgshelper

import (
	"path/filepath"
	"strings"

	pb "github.com/da-moon/protoc-gen-hcplugin/internal/helper/pgshelper/proto"
	utils "github.com/da-moon/protoc-gen-hcplugin/internal/helper/utils"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

// ExtractBase -
func ExtractBase(file pgs.File, context pgsgo.Context, packages map[string]string, embedded bool) (*pb.Base, error) {
	result := &pb.Base{}
	imports, err := ExtractImports(file, context, embedded)
	if err != nil {
		return nil, err
	}
	for k, v := range packages {
		_, ok := imports[k]
		if !ok {
			pkg := pb.Package{
				PackageName: k,
				PackagePath: v,
			}
			imports[k] = pkg
		}
	}
	result.Imports = imports
	result.Package = context.PackageName(file).String()
	return result, nil
}

// ExtractImports -
func ExtractImports(file pgs.File, context pgsgo.Context, embedded bool) (map[string]pb.Package, error) {
	result := make(map[string]pb.Package)
	for _, srv := range file.Services() {
		for _, method := range srv.Methods() {
			fields := append(method.Input().Fields(), method.Output().Fields()...)
			for _, field := range fields {
				var err error
				value := ""
				key := ""
				if field.Type().Embed() != nil && embedded {
					value = context.ImportPath(field.Type().Embed()).String()
					key = field.Type().Embed().File().Descriptor().GetOptions().GetGoPackage()
					if len(key) == 0 {
						tmp := strings.Split(field.Type().Embed().File().Descriptor().GetPackage(), ".")
						key = tmp[len(tmp)-1]
					}
					if field.Type().Embed().BuildTarget() {
						value, err = utils.FromRelativePath(filepath.Join(filepath.Dir(value)))
						if err != nil {
							return nil, err
						}
						key = "base"
					}
				}
				if len(key) != 0 {
					pkg := pb.Package{
						PackageName: key,
						PackagePath: value,
					}
					result[key] = pkg
				}
			}
		}
	}
	return result, nil
}
