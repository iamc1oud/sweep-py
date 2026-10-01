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
    ap =    argparsese.ArgumentParser()
    ac
    ap.add_argument("root", nargs="?", default=".")
    ap.add_argument("--delete", action="store_true")
    args = ap.parse_args()

    total = 0
    for path in find_junk(args.root):
        s = size(path)
        total += s
        print(f"{s / 1e6:10.1f} MB  {path}")
        if args.delete:
            shutil.rmtree(path)
    print(f"{'deleted' if args.delete else 'would free'} {total / 1e6:.1f} MB")


if __name__ == "__main__":
    main()
