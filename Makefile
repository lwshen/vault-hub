.PHONY: generate-web-assets build-server build clean dev-server help

# Generate web assets (build frontend and copy to internal/web)
generate-web-assets:
	go run ./cmd/generate-web-assets

# Build server with embedded web assets
build-server: generate-web-assets
	go build -o vault-hub ./apps/server

# Build everything (web + server with embedded assets)
build: build-server

# Clean build artifacts
clean:
	rm -f vault-hub vault-hub-embedded
	rm -rf internal/web/dist
	rm -rf apps/web/dist

# Development server (without embedded assets)
dev-server:
	go run ./apps/server

# Quick development build (build web assets and run server)
dev-build:
	$(MAKE) generate-web-assets
	go run ./apps/server

help:
	@echo "Available targets:"
	@echo "  generate-web-assets - Build web app and copy to internal/web"
	@echo "  build-server        - Build server with embedded web assets"
	@echo "  build               - Build everything (web + server)"
	@echo "  clean               - Clean build artifacts"
	@echo "  dev-server          - Run development server (no embedded assets)"
	@echo "  dev-build           - Generate assets and run development server"
	@echo "  help                - Show this help message"