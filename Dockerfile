ARG NODE_VERSION=22
ARG GO_VERSION=1.24
ARG VERSION=dev
ARG COMMIT=unknown

FROM node:${NODE_VERSION}-alpine AS frontend-builder

WORKDIR /app

COPY apps/web ./

RUN corepack enable

RUN pnpm install --frozen-lockfile

RUN pnpm build

FROM golang:${GO_VERSION}-alpine AS backend-builder

ARG VERSION
ARG COMMIT

WORKDIR /app

RUN apk add --no-cache gcc libc-dev

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux

COPY . .

COPY --from=frontend-builder /app/apps/web/dist ./internal/embed/dist

RUN go mod download

RUN go build -ldflags="-X github.com/lwshen/vault-hub/internal/version.Version=${VERSION} -X github.com/lwshen/vault-hub/internal/version.Commit=${COMMIT}" -o vault-hub-server apps/server/main.go

FROM alpine:3.22

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S vaultuser && \
    adduser -u 1001 -S vaultuser -G vaultuser

COPY --from=backend-builder /app/vault-hub-server ./

# Change ownership of app directory
RUN chown -R vaultuser:vaultuser /app

# Switch to non-root user
USER vaultuser

EXPOSE 3000

CMD ["sh", "-c", "./vault-hub-server"]
