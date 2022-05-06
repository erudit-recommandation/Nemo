GOCMD = /usr/local/go/bin/go
BUILD_PATH = ./build
BINARY = $(BUILD_PATH)/website

.PHONY: build

build:
	mkdir -p build
	$(GOCMD) build -o $(BINARY) main.go

run: build
	$(BINARY)