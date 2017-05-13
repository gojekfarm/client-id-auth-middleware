.PHONY: all
all: build-deps build fmt vet lint test

GLIDE_NOVENDOR=$(shell glide novendor)

setup:
	mkdir -p $(GOPATH)/bin
	go get -u github.com/golang/lint/golint

build-deps:
	glide install

update-deps:
	glide update

compile:
	mkdir -p out/
	go build -race $(GLIDE_NOVENDOR)

build: build-deps compile fmt vet lint

fmt:
	go fmt $(GLIDE_NOVENDOR)

vet:
	go vet $(GLIDE_NOVENDOR)

lint:
	@for p in $(UNIT_TEST_PACKAGES); do \
		echo "==> Linting $$p"; \
		golint -set_exit_status $$p; \
	done

test:
	ENVIRONMENT=test go test -race $(GLIDE_NOVENDOR)