all: tidy	test	install

tidy:
	@echo "Cleaning up go.mod..."
	@go mod tidy -v

install:
	@echo "Installing parens to GOBIN..."
	@go install ./cmd/parens/

build:
	@echo "Building parens at ./bin/parens"
	@go build -o bin/parens ./cmd/parens/*.go

test:
	@go test -cover ./...


benchmark:
	@go test -bench=Benchmark -benchmem ./...
