test:
	go test ./...

bench:
	go test -benchmem -bench .

.PHONY: test bench repl

