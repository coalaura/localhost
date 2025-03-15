# Localhost

`localhost` is a simple, fast, and user-friendly web server designed to enhance your local web development experience. Whether you're serving static files, testing PHP applications, or working with Laravel, `localhost` simplifies the setup process with convenient features like live reload and easy HTTPS configuration.

## What it Does

- **Static File Serving**: Quickly serve any local files and directories in your browser.
- **Live Reload**: Automatically refreshes your browser whenever you save changes, saving you time and boosting productivity.
- **PHP and Laravel Ready**: Effortlessly run PHP projects or Laravel applications with minimal setup.
- **Secure Connections**: Optional built-in HTTPS support, with easy redirection from HTTP.
- **Automatic Browser Launch**: Optionally opens your default browser when the server starts.

## Getting Started

### Installation

Ensure Go is installed on your system, then run:

```bash
go build -o localhost
```

## Quick Start

Serve your current directory:

```bash
./localhost
```

Serve a specific directory:

```bash
./localhost -d path/to/your/directory
```

## Advanced Options

Customize your experience with these handy flags:

| Flag              | Short | Description                                    | Default                    |
|-------------------|-------|------------------------------------------------|----------------------------|
| `--port`          | `-p`  | Specify port number to run the server.         | `80` (HTTP), `443` (HTTPS) |
| `--directory`     | `-d`  | Specify the document root directory.           | Current directory (`pwd`)  |
| `--live`          | `-l`  | Enable automatic page reloads on file changes. | Enabled                    |
| `--cert`          | `-c`  | SSL certificate file (enables HTTPS).          | Disabled                   |
| `--key`           | `-k`  | SSL key file (required for HTTPS).             | Disabled                   |
| `--redirect`      | `-r`  | Redirect HTTP to HTTPS.                        | Disabled                   |
| `--open`          | `-o`  | Open browser automatically on server start.    | Disabled                   |
| `--verbose`       | `-v`  | Enable verbose PHP logging.                    | Disabled                   |

### Examples

**Run with HTTPS:**
```bash
./localhost --cert cert.pem --key key.pem
```

**Auto-launch Browser:**
```bash
./localhost --open
```

## Live Reload
When enabled, `localhost` automatically refreshes your browser whenever you modify files, ensuring your changes are instantly visible without manual refreshes.
