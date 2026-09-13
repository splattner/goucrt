# Vendored Core-API spec

This directory is a partial, read-only vendor copy of
[unfoldedcircle/core-api](https://github.com/unfoldedcircle/core-api), pinned to the commit recorded in
[`VENDORED_COMMIT.txt`](VENDORED_COMMIT.txt).

It exists so goucrt's tests can check its hand-maintained protocol constants (entity features, device classes,
commands, attributes) against the spec's own machine-readable schema and worked examples, instead of relying on
someone noticing drift by eye. See `internal/spec` for the code that reads it and `pkg/entities/*_spec_test.go`
for the tests that use it.

## Contents

- `integration-api/UCR-integration-asyncapi.yaml` — the integration driver AsyncAPI definition. Entity `features`
  and `device_class` are enumerated here as JSON Schema `enum`s and are what `internal/spec` parses.
- `doc/entities/*.md` — the per-entity-type reference docs. `cmd_id`, attribute names and state values are **not**
  present as structured schema in the YAML; they only exist in these docs' prose tables and worked JSON examples,
  which is what the golden fixtures under `pkg/entities/testdata/` were extracted from.

## Updating

1. Re-fetch the two paths above from core-api's `main` branch.
2. Update `VENDORED_COMMIT.txt` with the new commit hash/date.
3. Re-run `go run ./tools/extractgoldens` (see that package's doc comment) to refresh
   `pkg/entities/testdata/`.
4. Run `go test ./...` — spec-verification test failures point at exactly what changed (a renamed feature, a
   new/removed command, etc.) and need a human decision about whether/how to follow suit in goucrt's own types.

## License

core-api is licensed under [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/) by Unfolded Circle ApS.
This vendored copy is unmodified reference material used only to drive goucrt's own tests at build/test time; it
is not part of goucrt's compiled output and does not change goucrt's own license (see the repository root
`LICENSE`).
