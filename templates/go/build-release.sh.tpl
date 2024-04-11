# build with tinygo to use wasi wasm target, but do not include debug info, making binary much smaller
tinygo build -scheduler=none --no-debug -o {{.OutputPath}} -target wasi {{.SourceFile}}