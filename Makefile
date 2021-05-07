.PHONY: default
default: clean lint vet test

.PHONY: clean
clean:
	find . -name \*.coverprofile -delete

.PHONY: lint
lint:
	golint -set_exit_status ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -cover ./...

.PHONY: run-example
run-example:
	go run example/main.go
