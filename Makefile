NAME=videotranscoder
VERSION=$(shell cat VERSION)
BUILD=$(shell git rev-parse --short HEAD)
LD_FLAGS="-w -X main.version=$(VERSION) -X main.build=$(BUILD)"

clean:
	rm -rf _build/ release/

build:
	go mod download
	CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o ${NAME}

build-all:
	mkdir -p _build
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o _build/${NAME}-$(VERSION)-darwin-amd64
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o _build/${NAME}-$(VERSION)-linux-amd64
	GOOS=linux   GOARCH=arm   CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o _build/${NAME}-$(VERSION)-linux-arm
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o _build/${NAME}-$(VERSION)-linux-arm64
	# GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o _build/${NAME}-$(VERSION)-windows-amd64
	cd _build; sha256sum * > sha256sums.txt

smoketest:
	_build/${NAME}-$(VERSION)-linux-amd64 -help | grep "version ${VERSION}"

release:
	mkdir release
	cp _build/* release
	cd release; sha256sum --quiet --check sha256sums.txt && \
	gh release create $(VERSION) -d -t v$(VERSION) *

.PHONY: build