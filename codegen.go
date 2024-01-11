package main

import (
	"encoding/json"
	"errors"
	"github.com/extism/go-pdk"
	t "github.com/spirefy/go-plugin-engine/types"
)

//export register
func register() int32 {
	// input := pdk.Input()
	// This extension point expects a Menu type to be provided by an extension.
	// It will call each extension once to obtain the name, parent name and onclick func name to call
	ep := t.ExtensionPoint{
		Id:          "codegen.plugins.loaders",
		Description: "This extension point provides the functionality for extensions to provide source format loaders that will add to this plugins unified in memory structure that is then passed to generators",
		Name:        "Codegen Source Loader ExtensionPoint",
	}

	p := &t.Plugin{
		Id:              "spirefy.codegen",
		Name:            "Spirefy Codegen",
		Version:         "1.0.0",
		Description:     "A plugin that provides code generation capabilities with additional extension points for other plugins to contribute to",
		ExtensionPoints: []t.ExtensionPoint{ep},
	}

	o, e := json.Marshal(p)
	if nil != e {
		return 1
	} else {
		pdk.Output(o)
	}

	return 0
}

func generate() {
	pdk.Log(pdk.LogDebug, "Generating code")
}

//export handleEvent
func handleEvent() int32 {
	input := pdk.Input()
	evt := t.Event{}
	err := json.Unmarshal(input, &evt)

	if nil != err {
		pdk.SetError(errors.New("problem unmarshalling the event"))
		return 1
	}

	switch evt.Id {
	case "generate":
		break
	}
	return 0
}

func main() {}
