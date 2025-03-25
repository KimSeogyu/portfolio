
# Portfolio

## 소개 (Introduction)

Portfolio는 개인 포트폴리오용 프로젝트를 개발하고 있는 저장소입니다.

## Build & Development

### Requirements

- Go 1.24 이상
- Make
- [Buf CLI](https://buf.build/docs/installation) (Protocol Buffers 관리)
- Node.js 22.14.0 이상
- pnpm 10.7.0 이상

This project uses Make for build automation and development workflows. Here are the available commands:

### Main Commands

- `make build` - Build all components
- `make clean` - Clean all build artifacts
- `make test` - Run all tests
- `make up` - Start all services with Docker Compose
- `make down` - Stop all services
- `make help` - Show all available commands

### Backend-specific Commands

- `make backend` - Run backend make command
- `make backend-build` - Build backend application
- `make backend-clean` - Clean backend build artifacts
- `make backend-test` - Run backend tests
- `make backend-run` - Run backend locally
- `make backend-up` - Start backend services
- `make backend-down` - Stop backend services

### Buf-specific Commands

- `make buf` - Run buf make command
- `make buf-version` - Display current buf version
- `make buf-gen` - Generate proto files and API documentation

### Examples

```bash
# Build the entire project (including proto generation)
make build

# Run all tests
make test

# Start all services
make up

# Stop all services
make down

# Build only backend
make backend-build

# Run backend tests
make backend-test

# Generate proto files and API documentation
make buf-gen

# Display current buf version
make buf-version
```
