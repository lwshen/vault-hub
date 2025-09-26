# Production Dockerfile - expects pre-built binaries from CI
FROM alpine:3.22

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S vaultuser && \
    adduser -u 1001 -S vaultuser -G vaultuser

# Copy pre-built binary and frontend assets (expected to be built by CI)
# Use build arg to determine which binary to copy based on target platform
ARG TARGETARCH=amd64
COPY bin/vault-hub-server-linux-${TARGETARCH} ./vault-hub-server
COPY internal/embed/dist ./internal/embed/dist

# Change ownership of app directory
RUN chown -R vaultuser:vaultuser /app

# Switch to non-root user
USER vaultuser

EXPOSE 3000

CMD ["sh", "-c", "./vault-hub-server"]
