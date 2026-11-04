# Project Structure

Modular mono-repo with separate Python and Go implementations.

## Directory Layout

```
sweep-py/
├── README.md                 # Main documentation
├── STRUCTURE.md             # This file
├── GOLANG_PORT.md           # Go port details & verification
├── pyproject.toml           # Python package config
├── install.sh               # Installation script
│
├── python/                  # Python implementation
│   ├── main.py              # Python entry point
│   └── README.md            # Python-specific docs
│
└── go/                      # Go implementation
    ├── main.go              # Go entry point
    ├── go.mod               # Go module definition
    ├── sweep                # Compiled binary (pre-built)
    └── README.md            # Go-specific docs
```

## Independence

Each version is **completely independent**:

- **Python** (`python/main.py`): Standalone script, no build required
- **Go** (`go/main.go`): Standalone module, compiles to binary

They share:
- Same logic and behavior
- Common project rules (marker files → junk folders)
- Same CLI interface
- Documentation references

## Why Separate?

1. **Clarity**: No ambiguity about which version you're using
2. **Maintenance**: Changes to one version don't affect the other's build
3. **Dependency Management**: Python uses `uv`, Go uses `go.mod` separately
4. **Flexibility**: Users can choose based on their environment
5. **CI/CD**: Each can have independent build pipelines

## Using the Versions

### Python

```bash
cd python
python main.py [root] [--delete]
```

Install globally:
```bash
cd ..
uv tool install --editable .
sweep [root] [--delete]
```

### Go

```bash
cd go
./sweep [root] [--delete]
```

Build:
```bash
cd go
go build -o sweep main.go
```

Install globally:
```bash
cd go
go build -o sweep main.go
sudo mv sweep /usr/local/bin/
```

## Keeping Versions in Sync

When adding features or fixing bugs:

1. **Edit both implementations**:
   - Python: `python/main.py` (edit `RULES` dict)
   - Go: `go/main.go` (edit `rules` map)

2. **Update corresponding README**:
   - Python: `python/README.md`
   - Go: `go/README.md`

3. **Test both versions**:
   ```bash
   cd python && python main.py /test/path
   cd ../go && ./sweep /test/path
   ```

4. **Update main README** if behavior changes:
   - `README.md` (top-level)
   - `GOLANG_PORT.md` (if Go-specific)

## Supported Rules

Both versions share identical rules (marker file → junk folders):

```python
# Python (python/main.py)
RULES = {
    "package.json": ["node_modules", ".next", ".nuxt", ".turbo", ".parcel-cache"],
    "pyproject.toml": [".venv", "venv", "__pycache__", ...],
    # etc.
}
```

```go
// Go (go/main.go)
var rules = map[string][]string{
    "package.json": {"node_modules", ".next", ".nuxt", ".turbo", ".parcel-cache"},
    "pyproject.toml": {".venv", "venv", "__pycache__", ...},
    // etc.
}
```

When adding a new project type, update **both**.

## Testing

Test both versions with the same test case:

```bash
# Create test structure
mkdir -p /tmp/test/{node_app,py_app,rust_app}
echo '{}' > /tmp/test/node_app/package.json
mkdir /tmp/test/node_app/node_modules

# Test Python
python python/main.py /tmp/test

# Test Go
go/sweep /tmp/test

# Output should be identical
```

## CI/CD Considerations

If setting up continuous integration:

- **Python tests**: Test with Python 3.13+, use `uv`
- **Go tests**: Test with Go 1.23+, build to all platforms
- **Integration tests**: Run both versions against same test cases
- **Documentation**: Verify links in both README files
