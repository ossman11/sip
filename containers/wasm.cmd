SET GOOS=js
SET GOARCH=wasm
go build -o main.wasm || (printf "Failed to build go build -o main.wasm" && exit 1)
