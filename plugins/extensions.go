package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/extism/go-pdk"
	ct "github.com/spirefy/cli/plugin"
	codegenTypes "github.com/spirefy/go-codegen/types"
	"github.com/spirefy/go-pdk/types"
)

// GetExtensions
//
// The CLI ExtensionPoint 'spirefy.cli.commandline' requires an extension to provide one or more cli/plugin/Option
// objects to indicate each command line parameter and type of the value. This call returns the list of options for
// the CLI ExtensionPoint and when the function is called will be able to load all sources first, then iterate on all
// generator provided plugins each getting the complete in memory structure built up by all the sources.
func GetExtensions() []types.Extension {
	sourceCommandLine := &ct.Option{
		Name:        "Source",
		Description: "This option adds a source option. It allows for comma separated file or url path names to where source files can be located",
		Option:      "source",
		Type:        "stringOrUrl",
		Default:     "",
	}

	targetCommandLine := &ct.Option{
		Name:        "Target",
		Description: "This option adds a target option. It allows for comma separated list of generator names provided by generator plugins to specify which plugins to use to generate code with",
		Option:      "target",
		Type:        "stringOrUrl",
		Default:     "",
	}

	options := ct.Options{sourceCommandLine, targetCommandLine}
	t, _ := json.Marshal(options)

	loaderExtension := types.CreateExtension(
		"spirefy.cli.commandLineExtensionForCodegen",
		"Codegen Source and Target Options",
		"spirefy.cli.commandline",
		"Adds source and target options to the cli command line",
		"loadAndGenerate",
		t,
		nil)

	return []types.Extension{*loaderExtension}
}

type source struct {
	url        string
	path       string
	sourceType string
}

type sourceRet struct {
	codegenTypes.Resources
	codegenTypes.Components
	codegenTypes.Workflows
}

func processSource(src source) sourceRet {
	ret := sourceRet{
		Resources:  make(codegenTypes.Resources, 0),
		Components: make(codegenTypes.Components, 0),
		Workflows:  make(codegenTypes.Workflows, 0),
	}

	fmt.Println("Loading source: ", src)
	return ret
}

type inputOptions struct {
	Name  string
	Value string
}

//export loadAndGenerate
func loadAndGenerate() int32 {
	input := pdk.Input()

	s := make([]inputOptions, 0)

	err := json.Unmarshal(input, &s)
	if nil != err {
		pdk.Log(pdk.LogDebug, "Error with unmarshal of input: %s"+err.Error())
		return 1
	}

	/*
		srcRet := processSource(s)
		var output []byte
		output, err = json.Marshal(srcRet)
		pdk.Output(output)
	*/
	pdk.Log(pdk.LogDebug, "Got Source Load done")
	return 0
}
