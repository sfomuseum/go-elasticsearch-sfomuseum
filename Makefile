GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/es-sfomuseum-index cmd/es-sfomuseum-index/main.go
