package main

import (
	"encoding/json"
	"github.com/extism/go-pdk"
	"github.com/spirefy/go-codegen/plugins"
	hostfuncs "github.com/spirefy/go-pdk"
	t "github.com/spirefy/go-pdk/types"
)

//export start
func start() int32 {
	pdk.Log(pdk.LogDebug, "STARTING THE CODEGEN PLUGIN")

	hostfuncs.RegisterPlugin("spirefy.codegen", "Spirefy Codegen", "1.0.0", "1.0.0", "A plugin that provides code generation capabilities with additional extension points for other plugins to contribute to", plugins.GetCodegenExtensionPoints(), plugins.GetExtensions())

	return 0
}

//export load
func load() int32 {
	e := make([]t.Extension, 0)
	dta := pdk.Input()
	if nil != dta && len(dta) > 0 {
		err := json.Unmarshal(pdk.Input(), &e)

		if nil != err {
			pdk.SetError(err)
			return 1
		}
	}

	pdk.Log(pdk.LogInfo, "Loading source for extensions")

	for _, e := range e {
		pdk.Log(pdk.LogDebug, "Loading extension: "+e.Name)
	}

	return 0
}

func main() {}
