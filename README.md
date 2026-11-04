# sweep

Find and delete build/dependency folders (`node_modules`, `.venv`, `target`, ...) based on the project files next to them.

A folder is only flagged if its marker file sits in the same directory. A `target/` next to `Cargo.toml` is junk; a `target/` on its own is left alone.

## Implementations

This project contains two independent implementations:

- **[Python](./python/)** - Original implementation, zero dependencies
- **[Go](./go/)** - Compiled version, faster runtime, no interpreter needed

Choose either based on your preference. Both have identical behavior.

## Quick Start

### Python
```bash
cd python
python main.py [root]           # dry run
python main.py [root] --delete  # delete
```

### Go
```bash
cd go
./sweep [root]           # dry run
./sweep [root] --delete  # delete
```

Or build from source:
```bash
cd go
go build -o sweep main.go
./sweep
```

## Supported Project Types

| Marker file | Folders flagged |
|---|---|
| `package.json` | `node_modules`, `.next`, `.nuxt`, `.turbo`, `.parcel-cache` |
| `pyproject.toml` | `.venv`, `venv`, `__pycache__`, `.pytest_cache`, `.mypy_cache`, `.ruff_cache`, `*.egg-info` |
| `requirements.txt` | `.venv`, `venv`, `__pycache__`, `.pytest_cache`, `.mypy_cache`, `.ruff_cache` |
| `setup.py` | `.venv`, `venv`, `__pycache__`, `*.egg-info`, `build`, `dist` |
| `Cargo.toml` | `target` |
| `pom.xml` | `target` |
| `build.gradle*` | `build`, `.gradle` |
| `composer.json` | `vendor` |
| `*.csproj` | `bin`, `obj` |
| `pubspec.yaml` | `.dart_tool`, `build` |
| `Package.swift` | `.build` |
| `Gemfile` | `.bundle` |

Skipped: `.git`, symlinked folders, and anything inside an already-flagged folder. Go `vendor` is not touched because it is often committed.

## Usage Examples

```bash
# List all junk folders in current directory
sweep

# List all junk folders starting from /home/user/projects
sweep /home/user/projects

# Delete a specific junk folder (index 2)
sweep 2 --delete

# Delete all junk folders
sweep . --delete

# Dry run first, then delete
sweep
# review output
sweep . --delete
```

Nothing is deleted without `--delete`, and there is no confirmation prompt when you pass it. Run the dry run first.

## Installation

### Python (Default)

One line (installs [uv](https://docs.astral.sh/uv/) if missing):
```bash
curl -LsSf https://raw.githubusercontent.com/iamc1oud/sweep-py/main/install.sh | sh
```

Or from a clone (Python 3.13+ and uv required):
```bash
cd python
uv tool install --editable .
```

Uninstall: `uv tool uninstall sweep`

### Go

One line (installs Go binary):
```bash
SWEEP_VERSION=go curl -LsSf https://raw.githubusercontent.com/iamc1oud/sweep-py/main/install.sh | sh
```

Or build from source (Go 1.23+ required):
```bash
cd go
go build -o sweep main.go
sudo mv sweep /usr/local/bin/
```

Or download precompiled binary (when available).

### Install Options

Use environment variables with the universal installer:

```bash
# Python (default)
./install.sh

# Go
SWEEP_VERSION=go ./install.sh

# Custom repository (if forked)
SWEEP_REPO=https://github.com/user/fork/archive/refs/heads/main.tar.gz ./install.sh

# Go from custom repo
SWEEP_VERSION=go SWEEP_REPO=https://... ./install.sh
```

## Development

To add a language:
- **Python**: Add entry to `RULES` dict in `python/main.py`
- **Go**: Add entry to `rules` map in `go/main.go`

Keep both versions in sync.

## TODO

- [ ] Support user-defined targets for automatic deletion (custom config file / CLI flag)
- [ ] TUI for interactive folder selection and deletion
- [ ] Precompiled Go binaries for common platforms
- [ ] Color output for better readability
