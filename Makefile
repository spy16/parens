all: test	install

install:
	@echo "Installing parens to GOBIN..."
	@go install ./cmd/parens/

build:
	@echo "Building parens at ./bin/parens"
	@go build -o bin/parens ./cmd/parens/*.go

test:
	@go test -cover ./. ./lexer/... ./parser/... ./reflection/...

