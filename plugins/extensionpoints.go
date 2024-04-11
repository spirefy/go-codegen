package plugins

import t "github.com/spirefy/go-pdk/types"

func GetCodegenExtensionPoints() []t.ExtensionPoint {
	// loaderExtensionPoint
	// This extension point sets up source loaders that will parse their supported source format in to some portion of
	// the unified in memory structure. Each loader will accept as input the location of the file to load which can be
	// a local relative file path or a URL to a file available over the internet. The passed in structure to each
	// extension is:
	//
	// {
	//   "url": "<url if its a url or left out/null/empty if not",
	//   "path": "<local relative path to a file on the file system if it's a local file, or null/empty if not",
	//   "sourceType": "optional helpful type of file.. may be useful for loaders, optional not needed",
	// }
	//
	// The return structure of a loader should be a
	loaderExtensionPoint := t.CreateExtensionPoint(
		"codegen.plugins.loaders",
		"Codegen Source Loader ExtensionPoint",
		"1.0.0",
		"This extension point provides the functionality for extensions to provide source format loaders that will add to this plugins unified in memory structure that is then passed to generators",
		nil)

	generatorExtensionPoint := t.CreateExtensionPoint(
		"codegen.plugins.generators",
		"Codegen Target Generator ExtensionPoint",
		"1.0.0",
		"This extension point provides the functionality for extensions to provide target generators that will use the unified in memory structure to generate output from",
		nil)

	return []t.ExtensionPoint{*loaderExtensionPoint, *generatorExtensionPoint}
}
