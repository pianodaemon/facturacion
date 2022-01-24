GOCMD=go
GOBUILD=$(GOCMD) build -ldflags="-w -s"

INTERVIEW_EXE=interview

all: build

build: fmt
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	$(GOBUILD) -o $(INTERVIEW_EXE) cmd/http/run.go

fmt:
	$(GOCMD) fmt ./...

clean:
	rm -f $(INTERVIEW_EXE)
