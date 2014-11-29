.PHONY: build doc lint run vendor_clean vendor_get vendor_update vet install

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
	&& mkdir -p ./_deps/src/github.com/MovingtoMars/gtkgain/src \
	&& cp -r ./src/library ./_deps/src/github.com/MovingtoMars/gtkgain/src

build: vet deps_get
	@go build -tags=gtk_3_10 -v -o ./bin/gtkgain ./src/gtkgain

doc:
	godoc -http=:6060 -index

run: build
	./bin/gtkgain

deps_update: deps_get
	rm -rf `find ./_deps/src -type d -name .git` \
	&& rm -rf `find ./_deps/src -type d -name .hg` \
	&& rm -rf `find ./_deps/src -type d -name .bzr` \
	&& rm -rf `find ./_deps/src -type d -name .svn`

vet:
	go vet ./...

install:
	@mkdir -p ${INSTALL_ROOT}/usr/bin \
	&& cp ./bin/gtkgain ${INSTALL_ROOT}/usr/bin \
	&& mkdir -p ${INSTALL_ROOT}/usr/share/applications \
	&& cp ./gtkgain.desktop ${INSTALL_ROOT}/usr/share/applications \
	&& mkdir -p ${INSTALL_ROOT}/usr/share/icons/hicolor/22x22/apps \
	&& cp ./src/icons/22x22/gtkgain.png ${INSTALL_ROOT}/usr/share/icons/hicolor/22x22/apps \
	&& mkdir -p ${INSTALL_ROOT}/usr/share/icons/hicolor/24x24/apps \
	&& cp ./src/icons/24x24/gtkgain.png ${INSTALL_ROOT}/usr/share/icons/hicolor/24x24/apps \
	&& mkdir -p ${INSTALL_ROOT}/usr/share/icons/hicolor/32x32/apps \
	&& cp ./src/icons/32x32/gtkgain.png ${INSTALL_ROOT}/usr/share/icons/hicolor/32x32/apps \
	&& mkdir -p ${INSTALL_ROOT}/usr/share/icons/hicolor/48x48/apps \
	&& cp ./src/icons/48x48/gtkgain.png ${INSTALL_ROOT}/usr/share/icons/hicolor/48x48/apps \
	&& mkdir -p ${INSTALL_ROOT}/usr/share/icons/hicolor/64x64/apps \
	&& cp ./src/icons/64x64/gtkgain.png ${INSTALL_ROOT}/usr/share/icons/hicolor/64x64/apps \
	&& mkdir -p ${INSTALL_ROOT}/usr/share/icons/hicolor/128x128/apps \
	&& cp ./src/icons/128x128/gtkgain.png ${INSTALL_ROOT}/usr/share/icons/hicolor/128x128/apps \
	&& mkdir -p ${INSTALL_ROOT}/usr/share/icons/hicolor/256x256/apps \
	&& cp ./src/icons/256x256/gtkgain.png ${INSTALL_ROOT}/usr/share/icons/hicolor/256x256/apps \
	&& xdg-icon-resource forceupdate --theme hicolor &> /dev/null
