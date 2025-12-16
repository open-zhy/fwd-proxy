#!/bin/sh
set -e

REPO="open-zhy/fwd-proxy"
BINARY_NAME="fwd-proxy"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux)
    ;;
  darwin)
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# Detect Architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)
    ARCH="amd64"
    ;;
  aarch64|arm64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported Architecture: $ARCH"
    exit 1
    ;;
esac

DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}-${OS}-${ARCH}"

echo "Downloading ${BINARY_NAME} for ${OS}/${ARCH}..."
echo "URL: ${DOWNLOAD_URL}"

# Download
curl -fsSL -o "${BINARY_NAME}" "${DOWNLOAD_URL}"
chmod +x "${BINARY_NAME}"

echo "Downloaded to ./${BINARY_NAME}"
echo "To install to /usr/local/bin, run:"
echo "  sudo mv ${BINARY_NAME} /usr/local/bin/"
