package pgshelper

import (
	"path/filepath"
	"strings"

	pb "github.com/da-moon/protoc-gen-hcplugin/internal/helper/pgshelper/proto"
	utils "github.com/da-moon/protoc-gen-hcplugin/internal/helper/utils"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

// ExtractMethod -
func ExtractMethod(method pgs.Method, context pgsgo.Context) (*pb.Method, error) {
	result := &pb.Method{}
	result.InputType = context.Name(method.Input()).String()
	result.OutputType = context.Name(method.Output()).String()
	result.Name = utils.NewString(context.Name(method).String())
	for _, field := range method.Input().Fields() {
		e, err := ExtractField(field, context)
		if err != nil {
			return nil, err
		}
		if e != nil {
			result.InputFields = append(result.InputFields, *e)
		}
	}

	for _, field := range method.Output().Fields() {
		e, err := ExtractField(field, context)
		if err != nil {
			return nil, err
		}
		if e != nil {
			result.OutputFields = append(result.OutputFields, *e)
		}
	}
	return result, nil
}

// ExtractField -
func ExtractField(field pgs.Field, context pgsgo.Context) (*pb.Field, error) {
	result := &pb.Field{}
	result.GoGoFieldOptions = &pb.GoGoFieldOptions{}
	result.GoGoFieldOptions.Nullable = true

	var err error
	rPath := ""
	fPackage := ""
	vName := context.Name(field).String()
	vType := context.Type(field).String()
	vPrefix := ""
	if field.Descriptor().GetOptions() != nil {
		opts := field.Descriptor().GetOptions()

		extensions := strings.Split(opts.String(), " ")
		for _, v := range extensions {
			v = strings.TrimSpace(strings.ReplaceAll(v, `"`, ""))
			if len(v) != 0 {
				splitted := strings.Split(v, ":")
				// Custom Name
				if splitted[0] == "65004" {
					result.GoGoFieldOptions.CustomName = splitted[1]
				}
				// Nullable
				if splitted[0] == "65001" {
					if splitted[1] == "0" {

						result.GoGoFieldOptions.Nullable = false
					} else {
						result.GoGoFieldOptions.Nullable = true
					}
				}
				// custom type
				if splitted[0] == "65003" {
					result.GoGoFieldOptions.CustomType = splitted[1]
				}
				if splitted[0] == "65010" {
					if splitted[1] == "0" {
						result.GoGoFieldOptions.STDTime = false
					} else {
						result.GoGoFieldOptions.STDTime = true
					}
				}
				if splitted[0] == "65011" {
					if splitted[1] == "0" {
						result.GoGoFieldOptions.STDDuration = false
					} else {
						result.GoGoFieldOptions.STDDuration = true
					}
				}
				if splitted[0] == "62023" {
					result.GoGoFieldOptions.EnumCustomName = splitted[1]
				}
			}
		}
	}
	if result.GoGoFieldOptions != nil && len(result.GoGoFieldOptions.CustomName) != 0 {
		vName = result.GoGoFieldOptions.GetCustomName()
	}
	if result.GoGoFieldOptions != nil && len(result.GoGoFieldOptions.CustomType) != 0 {
		vType = result.GoGoFieldOptions.GetCustomType()
	}
	if result.GoGoFieldOptions != nil && !result.GoGoFieldOptions.Nullable {
		if strings.ContainsAny(vType, "*") {
			vPrefix = ""
			vType = strings.ReplaceAll(vType, `*`, "")
			vType = strings.TrimSpace(vType)

		}
	}
	if strings.ContainsAny(vType, "*") {
		vPrefix = "*"
		vType = strings.ReplaceAll(vType, `*`, "")
		vType = strings.TrimSpace(vType)
	}
	if field.Type().Embed() != nil {

		rPath = context.ImportPath(field.Type().Embed()).String()

		fPackage = field.Type().Embed().File().Descriptor().GetOptions().GetGoPackage()
		if len(fPackage) == 0 {
			tmp := strings.Split(field.Type().Embed().File().Descriptor().GetPackage(), ".")
			fPackage = tmp[len(tmp)-1]
		}
		if field.Type().Embed().BuildTarget() {
			rPath, err = utils.FromRelativePath(filepath.Join(filepath.Dir(rPath)))
			if err != nil {
				return nil, err
			}
			fPackage = "base"
		}
	}

	result.VariableName = utils.NewKV(vPrefix, vName)
	splitted := strings.Split(vType, ".")
	if len(splitted) > 1 {
		vType = splitted[1]
	}
	result.VariableType = utils.NewKV(fPackage, vType)
	result.Package = utils.NewKV(fPackage, rPath)
	return result, nil
}
