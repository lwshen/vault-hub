# Official OpenAPI Generator – Preparation

This directory contains the scaffolding required to swap `oapi-codegen` with the
official [OpenAPI Generator CLI](https://openapi-generator.tech/) during
**Phase 4** of the migration. The goal is to keep the existing `go generate`
workflow working today while we stage the new tooling in parallel.

## Files

- `config-go-server.yaml` – configuration used when generating the Go server
  stubs and models. It does **not** write into the current `generated.go` so we
  can inspect outputs safely in a disposable directory.
- `config-go-client.yaml` – placeholder for the Go client configuration. The CLI
  is currently used by the published `vault-hub-go-client` module and will be
  wired after the server swap is complete.

## Usage

```bash
# From the repository root
go generate ./packages/api/...
```

`go generate` invokes `packages/api/generate-openapi.sh`, which bundles the spec
and runs the official CLI to produce preview artifacts under
`.openapi-generator/`. The script can also be called directly if you need to
avoid updating other packages:

```bash
cd packages/api
sh generate-openapi.sh
```

Both commands intentionally write to temporary directories so we can iterate
without overwriting the Fiber-based artifacts.

## Next Steps

1. Finalize the generator options in both config files once we validate Echo
   integration and strict handler requirements.
2. Replace the existing `go:generate` directives with the CLI invocation and
   remove the `oapi-codegen` tool reference from `go.mod`.
3. Propagate the updated artifacts to the CLI/cron consumers and adjust CI to
   rely on the official generator.
