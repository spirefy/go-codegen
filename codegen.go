package main

import (
	"encoding/json"
	"github.com/extism/go-pdk"
	"github.com/spirefy/go-codegen/plugins"
	hostfuncs "github.com/spirefy/go-pdk"
	t "github.com/spirefy/go-pdk/types"
)

//export register
func register() int32 {
	p := &t.Plugin{
		Id:              "spirefy.codegen",
		Name:            "Spirefy Codegen",
		Version:         "1.0.0",
		Description:     "A plugin that provides code generation capabilities with additional extension points for other plugins to contribute to",
		ExtensionPoints: plugins.GetCodegenExtensionPoints(),
		Extensions:      plugins.GetExtensions(),
	}

	o, e := json.Marshal(p)
	if nil != e {
		return 1
	} else {
		pdk.Output(o)
	}

	hostfuncs.CallExtension([]byte{})
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
