.PHONY: serve
serve:
	go run ./cmd

.PHONY: test
test:
	go test ./...
