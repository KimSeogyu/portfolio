.PHONY: all build clean test run up down backend buf buf-version buf-gen

# 기본 타겟
all: backend

# Backend 관련 명령어들
backend:
	@$(MAKE) -C backend

backend-build:
	@$(MAKE) -C backend build

backend-clean:
	@$(MAKE) -C backend clean

backend-test:
	@$(MAKE) -C backend test

backend-run:
	@$(MAKE) -C backend run

backend-up:
	@$(MAKE) -C backend up

backend-down:
	@$(MAKE) -C backend down

backend-generate-mocks:
	@$(MAKE) -C backend generate-mocks

# Buf 관련 명령어들
buf:
	@$(MAKE) -C buf

buf-version:
	@$(MAKE) -C buf version

buf-gen:
	@$(MAKE) -C buf gen

# 전체 프로젝트 명령어들
build: buf-gen backend-build

clean: backend-clean buf-clean

test: backend-test

up: backend-up

down: backend-down

# 도움말
help:
	@echo "Available targets:"
	@echo "  all          - Build all components (default)"
	@echo "  build        - Build all components"
	@echo "  clean        - Clean all build artifacts"
	@echo "  test         - Run all tests"
	@echo "  up           - Start all services"
	@echo "  down         - Stop all services"
	@echo "  backend      - Run backend make command"
	@echo "  backend-*    - Run specific backend command (build/clean/test/run/up/down)"
	@echo "  buf          - Run buf make command"
	@echo "  buf-version  - Display buf version"
	@echo "  buf-gen      - Generate proto files and documentation"