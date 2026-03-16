---
sidebar_position: 2
title: Installation
---

# Installation

jwx is available on macOS, Linux, and Windows. Choose the method that works best for your platform.

## macOS

### Homebrew (recommended)

The easiest way to install jwx on macOS is via [Homebrew](https://brew.sh/):

```bash
brew install manimovassagh/tap/jwx
```

This installs the latest release and makes `jwx` available in your PATH immediately. To upgrade later:

```bash
brew upgrade jwx
```

### Manual download

Download the binary from the [releases page](https://github.com/manimovassagh/jwx/releases):

```bash
# Apple Silicon (M1/M2/M3/M4)
curl -sL https://github.com/manimovassagh/jwx/releases/latest/download/jwx_darwin_arm64.tar.gz | tar xz
sudo mv jwx /usr/local/bin/

# Intel
curl -sL https://github.com/manimovassagh/jwx/releases/latest/download/jwx_darwin_amd64.tar.gz | tar xz
sudo mv jwx /usr/local/bin/
```

## Linux

### Binary download

Download the pre-built binary for your architecture:

```bash
# x86_64 (most common)
curl -sL https://github.com/manimovassagh/jwx/releases/latest/download/jwx_linux_amd64.tar.gz | tar xz
sudo mv jwx /usr/local/bin/

# ARM64 (Raspberry Pi 4+, AWS Graviton, etc.)
curl -sL https://github.com/manimovassagh/jwx/releases/latest/download/jwx_linux_arm64.tar.gz | tar xz
sudo mv jwx /usr/local/bin/
```

### APT repository (Debian/Ubuntu)

If an APT repository is available for your distribution:

```bash
# Add the repository (check the releases page for the latest instructions)
echo "deb [trusted=yes] https://manimovassagh.github.io/jwx/apt ./" | sudo tee /etc/apt/sources.list.d/jwx.list
sudo apt update
sudo apt install jwx
```

### Verify installation

```bash
jwx version
```

## Windows

### Chocolatey

If you use [Chocolatey](https://chocolatey.org/):

```powershell
choco install jwx
```

### Manual download

Download and extract the Windows binary:

```powershell
# Download the latest release
Invoke-WebRequest -Uri "https://github.com/manimovassagh/jwx/releases/latest/download/jwx_windows_amd64.zip" -OutFile "$env:TEMP\jwx.zip"

# Extract to a local directory
Expand-Archive "$env:TEMP\jwx.zip" -DestinationPath "$env:LOCALAPPDATA\jwx" -Force

# Add to PATH for the current session
$env:PATH += ";$env:LOCALAPPDATA\jwx"
```

To make it permanent, add `%LOCALAPPDATA%\jwx` to your system PATH through **Settings > System > About > Advanced system settings > Environment Variables**.

## Go install

If you have [Go](https://go.dev/dl/) 1.23 or later installed:

```bash
go install github.com/manimovassagh/jwx/cmd/jwx@latest
```

Make sure `$GOPATH/bin` (or `$HOME/go/bin`) is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Add this line to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) to make it permanent.

## Build from source

Clone the repository and build with Make:

```bash
git clone https://github.com/manimovassagh/jwx.git
cd jwx
make build
```

The binary is placed in `bin/jwx`. To install it to your GOPATH:

```bash
make install
```

### Build requirements

- **Go 1.23+** -- verify with `go version`
- **Make** -- available on most Unix systems; on Windows, use `go build` directly

### Building without Make

If you don't have Make, you can build directly with Go:

```bash
go build -o jwx ./cmd/jwx
```

## Verifying your installation

After installation, verify everything works:

```bash
# Check the version
jwx version

# Decode a sample token
jwx eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

If you see colorized output with the decoded header and payload, you're all set.

## Troubleshooting

### `jwx: command not found`

Make sure the binary is in a directory listed in your PATH. For Go installs:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Permission denied on Linux/macOS

If you get a permission error when running the binary:

```bash
chmod +x /usr/local/bin/jwx
```

### macOS Gatekeeper warning

On macOS, you may see a "cannot be opened because the developer cannot be verified" warning. To allow it:

1. Open **System Settings > Privacy & Security**
2. Click **Allow Anyway** next to the jwx message
3. Run `jwx` again and click **Open**

Or bypass Gatekeeper from the command line:

```bash
xattr -d com.apple.quarantine /usr/local/bin/jwx
```
