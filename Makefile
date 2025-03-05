build:
	@CGO_ENABLED=0 go build -o bin/rouge.exe

build_small:
	@CGO_ENABLED=0 go build -ldflags "-s -w -buildid=" -trimpath -o bin/rouge.exe

build_tiny:
	GOOS=windows tinygo build -o bin/rouge.exe

run: build
	@./bin/rouge.exe

test:
	@go test ./... -v