package utils

import (
	"strings"

	pgs "github.com/lyft/protoc-gen-star"

	pb "github.com/da-moon/protoc-gen-hcplugin/internal/helper/utils/proto"
)

// NewString - Helpper functio to create a new wrapped string
func NewString(value string) pb.String {
	result := pb.String{}
	result.Original = value
	result.LowerCamelCase = pgs.Name(value).LowerCamelCase().String()
	result.UpperCamelCase = pgs.Name(value).UpperCamelCase().String()
	result.LowerSnakeCase = pgs.Name(value).LowerSnakeCase().String()
	result.UpperSnakeCase = pgs.Name(value).UpperSnakeCase().String()
	result.LowerDotCase = pgs.Name(value).LowerDotNotation().String()
	result.UpperDotCase = pgs.Name(value).UpperDotNotation().String()
	result.LowerParamCase = strings.ReplaceAll(pgs.Name(value).LowerSnakeCase().String(), "_", "-")
	result.UpperParamCase = strings.ReplaceAll(pgs.Name(value).UpperSnakeCase().String(), "_", "-")
	return result
}

// NewKV -
func NewKV(key, value string) pb.KV {
	result := pb.KV{}
	result.Key = NewString(key)
	result.Value = NewString(value)
	return result
}
