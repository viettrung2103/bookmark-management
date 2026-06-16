IMAGE_NAME=viettrung21/bookmark-service
TAG=latest

.PHONY: run swagger dev-run test build-push-vm deploy-all

run:
	go run cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go

docker:
	docker build -t test_img:latest .
	docker run --rm -p 8080:8080 test_img:latest

dev-run: docker swagger run

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

# 1. Build for VM (AMD64) and push immediately to Docker Hub
build-push-vm:
	@echo "🚀 Building image for linux/amd64 (VM)..."
	docker buildx build --platform linux/amd64 -t $(IMAGE_NAME):$(TAG) --push .

# 2. OPTIONAL: Build for both your local Mac (ARM64) and your VM (AMD64)
build-push-multi:
	@echo "🚀 Building multi-arch image (ARM64 + AMD64)..."
	docker buildx build --platform linux/amd64,linux/arm64 -t $(IMAGE_NAME):$(TAG) --push .