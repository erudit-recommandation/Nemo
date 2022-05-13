GOCMD = /usr/local/go/bin/go
BUILD_PATH = ./build
BINARY = $(BUILD_PATH)/website

.PHONY: build

build:
	mkdir -p build
	$(GOCMD) build -o $(BINARY) main.go

run: build
	$(BINARY)

clear-docker:
	docker system prune -a -f

full-clear-docker:
	docker system prune -a -f --volumes
run-docker:
	sudo chmod -R a+rwx ./data/
	docker-compose up --build