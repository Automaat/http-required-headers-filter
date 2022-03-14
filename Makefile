golangci_lint := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.0

.PHONY: build
build:
	tinygo build -o required-header-filter.wasm -scheduler=none -target=wasi main.go

.PHONY: build-example
build-example:
	tinygo build -o ./example/required-header-filter.wasm -scheduler=none -target=wasi main.go

.PHONY: run-example
run-example:
	func-e run -c example/envoy.yaml

.PHONY: test
test:
	go test -tags=proxytest

.PHONY: lint
lint:
	go run $(golangci_lint) run --build-tags proxytest