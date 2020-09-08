PLUGIN_NAME := pack

.PHONY: build
build: build_linux build_mac build_windows

build_linux: export GOARCH=amd64
build_linux: export CGO_ENABLED=0
#build_linux: export GO111MODULE=on
build_linux: export GOPROXY=https://gocenter.io
build_linux:
	@GOOS=linux go build -v -mod=vendor -o bin/linux/amd64/helmpack .  # linux

link_linux:
	@cp bin/linux/amd64/helmpack /usr/local/bin/helmpack

build_mac: export GOARCH=amd64
build_mac: export CGO_ENABLED=0
#build_mac: export GO111MODULE=on
build_mac: export GOPROXY=https://gocenter.io
build_mac:
	@GOOS=darwin go build -v -mod=vendor -o bin/darwin/amd64/helmpack . # mac osx

link_mac:
	@cp bin/darwin/amd64/helmpack /usr/local/bin/helmpack

build_windows: export GOARCH=amd64
#build_windows: export GO111MODULE=on
build_windows: export GOPROXY=https://gocenter.io
build_windows:
	@GOOS=windows go build -v -mod=vendor -o bin/windows/amd64/helmpack .  # windows

link_windows:
	@cp bin/windows/amd64/helmpush ./bin/helmpush

.PHONY: clean
clean:
	@git status --ignored --short | grep '^!! ' | sed 's/!! //' | xargs rm -rf


.PHONY: install
install:
	HELM_PUSH_PLUGIN_NO_INSTALL_HOOK=1 helm plugin install $(shell pwd)

.PHONY: remove
remove:
	helm plugin remove $(PLUGIN_NAME)