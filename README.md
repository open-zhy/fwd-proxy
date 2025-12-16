# fwd-proxy

`fwd-proxy` is a lightweight, cross-platform HTTP Reverse Proxy written in Go. It is designed to be a simple, easy-to-deploy tool for forwarding requests to a backend server, with built-in support for CORS, custom headers, and cross-platform compatibility.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.24.3-blue)

## Features

- **HTTP/TCP Proxying**: Forwards HTTP requests to a specified target server.
- **CORS Support**: Easily enable default CORS headers with a single flag.
- **Custom Headers**: Inject custom response headers (e.g., for testing or specific client requirements).
- **Cross-Platform**: Binaries available for Linux, macOS, and Windows (amd64/arm64).
- **Zero Dependency**: Single static binary, no dependencies required on the host.

## Installation

### Quick Install (Linux/macOS)

You can install the latest release using the installation script:

```bash
curl -fsSL https://raw.githubusercontent.com/open-zhy/fwd-proxy/main/install.sh | sh
```

### Manual Download

Download the latest binary for your OS and Architecture from the [Releases](https://github.com/open-zhy/fwd-proxy/releases) page.

### From Source

Requirements: Go 1.24+

```bash
git clone https://github.com/open-zhy/fwd-proxy.git
cd fwd-proxy
make build
# Binary will be in bin/fwd-proxy
```

## Usage

```bash
./fwd-proxy -port <PORT> -target <TARGET_URL> [OPTIONS]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-port` | Local port to listen on. | `8080` |
| `-target` | **Required.** The target server URL (e.g., `http://localhost:9000`). | |
| `-cors` | Enable default CORS headers (`Access-Control-Allow-Origin: *`, etc.). | `false` |
| `-header` | Custom header to add to response. Format: `Key:Value`. Can be used multiple times. | |
| `-version`| Print version information and exit. | |

### Examples

**Basic Forwarding:**
Forward traffic from localhost:8080 to localhost:3000:
```bash
./fwd-proxy -target http://localhost:3000
```

**With Custom Port:**
Listen on port 9090 and forward to a remote API:
```bash
./fwd-proxy -port 9090 -target https://api.example.com
```

**Testing with CORS and Custom Headers:**
Useful for frontend development against an API that doesn't send CORS headers:
```bash
./fwd-proxy -target http://localhost:3000 -cors -header "X-Debug: True" -header "Cache-Control: no-cache"
```

## Development

### Running Tests
```bash
make test
```

### Building Cross-Platform Binaries
```bash
make build-all
```
This command generates binaries for all supported platforms in the `bin/` directory.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.
