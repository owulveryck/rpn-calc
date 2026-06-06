GOROOT := $(shell go env GOROOT)
WASMEXEC := $(firstword $(wildcard $(GOROOT)/lib/wasm/wasm_exec.js $(GOROOT)/misc/wasm/wasm_exec.js))

.PHONY: build test test-coverage run clean

build: web/rpn.wasm web/wasm_exec.js

web/rpn.wasm: cmd/wasm/main.go engine/*.go
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o $@ ./cmd/wasm/
	@command -v wasm-opt >/dev/null 2>&1 && { echo "Optimizing with wasm-opt..."; wasm-opt -Oz -o $@ $@; } || true

web/wasm_exec.js: $(WASMEXEC)
	cp "$(WASMEXEC)" web/

test:
	go test -v -coverprofile=coverage.out ./engine/...
	go tool cover -func=coverage.out | tail -1

test-coverage: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

run: build
	@echo "Serving on http://localhost:8080"
	@go run ./server/

clean:
	rm -f web/rpn.wasm web/wasm_exec.js coverage.out coverage.html
