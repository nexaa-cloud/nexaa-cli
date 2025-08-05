# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the Tilaa/Nexaa CLI - a Go-based command-line tool for managing cloud resources on the Tilaa Serverless Platform. The CLI provides commands for containers, container jobs, namespaces, registries, volumes, and cloud database clusters.

## Development Commands

### Building
- `go build -o tilaa .` - Build the binary
- `make setup` - Set up vendor dependencies

### Testing
- `go test -race ./...` - Run all tests with race detection
- `make test` - Run tests via Makefile

### Linting
- `make lint` - Run golangci-lint via Docker

### Development Environment
- Use `TILAA_ENV=dev` environment variable to target staging instead of production
- `go run . <args>` - Run during development
- Authentication tokens are stored in `./auth.json` (dev) or `~/.tilaa/auth.json` (prod)

### Code Generation
- `go generate` - Generate GraphQL client code via genqlient (defined in main.go:32)

## Architecture

### Core Structure
- **main.go**: Entry point with environment configuration and config loading
- **cmd/**: Cobra CLI command definitions (root.go contains main CLI setup)
- **api/**: GraphQL client code and generated API bindings
- **config/**: Environment configuration and authentication token management
- **graphql/**: Custom GraphQL query builder and marshaling strategies  
- **operations/**: GraphQL operation definitions (.graphql files)

### Authentication Flow
- CLI requires login via `nexaa login` before most operations
- OAuth tokens stored in JSON config files with automatic expiration checking
- Environment-specific auth endpoints (dev vs prod)

### GraphQL Integration
- Uses genqlient for GraphQL code generation from schema.graphql
- Custom query builder in graphql/ package for dynamic query construction
- Operations defined in operations/*.graphql files

### Key Components
- **Cobra CLI framework** for command structure and shell completion
- **Environment switching** between dev/staging and production
- **Vendor directory** for dependency management
- **Custom GraphQL marshaling strategies** for different parameter types

### Resource Management
The CLI manages these primary resource types:
- Containers and Container Jobs
- Namespaces  
- Registries
- Volumes
- Cloud Database Clusters

Each resource type has its own command file in cmd/ and corresponding GraphQL operations.