tinygo build -scheduler=none -o ./out/codegen.wasm -target wasi codegen.go
tar -czvf ../plugins/codegen.tar.gz -C ./ plugin.yaml -C ./out/ codegen.wasm

