ORGNAME := wabarc
PROJECT := starter
HOMEDIR ?= /go/src/github.com/${ORGNAME}/${PROJECT}
DOCKER ?= $(shell command -v docker || command -v podman)
IMAGE := wabarc/golang-chromium:dev

all: build \
	clean

build: submodule \
	patch \
	buster \
	movext \
	starter

buster:
	@echo "Copy secrets.json to source entry..."
	cp hack/secrets.json external/buster/secrets.json
	@echo "Packaging buster extension..."
	$(DOCKER) run -i --rm -v `pwd`:/workspace node:14-alpine sh -c 'cd /workspace/external/buster; \
		apk add --no-cache build-base automake autoconf libtool nasm libpng-dev zlib-dev; \
		yarn; \
		yarn add mozjpeg; \
		yarn add --dev --platform=linuxmusl sharp; \
		yarn build:prod:chrome'

starter:
	@echo "Building starter..."
	@rm -f starter
	go build -trimpath --ldflags "-s -w" -o starter *.go

submodule:
	@echo "Updating git submodule..."
	@rm -rf external/*
	git submodule update --init --recursive

movext:
	@echo "Moving dist to extensions..."
	@rm -rf extensions/buster
	cp -r external/buster/dist/chrome extensions/buster
	@rm -rf extensions/bypass-paywall
	cp -r external/bypass-paywall extensions/bypass-paywall

patch:
	@echo "Apply patches..."
	git apply --directory=external/bypass-paywall/ patches/bypass-paywall-src-js-options.patch

clean:
	@echo "Cleaning dist..."
	rm -rf demo/* external/* starter installed-extensions.png

run:
	@echo "-> Running docker container"
	$(DOCKER) run --memory="500m" -ti --rm -e DISPLAY=:99.0 -v ${PWD}:${HOMEDIR} ${IMAGE} sh -c "\
		cd ${HOMEDIR} && \
		go get -v && \
		sh"

demo:
	@echo "-> Running docker container"
	$(MAKE) build
	$(DOCKER) run --memory="500m" -i --rm -e DISPLAY=:99.0 -v ${PWD}:${HOMEDIR} ${IMAGE} sh -c "\
		cd ${HOMEDIR} && \
		sh hack/demo.sh && \
		sh hack/0x0.st.sh installed-extensions.png"

.PHONY: all
