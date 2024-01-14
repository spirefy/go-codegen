package plugins

import t "github.com/spirefy/go-pdk/types"

func GetCodegenExtensionPoints() []t.ExtensionPoint {
	// This extension point expects a Menu type to be provided by an extension.
	// It will call each extension once to obtain the name, parent name and onclick func name to call
	ep := t.ExtensionPoint{
		Id:          "codegen.plugins.loaders",
		Description: "This extension point provides the functionality for extensions to provide source format loaders that will add to this plugins unified in memory structure that is then passed to generators",
		Name:        "Codegen Source Loader ExtensionPoint",
	}

	return []t.ExtensionPoint{ep}
}
