# These are the values we want to pass for Version and BuildTime
VERSION=1.2.0
BUILD_TIME=$(shell date +%s)

# Setup the -ldflags option for go build here, interpolate the variable values

LDFLAGS += -X \"main.Version=$(VERSION)\"
LDFLAGS += -X \"main.BuildTime=$(BUILD_TIME)\"

OSARCH=darwin/amd64 linux/386 linux/amd64 linux/arm windows/386 windows/amd64

.PHONY: pack build clean

build:
	go build -ldflags "$(LDFLAGS)"

install:
	go install -ldflags "$(LDFLAGS)"

cp:
	for f in $(wildcard out/*); do cp LICENSE README.md $$f; done;

pack:
	rm -rf "./out/"
	gox -ldflags="$(LDFLAGS)" -osarch="$(OSARCH)" -output="./out/{{.OS}}_{{.Arch}}/chatlab"
	./scripts/package.sh

clean:
	go clean -i ./...

clean-mac: clean
	find . -name ".DS_Store" -print0 | xargs -0 rm
