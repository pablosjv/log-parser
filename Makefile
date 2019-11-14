BIN_NAME := log-parser
VERSION := $(shell grep "const Version " internal/version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
IMAGE_NAME := "pablosjv/log-parser"

default: help

.PHONY: help
help: ## Print this message.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


build: ## Build the binary. Place it under bin
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -o bin/${BIN_NAME} cmd/*.go

.PHONY: test
test:  ## Run tests
	go test ./...

build-docker: ## Build the docker image
	@echo "building image ${BIN_NAME} ${VERSION} $(GIT_COMMIT)"
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) --label version.type=dev -t $(IMAGE_NAME):local .

tag: ## Tag the docker images with version and commit
	@echo "Tagging: latest ${VERSION} $(GIT_COMMIT)"
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):$(GIT_COMMIT)
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):${VERSION}
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):latest

push: tag  ## Push to the docker repository
	@echo "Pushing docker image to registry: latest ${VERSION} $(GIT_COMMIT)"
	docker push $(IMAGE_NAME):$(GIT_COMMIT)
	docker push $(IMAGE_NAME):${VERSION}
	docker push $(IMAGE_NAME):latest

clean:  ## Remove the dev generated files and images
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}
	-docker rmi -f $(shell docker images -q --filter label=version.type=dev)
	-docker rmi -f $(shell docker images -q --filter "dangling=true")
