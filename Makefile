.PHONY: build clean install uninstall test

NAME=cfgrr
BUILD_VERSION=$(shell git describe --tags --always)
BUILD_DATE=$(shell date -R)
INSTALL_DIR=${HOME}/go/bin
INSTALL_TARGET=$(INSTALL_DIR)/$(NAME)

build:
	GOOS=linux go build -ldflags '-s -w -X "main.version=$(BUILD_VERSION)" -X "main.builddate=$(BUILD_DATE)"' -o bin/$(NAME)

clean:
	rm -rf bin

cp_build:
	mkdir -p $(INSTALL_DIR) && cp bin/$(NAME) $(INSTALL_DIR)

install: build cp_build clean

uninstall:
	rm -rf $(INSTALL_TARGET)

test:
	go test -v ./...
