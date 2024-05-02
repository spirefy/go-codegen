package main

import (
	"encoding/json"
	"github.com/extism/go-pdk"
)

//export start
func start() int32 {
	pdk.Log(pdk.LogDebug, "Codegen Plugin was called start()")
	return 0
}

type loadinput struct {
	Name  string
	Value string
}

//export loadAndGenerate
func loadAndGenerate() int32 {
	pdk.Log(pdk.LogDebug, "Loading and generating code")
	input := pdk.Input()
	inputs := make([]loadinput, 0)
	err := json.Unmarshal(input, &inputs)
	if nil != err {
		pdk.Log(pdk.LogDebug, "Problem unmarshalling to slice of struct")
	}

	for _, s := range inputs {
		pdk.Log(pdk.LogDebug, "Name is "+s.Name)
		pdk.Log(pdk.LogDebug, "Value is "+s.Value)
	}

	return 0
}

func main() {}
