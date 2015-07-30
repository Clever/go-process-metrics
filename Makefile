SHELL := /bin/bash
PKGS := \
github.com/Clever/go-process-metrics/metrics

GOLINT := $(GOPATH)/bin/golint

.PHONY: test

all: build

test: $(PKGS)

$(GODEP):
	@go get github.com/tools/godep

$(GOLINT):
	@go get github.com/golang/lint/golint

$(PKGS): $(GOLINT) $(GODEP)
	$(GODEP) go install $@
	@gofmt -w=true $(GOPATH)/src/$@/*.go
	@echo "LINTING..."
	$(GOPATH)/bin/golint $(GOPATH)/src/$@/*.go
	@echo "VETTING..."
	go vet $(GOPATH)/src/$@/*.go
	@echo ""
ifeq ($(COVERAGE),1)
	go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
	go tool cover -html=$(GOPATH)/src/$@/c.out
else
	@echo "TESTING..."
	go test $@ -test.v
	@echo ""
endif
