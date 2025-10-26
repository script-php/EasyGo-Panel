# EasyGo Panel

A comprehensive web server management panel written in Go, similar to Hestia CP.

## Features

- **Dual Interface**: CLI and Web panel
- **Single Binary**: All components embedded
- **PAM Authentication**: System user integration
- **Service Management**: Apache, Nginx, PHP-FPM, DNS, Mail, Databases
- **SSL Management**: Let's Encrypt with auto-renewal
- **Security**: Firewall, Fail2ban, IP lists
- **Backup & Cron**: Automated backup and task scheduling

## Build

```bash
go build -o easygo cmd/easygo/main.go
```

## Usage

### Web Panel
```bash
./easygo web
```

### CLI Commands
```bash
./easygo help
./easygo apache install
./easygo php install 8.2
./easygo ssl create example.com
```

## Requirements

- Linux OS
- Root privileges for system management
- PAM development libraries

## Architecture

- `cmd/` - Main application entry point
- `internal/` - Internal packages (cli, web)
- `pkg/` - Shared packages (actions, auth)
- `web/` - Static assets and templates