DOCKER=docker run --rm \
         -v ./:/usr/src/app \
         -w /usr/src/app \
         -p 8080:8080 \
         golang:1.24

SAFE_GIT=git config --global --add safe.directory /usr/src/app

COMMAND ?= start

APP_NAME ?= app-test


#### Linting and formatting commands
docker-lint:
	docker run --rm -v ./:/app -w /app golangci/golangci-lint:latest golangci-lint run
.PHONY: docker-lint

docker-test:
	$(DOCKER) go test ./...
.PHONY: docker-test

##### Build application in docker -> make docker-build
docker-build:
	$(DOCKER) bash -c "$(SAFE_GIT) && go build -v -o $(APP_NAME)"
.PHONY: docker-build

##### Run application in docker -> make docker-run
docker-run:
	$(DOCKER) ./$(APP_NAME) $(COMMAND)
.PHONY: docker-run

##### Build application from Golang -> make go-build
go-build:
	go build -v -o $(APP_NAME)
.PHONY: go-build

##### Run application from Golang -> make go-run
go-run:
	go run main.go $(COMMAND)
.PHONY: go-run

