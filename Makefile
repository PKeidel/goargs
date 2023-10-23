.PHONY: release $(PLATFORMS) dev help test cover bench build build-release build-gcc install run clean update-libs

dev: ## Runs the app locally in dev mode
	go run -tags debug ./cmd/ --host localhost --debug --user PKeidel

build:
	go build -o goargsgenerate ./cmd/generate

help: ## Prints all available make commands
	@grep -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":"}; {printf "\033[36m%-30s\033[0m\n", $$1}'

test: ## Runs the tests for the app
	go test .

testd: ## Runs the tests for the app
	go test -tags debug -v .

cover: ## Runs the tests with coverage enabled and opens the report in a local browser
	go test -race -coverprofile=coverage.out -v .
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'
	@go tool cover -html=coverage.out

bench: ## Runs the benchmarks
	go test -bench=. .

clean: ## Cleans everything
	go clean
	rm ${BINARY_NAME}

fmt:
	find . -name '*.go' -exec bash -c 'echo "{}"; gofmt -w {}' \;

update-libs: ## Update all go libs
	go get -u ./...
