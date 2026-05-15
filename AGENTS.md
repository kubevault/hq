# AGENTS.md

This file provides guidance to coding agents (e.g. Claude Code, claude.ai/code) when working with code in this repository.

## Repository purpose

Go module `github.com/kubevault/hq` — a tiny "**jq for HCL**" CLI. Converts JSON on stdin to HCL on stdout (default), or HCL on stdin to JSON on stdout (with `--reverse`). Used by KubeVault tooling to massage Vault HCL policy files.

The produced binary is `hq`. There is **no Makefile** — this is a one-file utility you build with `go build`.

## Architecture

- `main.go` — the entire program. Two functions: `toHCL()` (default; JSON stdin → HCL stdout) and `toJSON()` (HCL stdin → JSON stdout). Versioning via `-version` flag, backed by `Version` ldflag.
- `hack/testdata/` — input fixtures for ad-hoc testing (no automated tests are configured).
- `go.mod` — `module github.com/kubevault/hq` (note: uses the GitHub URL, not a vanity URL). Single non-stdlib dep: `github.com/hashicorp/hcl`.

## Common commands

There is no Makefile.

- Build:

  ```
  go build -ldflags "-X main.Version=$(git describe --tags --always)" -o bin/hq
  ```

- Convert JSON → HCL:

  ```
  cat policy.json | hq > policy.hcl
  ```

- Convert HCL → JSON:

  ```
  cat policy.hcl | hq --reverse > policy.json
  ```

- Version:

  ```
  hq --version
  ```

## Conventions

- Module path is `github.com/kubevault/hq` (the GitHub URL, **not** a `kubevault.dev` vanity URL). Imports must use that.
- License: see `LICENSE`.
- Sign off commits (`git commit -s`).
- Keep this binary tiny — the whole point is "single-file utility, minimal deps". Resist the urge to refactor into packages or pull in a CLI framework. If you need flags beyond `--reverse` / `--version`, prefer Go's `flag` package over Cobra.
- The version string is plumbed via `-ldflags "-X main.Version=..."`; CI / release scripts rely on that. Do not rename `Version`.
