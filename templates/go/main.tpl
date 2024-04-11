package {{.Package}}

import (
	"encoding/json"
	// extism Go PDK is required for the register exported function to use the pdk.Input() and pdk.Output() calls.
	"github.com/extism/go-pdk"
	// hostFuncs below should be uncommented if this plugin will need to make use of host funcs provided to plugins
	// see more on how to use this capability in the plugin engine documentation
	// hostFuncs "{{.GoPluginPDK}}"
	pluginTypes "{{.GoPluginPDK}}/types"
)

func GetExtensionPoints() []pluginTypes.ExtensionPoint {
	// This extension point expects a Menu type to be provided by an extension.
	// It will call each extension once to obtain the name, parent name and onclick func name to call
	ep := pluginTypes.ExtensionPoint{
		Id:          "... fill in extension point ID here ... ",
		Description: "... fill in a description of this extension point ...",
		Name:        "... fill in a name that could be displayed and readable ...",
	}

	return []pluginTypes.ExtensionPoint{ep}
}

func GetExtensions() []pluginTypes.Extension {
    // Here you would do any code necessary to for whatever extension point the extension here would contribute to.
    // This could be setting up the required extension point structure which the extension point should provide
    // documentation for on the structure needed.
    // The pluginTypes.Extension structure is then filled in with the MetaData property used to carry whatever
    // Extension Point data is required in a []byte format.. so in Go, using json.Marshal() would be necessary to
    // convert the struct type for example into the []byte data.
    //
    // Example:
    //
	// extensionPointRequiredType := struct {
	//   SomeStringValue   string
	//   SomeBoolValue     bool
	// }{
	//   SomeStringValue: "value",
	//   SomeBoolValue:   true,
	// }
	//
	// ext := pluginTypes.Extension {
	//   ExtensionPoint: "... fill in extension point ID that this extension contributes to ...",
	//   Description: "... fill in an optional meaningful description ...",
	//   Name: "... fill in a name ideally a unique one similar to the extension point ID suggested format ...",
	//   Func: "... fill in the name of the function that will be exported from this module that can be called by the Extism Host Function... ",
	//   MetaData: "... provide the []byte object, ideally a json marshalled struct.. ",
	// }

	// This is a simplified example using the above struct in the example above. The below is not complete and would
	// need the Func and MetaData values filled in as per the info above.

	e := pluginTypes.Extension{
		ExtensionPoint: "... fill in extension point ID that this extension contributes to ...",
		Description:    "... fill in an optional meaningful description ...",
		Name:           "... fill in a name ideally a unique one similar to the extension point ID suggested format ...",
	}

	return []pluginTypes.Extension{e}
}

// register
// This is a required function by the plugin engine. It must be present in all plugins. It returns the Plugin structure
// tot he Plugin Engine loading this plugin to register the id, name, version and option description, extension points
// and extensions
//
// The export is required by the tinygo compiler for the WASM/WASI target output.
//
//export register
func register() int32 {
	plg := &pluginTypes.Plugin{
		Id:              "{{.PluginId}}",
		Name:            "{{.PluginName}}",
		Version:         "{{.PluginVersion}}",
		Description:     "A plugin that provides code generation capabilities with additional extension points for other plugins to contribute to",
		ExtensionPoints: GetExtensionPoints(),
		Extensions:      GetExtensions(),
	}

    // Marshall the registration plugin structure
	p, e := json.Marshal(plg)
	if nil != e {
	    // Extism SDK requires a 1 for error status
		return 1
	} else {
	    // send the []byte data of the Plugin structure to the PluginEngine loading this plugin
		pdk.Output(p)
	}

    // Extism SDK requires a 0 for no error status
	return 0
}

func main() {}