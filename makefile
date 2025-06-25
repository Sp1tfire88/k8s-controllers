# APP_NAME = controller
# BUILD_DIR = build
# GO_FILES = main.go
# DOCKER_IMAGE = sp1tfire88/$(APP_NAME)
# GO_VERSION = 1.21
# COVERAGE_DIR = coverage
# COVERAGE_FILE = $(COVERAGE_DIR)/coverage.out

# .PHONY: all build run clean test docker docker-run docker-push tidy lint

# all: build

# # Создание папки и билд
# build:
# 	@mkdir -p $(BUILD_DIR)
# 	go build -o $(BUILD_DIR)/$(APP_NAME) $(GO_FILES)

# # Запуск бинарника
# run: build
# 	./$(BUILD_DIR)/$(APP_NAME)

# # Очистка
# clean:
# 	rm -rf $(BUILD_DIR)

# # Запуск тестов
# test:
# 	go test -v ./cmd

# # Docker билд
# docker:
# 	docker build -t $(DOCKER_IMAGE):latest .

# # Запуск контейнера
# docker-run:
# 	docker run --rm -p 8080:8080 $(DOCKER_IMAGE):latest

# # Пуш в Docker Hub
# docker-push:
# 	docker push $(DOCKER_IMAGE):latest

# # Приведение зависимостей в порядок
# tidy:
# 	go mod tidy

# # Линтинг (можно вызывать в GitHub Actions)
# lint:
# 	golangci-lint run

# coverage:
# 	@mkdir -p $(COVERAGE_DIR)
# 	go test -v -coverprofile=$(COVERAGE_FILE) ./cmd
# 	go tool cover -func=$(COVERAGE_FILE)
# 	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_DIR)/coverage.html

APP_NAME = controller
BUILD_DIR = build
DOCKER_IMAGE = $(APP_NAME)
TAG = latest
COVERAGE_DIR = coverage
COVERAGE_FILE = $(COVERAGE_DIR)/coverage.out

.PHONY: all build run clean test docker coverage

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) main.go

run: build
	./$(BUILD_DIR)/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR) $(COVERAGE_DIR)

test:
	go test -v ./...

coverage:
	@mkdir -p $(COVERAGE_DIR)
	go test -v -coverprofile=$(COVERAGE_FILE) ./cmd
	go tool cover -func=$(COVERAGE_FILE)
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_DIR)/coverage.html

docker:
	docker build -t $(DOCKER_IMAGE):$(TAG) .

# Линтинг (можно вызывать в GitHub Actions)
lint:
	golangci-lint run