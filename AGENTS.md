# Repository Guidelines

## Project Structure & Module Organization
- `apps/server` – Go Fiber backend; routes in `handler/`, middleware in `route/`, config helpers in `internal/`.
- `apps/cli` – Cobra CLI reusing logic from `internal/cli` and `internal/encryption`.
- `apps/web` – Vite + React UI with screens in `src/pages`, reusable pieces in `components/`, and state in `stores/`.
- `apps/cron` & `scripts/` – Automation entrypoints for scheduled syncs and release chores.
- `packages/api` – OpenAPI specs and generated clients; update via `go generate packages/api/tool.go`.
- Shared models live in `model/` and container assets in `docker/`.

## Build, Test, and Development Commands
- `go run ./apps/server/main.go` starts the API at `http://localhost:3000`.
- `go build -o tmp/main ./apps/server/main.go` or `go build -o vault-hub-cli ./apps/cli/main.go` produces binaries; add README `-ldflags` for release metadata.
- `JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./...` runs backend tests.
- `cd apps/web && pnpm install && pnpm run dev` launches the frontend; run `pnpm run build`, `pnpm run lint`, and `pnpm run typecheck` before merge.
- `golangci-lint run ./...` and `pnpm run lint --fix` mirror CI checks.

## Coding Style & Naming Conventions
- Format Go with `gofmt`; keep exported symbols PascalCase and prefer unexported helpers for wiring.
- CLI command files use hyphenated names (e.g., `list.go`) with snake_case flags.
- React code uses PascalCase components, camelCase hooks/utilities, and Tailwind utility classes inline; keep global CSS minimal.

## Testing Guidelines
- Place Go tests beside implementations as `*_test.go`, covering config, encryption, and database flows.
- Add Vitest + Testing Library and a `pnpm run test` script when introducing UI coverage.
- Generate ephemeral secrets with `openssl rand` and avoid committing or reusing the sample `data.db`.

## Commit & Pull Request Guidelines
- Follow Conventional Commits (`feat:`, `fix:`, `chore(scope):`).
- Rebase on `main`, make sure `.github/workflows/ci.yml` passes, and regenerate `packages/api` outputs when specs change.
- PRs should summarise scope, note schema or env updates, link issues, and include CLI output or screenshots for UX changes.

## Security & Configuration Tips
- Store `JWT_SECRET`, `ENCRYPTION_KEY`, and database credentials outside the repo.
- Document new OIDC or database variables in PRs and purge sensitive rows from shared `data.db` snapshots.
