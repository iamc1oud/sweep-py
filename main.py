"""Find deletable build/dependency folders, chosen by the file types next to them.

    python main.py [root]             # dry run, lists folders + sizes
    python main.py [root] --delete    # actually remove them
"""
import argparse
import fnmatch
import os
import shutil
from pathlib import Path

# marker file glob (must sit in the same dir) -> junk folder globs it implies
RULES = {
    "package.json": ["node_modules", ".next", ".nuxt", ".turbo", ".parcel-cache"],
    "pyproject.toml": [".venv", "venv", "__pycache__", ".pytest_cache", ".mypy_cache", ".ruff_cache", "*.egg-info"],
    "requirements.txt": [".venv", "venv", "__pycache__", ".pytest_cache", ".mypy_cache", ".ruff_cache"],
    "setup.py": [".venv", "venv", "__pycache__", "*.egg-info", "build", "dist"],
    "Cargo.toml": ["target"],
    "pom.xml": ["target"],
    "build.gradle*": ["build", ".gradle"],
    "composer.json": ["vendor"],
    "*.csproj": ["bin", "obj"],
    "pubspec.yaml": [".dart_tool", "build"],
    "Package.swift": [".build"],
    "Gemfile": [".bundle"],
}
SKIP = {".git"}


def find_junk(root):
    """Yield junk dirs under root. Matched dirs are not descended into."""
    for cur, dirs, files in os.walk(root):
        dirs[:] = [d for d in dirs if d not in SKIP]
        patterns = {p for m, ps in RULES.items() if fnmatch.filter(files, m) for p in ps}
        for d in list(dirs):
            path = Path(cur, d)
            if not path.is_symlink() and any(fnmatch.fnmatch(d, p) for p in patterns):
                dirs.remove(d)
                yield path


def size(path):
    return sum(f.stat().st_size for f in path.rglob("*") if f.is_file() and not f.is_symlink())


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("root", nargs="?", default=".")
    ap.add_argument("index", nargs="?", default=None, type=int)
    ap.add_argument("--delete", action="store_true")
    args = ap.parse_args()

    # If root is a digit string, treat it as index in current dir
    if args.root.isdigit():
        args.index = int(args.root)
        args.root = "."

    junk = [(path, size(path)) for path in find_junk(args.root)]
    total = sum(s for _, s in junk)

    # Delete specific index or all
    if args.index is not None:
        if 0 <= args.index < len(junk):
            path, s = junk[args.index]
            if args.delete:
                shutil.rmtree(path)
                print(f"deleted {s / 1e6:.1f} MB  {path}")
            else:
                print(f"would free {s / 1e6:.1f} MB  {path}")
        else:
            print(f"Index {args.index} out of range (0-{len(junk)-1})" if junk else "No junk folders found")
    else:
        # List all with indices
        for idx, (path, s) in enumerate(junk):
            print(f"[{idx}]  {s / 1e6:10.1f} MB  {path}")
            if args.delete:
                shutil.rmtree(path)
        print(f"{'deleted' if args.delete else 'would free'} {total / 1e6:.1f} MB")


if __name__ == "__main__":
    main()
