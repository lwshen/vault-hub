# Repository Guidelines

## Project Structure & Module Organization
- `apps/server` hosts the Go Fiber API; routes in `apps/server/handler`, middleware in `apps/server/route`, and shared config helpers in `apps/server/internal`.
- `apps/cli` provides the Cobra CLI backed by `internal/cli` logic and `internal/encryption` utilities.
- `apps/web` contains the Vite + React UI (`src/pages`, `src/components`, `src/stores`); run UI assets through `pnpm`.
- Do not edit files under `apps/web`; that directory is managed as an external codebase.
- `apps/cron` and `scripts/` supply scheduled jobs and release chores; keep them idempotent.
- Shared OpenAPI specs live in `packages/api`; regenerate clients with `go generate packages/api/tool.go`.
- Reusable models reside in `model/`; container assets live under `docker/`.

## Build, Test, and Development Commands
- `go run ./apps/server/main.go` launches the API at http://localhost:3000 for local dev.
- `go build -o tmp/main ./apps/server/main.go` and `go build -o vault-hub-cli ./apps/cli/main.go` compile server and CLI binaries; add release `-ldflags` when tagging.
- `JWT_SECRET=test ENCRYPTION_KEY=$(openssl rand -base64 32) go test ./...` runs backend unit tests.
- `golangci-lint run ./...` enforces Go lint rules; fix before commits.
- `pnpm --dir apps/web install && pnpm --dir apps/web run dev` starts the web app; run `pnpm --dir apps/web run build`, `pnpm --dir apps/web run lint`, and `pnpm --dir apps/web run typecheck` prior to merging.
- Use `air -c .air.toml` for backend live reload; the Air config automatically rebuilds the embedded frontend assets.

## Coding Style & Naming Conventions
- Format Go code with `gofmt`; exported types use PascalCase, internal helpers remain unexported.
- CLI command files adopt hyphenated filenames (`list.go`) and snake_case flags.
- React components use PascalCase; hooks/utilities use camelCase; apply Tailwind classes inline and keep global CSS minimal.
- Commit generated artifacts only when necessary; regenerate `packages/api` outputs after spec edits.

## Testing Guidelines
- Place Go tests in `*_test.go` next to their implementations; cover config, encryption, and database flows.
- Ensure secrets used in tests are ephemeral (`openssl rand`); never commit real credentials or `data.db`.
- Add Vitest + Testing Library with a `pnpm run test` script when introducing UI coverage.

## Commit & Pull Request Guidelines
- Follow Conventional Commits (`feat:`, `fix:`, `chore(scope):`); scope optional but recommended for clarity.
- Rebase onto `main` and verify `.github/workflows/ci.yml` remains green before opening a PR.
- PRs should summarize scope, note schema or env changes, link tracking issues, and include CLI output or screenshots for UI updates.

## Security & Configuration Tips
- Keep `JWT_SECRET`, `ENCRYPTION_KEY`, and database credentials in environment variables or secret storage.
- Document new OIDC/database variables in PRs, and scrub sensitive rows from shared `data.db` snapshots.

## Post-change Checklist
- When modifying `packages/api/openapi/api.yaml`, bump the patch segment of the `info.version` field before regenerating artifacts unless this branch has already updated the version relative to `main`.
- After editing files in `packages/api`, do not manually modify `packages/api/api.bundled.yaml` or `packages/api/generated.go`; run `go generate packages/api/tool.go` to regenerate them.
- Run `golangci-lint run ./...` after backend changes to confirm the Go codebase stays clean.
- Run `pnpm typecheck` and `pnpm lint` from `apps/web` after frontend changes to catch TypeScript and lint issues early.
