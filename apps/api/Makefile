.PHONY: gen run test lint build

## gen: regenerate internal/openapi from the contract
gen:
	go generate ./internal/openapi

## run: run the api on your machine (needs DATABASE_URL)
run:
	go run ./cmd/server

## test: unit tests; repository tests also hit postgres when TEST_DATABASE_URL is set
test:
	go test -race -cover ./...

## lint: formatting, vet, and generated code in sync with the contract
lint: gen
	test -z "$$(gofmt -l .)"
	go vet ./...
	git diff --exit-code -- internal/openapi

build:
	CGO_ENABLED=0 go build -trimpath -o bin/server ./cmd/server
