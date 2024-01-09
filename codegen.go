package main

import (
	"encoding/json"
	"github.com/extism/go-pdk"
	t "github.com/spirefy/go-plugin-engine/types"
)

// register
//
// This is an internal function that returns the plugineigne's t.Plugin object to register this plugin along with
// extension points and extensions if any
func register() *t.Plugin {
	// This extension point expects a Menu type to be provided by an extension.
	// It will call each extension once to obtain the name, parent name and onclick func name to call
	ep := t.ExtensionPoint{
		Id:          "codegen.plugins.loaders",
		Description: "This extension point provides the functionality for extensions to provide source format loaders that will add to this plugins unified in memory structure that is then passed to generators",
		Name:        "Codegen Source Loader ExtensionPoint",
		StartOnLoad: true,
		Schema:      t.Schema{},
	}

	ret := &t.Plugin{
		Id:              "spirefy.codegen",
		Name:            "Spirefy Codegen",
		Version:         "1.0.0",
		Description:     "A plugin that provides code generation capabilities with additional extension points for other plugins to contribute to",
		ExtensionPoints: []t.ExtensionPoint{ep},
	}

	return ret
}

//export pluginInit
func pinit() int32 {
	// input := pdk.Input()

	p := register()

	o, e := json.Marshal(p)
	if nil != e {
		return 1
	} else {
		pdk.Output(o)
	}

	return 0
}

func main() {}
