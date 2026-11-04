#!/bin/sh
# curl -LsSf https://raw.githubusercontent.com/iamc1oud/sweep-py/main/install.sh | sh
set -eu

REPO="${SWEEP_REPO:-https://github.com/iamc1oud/sweep-py/archive/refs/heads/main.tar.gz}"
VERSION="${SWEEP_VERSION:-python}"  # python or go

install_python() {
    if ! command -v uv >/dev/null 2>&1; then
        curl -LsSf https://astral.sh/uv/install.sh | sh
        PATH="$HOME/.local/bin:$PATH"
    fi

    UV_CACHE_DIR="$(mktemp -d)"
    export UV_CACHE_DIR UV_LINK_MODE=copy
    trap 'rm -rf "$UV_CACHE_DIR"' EXIT

    uv tool install --force "$REPO"
    echo "installed: sweep (Python) - run 'uv tool update-shell' if it's not on your PATH"
}

install_go() {
    if ! command -v go >/dev/null 2>&1; then
        echo "error: Go not found. Install from https://go.dev/dl"
        exit 1
    fi

    # Download and extract
    TMP_DIR="$(mktemp -d)"
    trap 'rm -rf "$TMP_DIR"' EXIT

    cd "$TMP_DIR"
    curl -LsSf "$REPO" | tar xz
    cd sweep-py-main/go

    go build -o sweep main.go

    # Install to $GOPATH/bin or $HOME/.local/bin
    INSTALL_DIR="${GOPATH:-$HOME/go}/bin"
    mkdir -p "$INSTALL_DIR"
    mv sweep "$INSTALL_DIR/"

    echo "installed: sweep (Go) at $INSTALL_DIR/sweep"
    echo "Add $INSTALL_DIR to \$PATH if not already there"
}

case "$VERSION" in
    python)
        install_python
        ;;
    go)
        install_go
        ;;
    *)
        echo "error: SWEEP_VERSION must be 'python' or 'go', got '$VERSION'"
        exit 1
        ;;
esac
