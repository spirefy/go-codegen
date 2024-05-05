package main

import (
	"encoding/json"
	"github.com/extism/go-pdk"
	"github.com/spirefy/go-codegen/types"
	"github.com/spirefy/go-pdk/hostfuncs"
	"strings"
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

func load(sources string) {
	srcs := strings.Split(sources, ",")
	for _, src := range srcs {
		data, err := hostfuncs.LoadFile(src)

		if nil != err {
			pdk.Log(pdk.LogDebug, "Err loading source: "+err.Error())
		}

		if nil != data && len(data) > 0 {
			// now loop through loader extensions, and pass the first 40 or so of bytes of data to them to determine if they
			// can load it or not
			extensions, err := hostfuncs.GetExtensionsForExtensionPoint("spirefy.plugins.codegen.loaders")
			if nil != err {
				pdk.Log(pdk.LogDebug, "Problem getting extensions: "+err.Error())
			} else {
				if nil != extensions && len(extensions) > 0 {
					pdk.Log(pdk.LogDebug, "We got extensions: "+string(len(extensions)))
					for _, ext := range extensions {
						pdk.Log(pdk.LogDebug, "extension: "+ext.Name)
						pdk.Log(pdk.LogDebug, "id: "+ext.Id)
						pdk.Log(pdk.LogDebug, "ep: "+ext.ExtensionPoint)
						pdk.Log(pdk.LogDebug, "func: "+ext.Func)
						extResp, err2 := hostfuncs.CallExtension(ext.Id, nil)
						if nil != err2 {
							pdk.Log(pdk.LogDebug, "error calling extnesion: "+err2.Error())
							pdk.SetError(err2)
						} else {
							pdk.Log(pdk.LogDebug, "We got some response back")
							if nil != extResp && len(extResp) > 0 {
								pdk.Log(pdk.LogDebug, "response is > 0 ")
								resp := types.LoadedResponse{}
								err3 := json.Unmarshal(extResp, &resp)
								if nil != err3 {
									pdk.Log(pdk.LogDebug, "ERROR UNMARSHALLING: "+err3.Error())
									pdk.SetError(err3)
								}
								pdk.Log(pdk.LogDebug, resp.Workflows[0].Id)

							}
						}
					}
				}
			}
		}
	}
}

func generate(targets string) {
	tgts := strings.Split(targets, ",")
	for _, tgt := range tgts {
		pdk.Log(pdk.LogDebug, "Target: "+tgt)
	}
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

	var sources, targets string

	for _, s := range inputs {
		switch s.Name {
		case "sources":
			sources = s.Value
		case "targets":
			targets = s.Value
		}
	}

	if len(sources) > 0 {
		load(sources)
	}

	if len(targets) > 0 {
		generate(targets)
	}

	return 0
}

func main() {}
