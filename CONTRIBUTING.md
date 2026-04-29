# Contributing to SEA Hub Backend API

Thank you for your interest in contributing to SEA Hub Backend API! This document provides guidelines and steps for contributing.

## Development Setup

```
git clone https://github.com/Engineering-Association-UofK/sea-api.git
cd sea-api
go mod download

# if running without docker:
cp .env.example .env # Fill with the needed fields
go run cmd/api/main.go

# if using docker
cp .env.compose.example .env.compose # Fill with the needed fields
docker-compose up --build
```

## Pull Request Process

1. Fork the repository and create your branch from `main`.
2. Make your changes.
3. Update documentation if needed.
4. Submit a pull request with a clear description of your changes.

## Reporting Issues

If you find a bug or have a feature request, please open an issue with:
- A clear title and description
- Steps to reproduce (for bugs)
- Expected vs actual behavior
