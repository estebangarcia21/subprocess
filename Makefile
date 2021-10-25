.PHONY: test

test:
	go test -v .

bench:
	go test -bench .
