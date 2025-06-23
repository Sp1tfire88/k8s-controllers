# Makefile

APP_NAME = controller
BUILD_DIR = build

.PHONY: all build run clean test docker

all: build

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) main.go

run: build
	./$(BUILD_DIR)/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...

docker:
	docker build -t sp1tfire88/$(APP_NAME):latest .

