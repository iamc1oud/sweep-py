# sweep - Go Implementation

Compiled Go version of the sweep utility. Faster runtime, no interpreter dependency.

## Usage

```bash
./sweep [root]             # dry run, lists folders + sizes
./sweep [root] --delete    # actually remove them
```

## Examples

```bash
# List all junk in current directory
./sweep

# List all junk in /home/user/projects
./sweep /home/user/projects

# Delete specific index
./sweep . 2 --delete

# Delete all
./sweep . --delete
```

## Building

```bash
go build -o sweep main.go
```

Creates a standalone binary (~2.6 MB, no dependencies).

### Install Globally

```bash
go build -o sweep main.go
sudo mv sweep /usr/local/bin/
```

Then use: `sweep [root] [--delete]` from anywhere

## Adding New Project Types

Edit the `rules` map in `main.go`:

```go
var rules = map[string][]string{
    "marker_file_pattern": {"junk_folder_pattern1", "junk_folder_pattern2"},
}
```

Example:
```go
"go.mod":      {"vendor"},  // Golang projects
```

## Performance

- Fast: Compiled to native code
- Single pass through directory tree
- Memory efficient
- ~2.6 MB standalone binary

## Requirements

- Go 1.23+ (for building)
- Any OS: Linux, macOS, Windows

## Comparison with Python

| Aspect | Python | Go |
|--------|--------|-----|
| Dependencies | None | Go compiler (build only) |
| Runtime | Requires Python 3.13+ | Standalone binary |
| Speed | Fast | Faster |
| Binary Size | ~3KB script | ~2.6 MB binary |
| Development | Easier | Faster execution |

Both versions have identical behavior and output format.
