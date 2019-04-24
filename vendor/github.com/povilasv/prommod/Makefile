LINT_FLAGS := run --deadline=120s
LINTER_EXE := golangci-lint
LINTER := ./bin/$(LINTER_EXE)
TESTFLAGS := -v -cover

GO111MODULE := on
all: $(LINTER) deps test lint build

$(LINTER):
#	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.15.0
#	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.15.0
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.15.0

.PHONY: lint
lint: $(LINTER)
	$(LINTER) $(LINT_FLAGS) ./...

.PHONY: deps
deps:
	go get .

.PHONY: build
build:
	go build .

.PHONY: test
test:
	go test $(TESTFLAGS) ./...
