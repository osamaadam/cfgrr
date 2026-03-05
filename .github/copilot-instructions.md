# Copilot Instructions for `cfgrr`

## Repo Snapshot

- `cfgrr` is a Go CLI (Cobra + Viper) for backing up, restoring, replicating, and syncing user config files.
- Core behavior depends on path handling, symlink safety, and map-file consistency (YAML/JSON).
- Primary command wiring lives in `cmd/`, core file operations in `configfile/` and `core/`.

## What Good Changes Look Like

- Keep changes small and focused; avoid broad refactors unless requested.
- Preserve command behavior and flags unless the request explicitly changes UX.
- Prefer explicit, actionable errors (with context) over generic failures.
- Follow existing Go style and package boundaries already used in the repo.

## Testing and Validation Rules

- Add or update tests when behavior changes.
- Prefer table-driven tests for command and file-operation logic.
- Use `t.TempDir()` and isolated test state for filesystem tests.
- Validate both success paths and failure paths for filesystem/path logic.

## Build Safety (Required After Any Code Change)

After making changes, always run these validations and only finish when they pass:

```sh
go test ./... -v
go build ./...
```

If a full test run is too heavy during iteration, run targeted tests first, but run both commands above before finalizing.

## File and Platform Safety Checks

- Treat `backup_dir` and user paths as platform-sensitive inputs.
- Avoid assuming Linux-only home paths (for example `/home/<user>`) when writing logic.
- Use `filepath` utilities for path joins and absolute/relative handling.
- Never silently delete or overwrite user files outside the expected backup/restore flow.

## Skills to Apply in This Repo

Use this lightweight checklist as execution skills:

1. Path-Portability Skill:
   Normalize and validate path inputs early; error clearly when config paths are invalid on the current OS.
2. Command-Contract Skill:
   Keep Cobra command arguments, help text, and side effects aligned with existing command contracts.
3. Filesystem-Test Skill:
   Cover symlink, permission, and map-file side effects with deterministic temp-dir tests.
4. Build-Gate Skill:
   Always end with a green `go test ./... -v` and `go build ./...`.

## PR/Change Summary Expectations

When summarizing changes:

- List behavior changes first.
- List tests added/updated.
- Include build/test commands run and their outcome.
