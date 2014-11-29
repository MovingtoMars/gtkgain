.PHONY: build doc lint run vendor_clean vendor_get vendor_update vet

# Prepend our _vendor directory to the system GOPATH
# so that import path resolution will prioritize
# our third party snapshots.
GOPATH := ${PWD}/_deps:${GOPATH}
export GOPATH

default: build deps_clean

deps_clean:
	rm -dRf ./_deps

deps_get: deps_clean
	GOPATH=${PWD}/_deps go get -d -u -v \
	github.com/MovingtoMars/gotk3/gtk \
	github.com/vchimishuk/chub/src/ogg/libvorbis \
	github.com/wtolson/go-taglib \
	&& mkdir ./_deps/src/github.com/MovingtoMars/gtkgain \
	&& cp -r ./src/library ./_deps/src/github.com/MovingtoMars/gtkgain

build: vet deps_get
	go build -tags=gtk_3_10 -v -o ./bin/gtkgain ./src/main

doc:
	godoc -http=:6060 -index

run: build
	./bin/gtkgain

deps_update: deps_get
	rm -rf `find ./_deps/src -type d -name .git` \
	&& rm -rf `find ./_deps/src -type d -name .hg` \
	&& rm -rf `find ./_deps/src -type d -name .bzr` \
	&& rm -rf `find ./_deps/src -type d -name .svn`

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet ./...
