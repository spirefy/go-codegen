# build with tinygo to use wasi wasm target, but do not include debug info, making binary much smaller
tinygo build -o ../plugins/openapiloader.wasm -no-debug -target wasi openapi.go