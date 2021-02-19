APP_NAME = milv
APP_PATH = tools/$(APP_NAME)
IMG_NAME := $(DOCKER_PUSH_REPOSITORY)$(DOCKER_PUSH_DIRECTORY)/$(APP_NAME)
TAG := $(DOCKER_TAG)

.PHONY: build-image
build-image:
	docker build -t $(APP_NAME):latest .

.PHONY: push-image
push-image:
	docker tag $(APP_NAME) $(IMG_NAME):$(TAG)
	docker push $(IMG_NAME):$(TAG)

.PHONY: release
release: component-check build-image push-image

.PHONY: deps
deps: vendor verify

.PHONY: check
check: test check-fmt

.PHONY: component-check
component-check: deps check

vendor:
	GO111MODULE=on go mod vendor

.PHONY: verify
verify:
	GO111MODULE=on go mod verify

.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

.PHONY: test
test:
	go test ./...

build:
	go build -o milv main.go

VERIFY_IGNORE := /vendor\|/automock
FILES_TO_CHECK = find . -type f -name "*.go" | grep -v "$(VERIFY_IGNORE)"
.PHONY: check-fmt
check-fmt:
	@if [ -n "$$(goimports -l $$($(FILES_TO_CHECK)))" ]; then \
		echo "✗ some files contain not propery formatted imports. To repair run make fmt"; \
		goimports -l $$($(FILES_TO_CHECK)); \
		exit 1; \
	fi;

.PHONY: fmt
fmt:
	goimports -w -l $$($(FILES_TO_CHECK))