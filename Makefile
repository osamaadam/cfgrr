.PHONY: build_linux build_macos build_windows build_native build clean install uninstall test cp_build

BUILD_VERSION=$(shell git describe --tags --always)
TAG_DATE=$(shell git log -1 --format=%cd --date=rfc $(BUILD_VERSION))
INSTALL_DIR=${HOME}/.local/bin
INSTALL_TARGET=$(INSTALL_DIR)/$(NAME)
NAME=cfgrr
LINUX_NAME=$(NAME)_linux
MACOS_NAME=$(NAME)_macos
WIN_NAME=$(NAME)_windows.exe

build_linux:
	$(info Building Linux binary...)
	GOOS=linux go build -ldflags '-s -w -X "main.version=$(BUILD_VERSION)" -X "main.tagdate=$(TAG_DATE)"' -o bin/$(LINUX_NAME)

build_macos:
	$(info Building MacOS binary...)
	GOOS=darwin go build -ldflags '-s -w -X "main.version=$(BUILD_VERSION)" -X "main.tagdate=$(TAG_DATE)"' -o bin/$(MACOS_NAME)

build_windows:
	$(info Building windows binary...)
	GOOS=windows go build -ldflags '-s -w -X "main.version=$(BUILD_VERSION)" -X "main.tagdate=$(TAG_DATE)"' -o bin/$(WIN_NAME)

build_native:
	$(info Building native binary...)
	go build -ldflags '-s -w -X "main.version=$(BUILD_VERSION)" -X "main.tagdate=$(TAG_DATE)"' -o bin/$(NAME)

build: build_linux build_macos build_windows

clean:
	$(info Cleaning up...)
	rm -rf bin

cp_build:
	$(info Copying binary to $(INSTALL_DIR))
	mkdir -p $(INSTALL_DIR) && cp bin/$(NAME) $(INSTALL_DIR)

install: build_native cp_build clean

uninstall:
	$(info Removing $(INSTALL_TARGET))
	rm -rf $(INSTALL_TARGET)

test:
	$(info Running tests...)
	go test -v ./...
