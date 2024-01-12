package main

import (
	"encoding/json"
	"errors"
	"github.com/extism/go-pdk"
	hostfuncs "github.com/spirefy/go-pdk"
	t "github.com/spirefy/go-pdk/types"
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

	l := t.Listener{
		Event:    "load",
		FuncName: "load",
	}

	l2 := t.Listener{
		Event:    "load",
		FuncName: "loadme",
	}

	p := &t.Plugin{
		Id:              "spirefy.codegen",
		Name:            "Spirefy Codegen",
		Version:         "1.0.0",
		Description:     "A plugin that provides code generation capabilities with additional extension points for other plugins to contribute to",
		ExtensionPoints: []t.ExtensionPoint{ep},
		Listeners:       []t.Listener{l, l2},
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

//export loadme
func loadme() int32 {
	pdk.Log(pdk.LogInfo, "Loading source for loadme extension event")

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
