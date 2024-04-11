package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/extism/go-pdk"
	ct "github.com/spirefy/cli/plugin"
	codegenTypes "github.com/spirefy/go-codegen/types"
	"github.com/spirefy/go-pdk/types"
)

func GetExtensions() []types.Extension {
	sourceCommandLine := ct.Option{
		Name:        "Source",
		Description: "This option adds a source option. It allows for comma separated file or url path names to where source files can be located",
		Option:      "source",
		Type:        "stringOrUrl",
		Default:     ".",
	}

	targetCommandLine := ct.Option{
		Name:        "Target",
		Description: "This option adds a target option. It allows for comma separated list of",
		Option:      "target",
		Type:        "stringOrUrl",
		Default:     ".",
	}

	options := ct.Options{sourceCommandLine}
	t, _ := json.Marshal(options)

	loaderExtension := types.CreateExtension(
		"spirefy.cli.commandLineExtension",
		"SourceOption",
		"spirefy.cli.commandline",
		"Adds a SORUCE option to the cli command line",
		"load",
		t,
		nil)

	options2 := ct.Options{targetCommandLine}
	t2, _ := json.Marshal(options2)

	generatorExtension := types.CreateExtension(
		"spirefy.cli.commandLineExtension",
		"GeneratorOption",
		"spirefy.cli.commandline",
		"Adds a TARGET option to the cli command line",
		"generate",
		t2,
		nil)

	return []types.Extension{*loaderExtension, *generatorExtension}
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

	fmt.Println("Loading srouce: ", src)
	return ret
}

//export load
func handleSourceLoad() int32 {
	input := pdk.Input()

	s := source{}

	err := json.Unmarshal(input, &s)
	if nil != err {
		pdk.Log(pdk.LogDebug, "Error with unmarshal of input: %s"+err.Error())
		return 1
	}

	srcRet := processSource(s)
	var output []byte
	output, err = json.Marshal(srcRet)
	pdk.Output(output)

	pdk.Log(pdk.LogDebug, "Got Source Load done")
	return 0
}

//export generate
func generate() int32 {
	pdk.Log(pdk.LogDebug, "Generating code")
	return 0
}
