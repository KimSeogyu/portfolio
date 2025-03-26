GOOS ?= linux
GOARCH ?= amd64
OUT_DIR = build
APPLICATIONS = board

.PHONY: build run clean test docker up down $(APPLICATIONS)

build: $(APPLICATIONS)

$(APPLICATIONS):
	@$(MAKE) build-app APPLICATION=$@
	@$(MAKE) docker-image APPLICATION=$@

build-app:
	@echo "Building ${APPLICATION}..."
	@go build $(BUILD_FLAGS) -o ${OUT_DIR}/${GOOS}/${GOARCH}/${APPLICATION} ./cmd/${APPLICATION}/main.go
	@echo "Build ${APPLICATION} success"

docker-image:
	@echo "Building docker image..."
	@docker build -f docker/Dockerfile.${APPLICATION} -t portfolio/${APPLICATION} .
	@echo "Build docker image success"

run:
	go run cmd/board/main.go

clean:
	rm -rf build
	rm -rf .config
	rm -rf volumes
	docker rmi portfolio/board

test:
	go test ./...

up: 
	$(MAKE) build GOOS=linux GOARCH=amd64
	docker compose -f docker/docker-compose.yml up

generate-mocks:
	mockgen -source=cmd/board/internal/comments/repository.go -destination=cmd/board/internal/comments/mocks/repository.go
	mockgen -source=cmd/board/internal/postings/repository.go -destination=cmd/board/internal/postings/mocks/repository.go
	mockgen -source=internal/cache/cache.go -destination=internal/cache/mocks/cache.go

down:
	docker compose -f docker/docker-compose.yml down