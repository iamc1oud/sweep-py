# sweep

Find and delete build/dependency folders (`node_modules`, `.venv`, `target`, ...) based on the project files next to them.

A folder is only flagged if its marker file sits in the same directory. A `target/` next to `Cargo.toml` is junk; a `target/` on its own is left alone.

## Install

One line (installs [uv](https://docs.astral.sh/uv/) if missing; temp download/build cache is deleted afterwards, only `sweep` stays):

```
curl -LsSf https://raw.githubusercontent.com/iamc1oud/sweep-py/main/install.sh | sh
```

Or from a clone (Python 3.13+ and uv required):

```
uv tool install --editable .
```

Uninstall: `uv tool uninstall sweep`

## Usage

```
sweep [root]            # dry run: list folders and sizes (default root: .)
sweep [root] --delete   # remove them
```

Nothing is deleted without `--delete`, and there is no confirmation prompt when you pass it. Run the dry run first.

## What it detects

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

To add a language, add an entry to `RULES` in `main.py`.
