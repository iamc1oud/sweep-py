# Sweep - Go Port

A complete port of the Python `sweep` utility to Go, maintaining 100% feature parity and behavior compatibility.

## Overview

`sweep` is a utility that finds and deletes build/dependency folders (`node_modules`, `.venv`, `target`, etc.) based on marker files (project configuration files) in the same directory.

## Features

✅ **Marker-based detection** - Only deletes folders when their marker file exists nearby (e.g., `.venv` next to `pyproject.toml`)

✅ **Smart traversal** - Doesn't descend into already-matched junk folders

✅ **Dry-run mode** - Preview what would be deleted before using `--delete`

✅ **Flexible arguments** - Accepts arguments in any order (root path, index, --delete flag)

✅ **Size calculation** - Shows how much space would be freed

✅ **Multiple formats** - List all, delete specific index, or delete all at once

## Building

```bash
go build -o sweep main.go
```

## Usage

```bash
sweep [root]            # dry run: list folders and sizes (default root: .)
sweep [root] --delete   # remove them
sweep [index]           # show specific junk folder (uses current dir as root)
sweep [index] --delete  # delete specific junk folder
```

### Examples

```bash
# List all junk folders in current directory
./sweep

# List all junk folders starting from /home/user/projects
./sweep /home/user/projects

# Delete a specific junk folder (index 2)
./sweep 2 --delete

# Delete all junk folders
./sweep . --delete
```

## Supported Project Types

| Marker File | Folders Detected | Reason |
|---|---|---|
| `package.json` | `node_modules`, `.next`, `.nuxt`, `.turbo`, `.parcel-cache` | Node.js/JavaScript |
| `pyproject.toml` | `.venv`, `venv`, `__pycache__`, `.pytest_cache`, `.mypy_cache`, `.ruff_cache`, `*.egg-info` | Python (modern) |
| `requirements.txt` | `.venv`, `venv`, `__pycache__`, `.pytest_cache`, `.mypy_cache`, `.ruff_cache` | Python (legacy) |
| `setup.py` | `.venv`, `venv`, `__pycache__`, `*.egg-info`, `build`, `dist` | Python (setuptools) |
| `Cargo.toml` | `target` | Rust |
| `pom.xml` | `target` | Java/Maven |
| `build.gradle*` | `build`, `.gradle` | Java/Gradle |
| `composer.json` | `vendor` | PHP |
| `*.csproj` | `bin`, `obj` | C# / .NET |
| `pubspec.yaml` | `.dart_tool`, `build` | Dart/Flutter |
| `Package.swift` | `.build` | Swift |
| `Gemfile` | `.bundle` | Ruby |

## Implementation Notes

### Key Differences from Python Version

1. **Argument Parsing**: Go's standard `flag` package requires flags before positional arguments, so this port implements custom argument parsing to match Python's flexible `argparse` behavior

2. **Directory Traversal**: Uses a custom recursive walk function instead of Python's `os.walk` to properly handle not descending into matched junk directories

3. **Glob Matching**: Uses Go's `filepath.Match` for pattern matching (supports `*` and `?` wildcards, compatible with Python's `fnmatch`)

4. **Binary Size**: The compiled Go binary (~2.6 MB) is significantly larger than the Python script, but offers faster runtime and no interpreter dependency

### Performance

- Go version is faster due to compiled code (no interpreter overhead)
- Single-pass directory traversal with linear time complexity
- Memory efficient: builds list of junk folders then processes them

## Testing

The port has been verified against the following scenarios:

✅ Detects Node.js project junk folders
✅ Detects Python project junk folders  
✅ Detects Rust project junk folders
✅ Dry-run listing works correctly
✅ Index-based selection works
✅ Deletion with --delete flag works
✅ Flexible argument order (root, index, --delete in any order)
✅ Symlink detection and skipping
✅ `.git` folder skipping
✅ Size calculation accuracy

## Behavior Parity Verification

| Feature | Python | Go | Status |
|---------|--------|----|----|
| Marker file detection | ✅ | ✅ | ✅ Identical |
| Junk folder matching | ✅ | ✅ | ✅ Identical |
| Directory traversal | ✅ | ✅ | ✅ Identical |
| Dry-run listing | ✅ | ✅ | ✅ Identical |
| Index-based operations | ✅ | ✅ | ✅ Identical |
| Batch deletion | ✅ | ✅ | ✅ Identical |
| Size calculation | ✅ | ✅ | ✅ Identical |
| Symlink handling | ✅ | ✅ | ✅ Identical |
| Error handling | ✅ | ✅ | ✅ Identical |

## File Structure

```
sweep-py/
├── main.py          # Original Python implementation
├── main.go          # Go port
├── go.mod          # Go module definition
├── sweep            # Compiled binary
└── GOLANG_PORT.md   # This file
```

## Future Enhancements

These features are planned in both versions:

- [ ] Custom configuration file support
- [ ] TUI for interactive folder selection
- [ ] Color output for better readability
- [ ] Progress bar for large directory trees

## Contributing

To keep the Go and Python versions in sync, any changes to one should be reflected in the other (unless language-specific).
