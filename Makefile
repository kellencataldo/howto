.PHONY: fmt

BIN = $(CURDIR)/bin
$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN)
	@echo downloading lint package: $(PACKAGE)
	@tmp=$$(mktemp -d); env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) go get $(PACKAGE) || ret=$$?; rm -rf $$tmp; exit $$ret

$(BIN)/golint: PACKAGE=golang.org/x/lint/golint

GOLINT = $(BIN)/golint
lint: | $(GOLINT)
	$(GOLINT) -set_exit_status ./...

serverht:
	go build ./cmd/htserver.go

deploy: server

howto:
	go build ./cmd/howto.go

clean:
	rm -f howto server
	rm -rf bin

fmt:
	go fmt ./...
