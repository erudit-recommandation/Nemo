GOCMD = go
BUILD_PATH = ./build
BINARY = $(BUILD_PATH)/website
TEST_PATH= ./test/...

.PHONY: build test

build:
	mkdir -p build
	$(GOCMD) build -o $(BINARY) main.go

run: build
	$(BINARY) -d

test:
	$(GOCMD) test $(TEST_PATH)

run-docker:
	docker run -it -p 8087:8087 erudit_recommandation

create-docker:
	docker build -t erudit_recommandation .

run-prod: build
	$(BINARY)