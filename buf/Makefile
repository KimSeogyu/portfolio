VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD)
version:
	@echo "Current version: $(VERSION)"

gen: clean
	buf generate --template buf.gen.yaml

	# generate proto files
	buf generate --template buf.gen.yaml
	
	# generate docs
	npx @redocly/cli@latest build-docs ../backend/internal/docs/openapi.yaml --output ../backend/internal/docs/openapi.html

clean:
	rm -rfd ../backend/internal/proto
	rm -rfd ../web/src/proto
	find ../backend/internal/docs -type f -not -name "embed.go" -delete
