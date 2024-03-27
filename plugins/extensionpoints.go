package plugins

import t "github.com/spirefy/go-pdk/types"

func GetCodegenExtensionPoints() []t.ExtensionPoint {
	// This extension point expects a Menu type to be provided by an extension.
	// It will call each extension once to obtain the name, parent name and onclick func name to call
	ep := t.CreateExtensionPoint(
		"codegen.plugins.loaders",
		"Codegen Source Loader ExtensionPoint",
		"1.0.0",
		"This extension point provides the functionality for extensions to provide source format loaders that will add to this plugins unified in memory structure that is then passed to generators",
		nil)

	return []t.ExtensionPoint{*ep}
}
