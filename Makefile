.PHONY: build clean install uninstall test

NAME=cfgrr
BUILD_VERSION=$(shell git describe --tags --always)
TAG_DATE=$(shell git log -1 --format=%cd --date=rfc $(BUILD_VERSION))
INSTALL_DIR=${HOME}/go/bin
INSTALL_TARGET=$(INSTALL_DIR)/$(NAME)

build:
	GOOS=linux go build -ldflags '-s -w -X "main.version=$(BUILD_VERSION)" -X "main.tagdate=$(TAG_DATE)"' -o bin/$(NAME)

clean:
	rm -rf bin

cp_build:
	mkdir -p $(INSTALL_DIR) && cp bin/$(NAME) $(INSTALL_DIR)

install: build cp_build clean

uninstall:
	rm -rf $(INSTALL_TARGET)

test:
	go test -v ./...
