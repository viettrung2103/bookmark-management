#.PHONY: run swagger dev-run test
#
#run:
#	go run cmd/api/main.go
#
#swagger:
#	swag init -g cmd/api/main.go
#
#dev-run: swagger run
#
#COVERAGE_EXCLUDE=mocks|main.go|test|pkg
#COVERAGE_THRESHOLD = 40

#test:
#	$(eval COVER_PKGS := $(shell go list ./... | grep -vE "cmd|docs|mocks" | tr '\n' ',' | sed 's/,$$//'))
#	$(eval TEST_PKGS  := $(shell go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./... | grep -vE "cmd|docs|mocks"))
#	go test $(TEST_PKGS) -coverprofile=coverage.tmp -covermode=atomic -coverpkg=$(COVER_PKGS) -p 1
#	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
#	go tool cover -html=coverage.out -o coverage.html
#	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | tr -d '%'); \
#	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
#	   echo "Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
#	   exit 1; \
#	else \
#	   echo "Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
#	fi

#.PHONY: run swagger dev-run test
IMG_NAME=viettrung21/bookmark-service
GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null)
BRANCH 	:= $(shell git rev-parse --abbrev-ref HEAD)

#IMG_TAG := dev

ifeq ($(BRANCH),main)
	IMG_TAG := dev
endif

ifneq ($(GIT_TAG),)
	IMG_TAG := $(GIT_TAG)
endif

export IMG_TAG

.PHONY: run swagger dev-run test docker-test docker-build docker-release docker-login
run:
	go run cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go output docs

dev-run: swagger run

COVERAGE_EXCLUDE=mocks|main.go|test|docs|test|config.go
COVERAGE_THRESHOLD = 50

test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | tr -d '%'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

COVERAGE_FOLDER =./coverage

docker-test:
	mkdir -p ${COVERAGE_FOLDER}
	docker buildx build --build-arg COVERAGE_EXCLUDE="$(COVERAGE_EXCLUDE)" --target test -t bookmark-service:dev --output $(COVERAGE_FOLDER) .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "X Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "checked Coverage($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

docker-build:
	docker build -t $(IMG_NAME):$(IMG_TAG) .

docker-release:
	docker push $(IMG_NAME):$(IMG_TAG)

DOCKER_USERNAME ?=
DOCKER_PASSWORD ?=

docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

generate-rsa-key:
	openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -pubout -in private.pem -out  public.pem