include make.vars

TAG = $(VERSION)

BINDIR = bin
CURDIR = $(shell pwd)

PREFIX = $(REPOSITORY)/$(IMAGE_NAME)

GO_MOD_TIDY = go mod tidy

DOCKER_BUILD = docker build --no-cache
BUILDER_IMAGE = $(BUILD_IMAGE):$(BUILD_TAG)
GO_BUILD = CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w'

