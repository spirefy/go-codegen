package plugins

import (
	"encoding/json"
	"github.com/extism/go-pdk"
	ct "github.com/spirefy/cli/plugin"
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

	options := ct.Options{sourceCommandLine, targetCommandLine}
	t, _ := json.Marshal(options)

	e := types.Extension{
		ExtensionPoint: "spirefy.cli.commandline",
		Description:    "Adds a SORUCE option to the cli command line",
		Name:           "SourceOption",
		Func:           "sourceLoad",
		MetaData:       t,
		Event:          "",
	}

	return []types.Extension{e}
}

//export sourceLoad
func handleSourceLoad() int32 {
	pdk.Log(pdk.LogDebug, "Got Source Load")
	return 0
}
