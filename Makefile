.PHONY: test unit-test integration lint vet run-example
test unit-test:
	go test -race -count=1 ./...

integration:
	go test -tags=integration -run '^TestMockIntegration' -count=1 -timeout=3m ./...

lint:
	test -z "$$(gofmt -l .)"

vet:
	go vet ./...

run-example:
	go run ./example
