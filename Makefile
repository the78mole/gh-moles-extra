.PHONY: update-version build test clean

# Update version in manifest from git tag
update-version:
	@./scripts/update-version.sh

# Build for current platform
build:
	@go build -o gh-moles .

# Build for all platforms
build-all:
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/gh-moles-linux-amd64 .
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/gh-moles-linux-arm64 .
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/gh-moles-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/gh-moles-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/gh-moles-windows-amd64.exe .

# Run tests
test:
	@go test -v ./...

# Clean build artifacts
clean:
	@rm -f gh-moles
	@rm -rf dist/

# Development setup
dev: update-version build