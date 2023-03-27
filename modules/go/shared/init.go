package module

import (
	"io"
	"strings"
	templatex "text/template"

	"github.com/gobuffalo/packr/v2"
)

var template *templatex.Template

func init() {
	template = templatex.New("shared_module").Funcs(templatex.FuncMap{
		"inc": func(input, value int) int {
			return input + value
		},
		"dec": func(input, value int) int {
			return input - value
		},
		"builtin": func(arg string) bool {
			switch arg {
			case "size":
				return true
			case "Size":
				return true
			case "type":
				return true
			default:
				return false

			}
		},
	})
	box := packr.New("shared_module_box", "./template")
	err := box.Walk(walkFN)
	if err != nil {
		panic(err)
	}
}

var walkFN = func(s string, file packr.File) error {
	var sb strings.Builder
	if _, err := io.Copy(&sb, file); err != nil {
		return err
	}
	var err error
	template, err = template.Parse(sb.String())
	return err
}
