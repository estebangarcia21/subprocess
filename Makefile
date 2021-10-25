.PHONY: test

build:
	@echo Building
	@go build .

test:
	@go test -v .

bench:
	@go test -bench .
