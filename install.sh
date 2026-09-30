#!/bin/sh
# curl -LsSf https://raw.githubusercontent.com/iamc1oud/sweep-py/main/install.sh | sh
set -eu

SRC="${SWEEP_SRC:-https://github.com/iamc1oud/sweep-py/archive/refs/heads/main.tar.gz}"

if ! command -v uv >/dev/null 2>&1; then
    curl -LsSf https://astral.sh/uv/install.sh | sh
    PATH="$HOME/.local/bin:$PATH"
fi

# throwaway cache: downloaded + build deps are deleted on exit, only the tool stays
UV_CACHE_DIR="$(mktemp -d)"
export UV_CACHE_DIR UV_LINK_MODE=copy
trap 'rm -rf "$UV_CACHE_DIR"' EXIT

uv tool install --force "$SRC"

echo "installed: sweep (run 'uv tool update-shell' if it's not on your PATH)"
