# Contributing to StackPulse

Thank you for your interest in contributing to StackPulse! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/sanchukanirupama/stackpulse.git
   cd stackpulse
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Build the project:
   ```bash
   make build
   ```

5. Run tests:
   ```bash
   make test
   ```

## Development Workflow

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
make install
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make format

# Lint code (requires golangci-lint)
make lint
```

## Project Structure

```
stackpulse/
├── cmd/                 # CLI commands
│   ├── root.go         # Root command
│   ├── watch.go        # Watch command
│   └── status.go       # Status command
├── internal/           # Internal packages
│   ├── alerts/         # Alert management
│   ├── config/         # Configuration
│   ├── display/        # Terminal dashboard
│   ├── metrics/        # Metrics collection
│   ├── monitor/        # Main monitoring logic
│   └── types/          # Type definitions
├── scripts/            # Build and release scripts
├── Makefile           # Build automation
├── go.mod             # Go module definition
└── main.go            # Application entry point
```

## Coding Standards

### Go Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors appropriately


## Pull Request Process

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and commit them with clear messages

3. Add or update tests as needed

4. Ensure all tests pass:
   ```bash
   make test
   ```

5. Update documentation if needed

6. Push your branch and create a pull request

7. Address any feedback from code review


## Issue Reporting

When reporting issues, please include:
- Operating system and version
- Steps to reproduce the issue
- Expected vs actual behavior

## Feature Requests

For feature requests, please:

- Check existing issues to avoid duplicates
- Clearly describe the use case

Thank you for contributing to StackPulse!