.PHONY: all build_linux build_macos build_windows build_native clean install uninstall test

BUILD_VERSION := $(shell git describe --tags --always)
TAG_DATE := $(shell git log -1 --format=%cd --date=rfc $(BUILD_VERSION))
INSTALL_DIR := ${HOME}/go/bin
INSTALL_TARGET := $(INSTALL_DIR)/$(NAME)
NAME := cfgrr
LINUX_NAME := $(NAME)_linux
MACOS_NAME := $(NAME)_macos
WIN_NAME := $(NAME)_windows.exe
BUILD_TAGS := '-s -w -X "main.version=$(BUILD_VERSION)" -X "main.tagdate=$(TAG_DATE)"'

all: build_linux build_macos build_windows

build_linux:
	$(info Building Linux binary...)
	GOOS=linux go build -ldflags $(BUILD_TAGS) -o bin/$(LINUX_NAME)

build_macos:
	$(info Building MacOS binary...)
	GOOS=darwin go build -ldflags $(BUILD_TAGS) -o bin/$(MACOS_NAME)

build_windows:
	$(info Building windows binary...)
	GOOS=windows go build -ldflags $(BUILD_TAGS) -o bin/$(WIN_NAME)

build_native:
	$(info Building native binary...)
	go build -ldflags $(BUILD_TAGS) -o bin/$(NAME)

clean:
	$(info Cleaning up...)
	rm -rf bin

install: |
	$(info Installing $(INSTALL_TARGET))
	mkdir -p $(INSTALL_DIR) && go install -ldflags $(BUILD_TAGS)

uninstall:
	$(info Removing $(INSTALL_TARGET))
	rm -rf $(INSTALL_TARGET)

test:
	$(info Running tests...)
	go test -v ./...
