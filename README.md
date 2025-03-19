# Mast - A Fast and Reliable File Downloader

Mast is a command-line file downloader written in Go that supports concurrent downloads, resumable downloads, and download queue management.

## Features

- Concurrent downloads with configurable workers
- Resumable downloads
- Download queue management
- Progress tracking
- Support for custom headers and cookies
- Configurable retry mechanism
- Clean and intuitive CLI interface

## Installation

### Quick Install (Recommended)

Download the latest release for your platform from the [releases page](https://github.com/bugass/mast/releases):

```bash
# For Linux (64-bit)
wget https://github.com/bugass/mast/releases/latest/download/mast-linux-amd64 -O mast
chmod +x mast
sudo mv mast /usr/local/bin/

# For Windows (64-bit)
# Download mast-windows-amd64.exe from releases
```

### From Source

If you want to build from source:

```bash
git clone https://github.com/bugass/mast.git
cd mast
go build -o mast
sudo mv mast /usr/local/bin/
```

## Usage

### Basic Download

```bash
# Download with default filename
mast download https://example.com/file.zip

# Download with custom filename
mast download https://example.com/file.zip -f custom.zip

# Download to specific location
mast download https://example.com/file.zip -l downloads/
```

### With Custom Headers

```bash
mast download https://example.com/file.zip --header "Authorization: Bearer token"
```

### With Cookies

```bash
mast download https://example.com/file.zip --cookie "session=abc123"
```

### Queue Management

List all downloads:
```bash
mast queue list
```

Check queue status:
```bash
mast queue status
```

Pause a download:
```bash
mast queue pause <task-id>
```

Resume a download:
```bash
mast queue resume <task-id>
```

Cancel a download:
```bash
mast queue cancel <task-id>
```

## Configuration

The following flags are available for the download command:

- `-f, --file`: Destination filename (default: filename from URL)
- `-l, --location`: Location to save the file (optional)
- `--cookie`: Cookies to send with the request (format: name=value)
- `--header`: Headers to send with the request (format: name:value)
- `--workers`: Number of download workers (default: 5)
- `--retries`: Maximum number of retry attempts (default: 3)
- `--resume`: Enable resumable downloads (default: true)

## License

MIT License - see LICENSE file for details 