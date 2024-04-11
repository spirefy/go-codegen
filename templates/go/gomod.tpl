module {{.Module}}

go {{.GoVersion}}


require (
    {{.GoPluginPDK}} {{.GoPluginPDKVer}}
    {{.GoPluginEngine}} {{.GoPluginEngineVer}}
)
