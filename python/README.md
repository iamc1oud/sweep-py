# sweep - Python Implementation

Original Python implementation of the sweep utility. Zero dependencies, Python 3.13+.

## Usage

```bash
python main.py [root]             # dry run, lists folders + sizes
python main.py [root] --delete    # actually remove them
```

## Examples

```bash
# List all junk in current directory
python main.py

# List all junk in /home/user/projects
python main.py /home/user/projects

# Delete specific index
python main.py . 2 --delete

# Delete all
python main.py . --delete
```

## Installation

### Via uv (recommended)
```bash
cd ..  # go to repo root
uv tool install --editable .
```

Then use: `sweep [root] [--delete]`

### Direct
```bash
python main.py [root] [--delete]
```

## Adding New Project Types

Edit the `RULES` dict in `main.py`:

```python
RULES = {
    "marker_file_pattern": ["junk_folder_pattern1", "junk_folder_pattern2"],
}
```

Example:
```python
"go.mod": ["vendor"],  # Golang projects
```

## Performance

- Fast: Single pass through directory tree
- No external dependencies
- Lightweight: ~3KB of code
