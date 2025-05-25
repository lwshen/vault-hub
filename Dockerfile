ARG NODE_VERSION=22
ARG GO_VERSION=1.24

FROM node:${NODE_VERSION}-alpine AS frontend-builder

ARG GITHUB_TOKEN

WORKDIR /app

COPY web ./

RUN corepack enable

RUN if [ -n "$GITHUB_TOKEN" ]; then \
        echo "//npm.pkg.github.com/:_authToken=$GITHUB_TOKEN" >> .npmrc; \
    fi && \
    pnpm install --frozen-lockfile

RUN pnpm build

FROM golang:${GO_VERSION}-alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache gcc libc-dev

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux

COPY . .

COPY --from=frontend-builder /app/dist ./web/dist

RUN go mod download

RUN go build -o vault-hub-server cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=backend-builder /app/vault-hub-server ./
COPY --from=backend-builder /app/web/dist ./web/dist

EXPOSE 3000

CMD ["sh", "-c", "./vault-hub-server"]
