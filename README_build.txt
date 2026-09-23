CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o magnetico-linux-arm64 -v -tags fts5 .
