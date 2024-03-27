package main

import (
	"github.com/spirefy/go-codegen/plugins"
	hostfuncs "github.com/spirefy/go-pdk"
)

//export start
func start() int32 {
	hostfuncs.RegisterPlugin("spirefy.codegen", "Spirefy Codegen", "1.0.0", "1.0.0", "A plugin that provides code generation capabilities with additional extension points for other plugins to contribute to", plugins.GetCodegenExtensionPoints(), plugins.GetExtensions())
	return 0
}

func main() {}
