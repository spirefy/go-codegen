package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/spirefy/go-pdk/types"
)

type CommandLineSchema struct {
	Name            string `json:"name"`
	ValueType       string `json:"valueType"`
	MultipleAllowed bool   `json:"multipleAllowed"`
}

func GetExtensions() []types.Extension {
	sourceCommandLine := CommandLineSchema{
		Name:            "source",
		ValueType:       "string",
		MultipleAllowed: true,
	}
	t, _ := json.Marshal(sourceCommandLine)

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
	fmt.Println("Handling a LOAD event")

	return 0
}
