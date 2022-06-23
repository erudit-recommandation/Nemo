GOCMD = go
BUILD_PATH = ./build
BINARY = $(BUILD_PATH)/website
TEST_PATH= ./test/...

.PHONY: build test

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

run-gemsim-service:
	cd ./gemsim_service && export FLASK_APP=app && flask run

test:
	$(GOCMD) test $(TEST_PATH)