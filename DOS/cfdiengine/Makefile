GOCMD=go
GOBUILD=$(GOCMD) build -ldflags="-w -s" 

FISCAL_EXE=./sidecar-fiscal-engine

all: build

build: fmt
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	$(GOBUILD) -o $(FISCAL_EXE) cmd/http/run.go

test:	build
	$(GOCMD) test -v ./...

fmt:
	$(GOCMD) fmt ./...

clean:
	rm -f $(FISCAL_EXE)
