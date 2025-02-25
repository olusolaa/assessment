# Product Catalog Sorting System

A flexible and extensible system for sorting product catalogs using Clean Architecture principles.

## Architecture

This project follows Clean Architecture principles with the following layers:

### Domain Layer
- Contains the core business logic and entities
- Defines interfaces that will be implemented by outer layers
- Located in the `domain` package

### Use Case Layer
- Contains application-specific logic
- Implements use cases that orchestrate domain entities
- Located in the `usecase` package

### Adapter Layer
- Translates between the use cases and external entities
- Implements interfaces defined in the domain layer
- Located in the `adapter` package

### Infrastructure Layer
- Contains configuration, persistence, and external services
- Implements interfaces defined in the domain layer
- Located in the `infrastructure` package

## Directory Structure

```
assessment/
├── domain/
│   ├── model/
│   │   └── product.go         # Core domain entities
│   ├── repository/
│   │   └── product_repo.go    # Repository interfaces
│   └── service/
│       └── sorter.go          # Core business logic interfaces
├── usecase/
│   ├── product_sorter.go      # Application use cases
│   └── product_paginator.go   # Pagination use cases
├── adapter/
│   ├── registry/
│   │   └── sorter_registry.go # Registry implementation
│   └── sorter/
│       ├── price_sorter.go    # Concrete sorter implementations
│       ├── date_sorter.go
│       └── ...
├── infrastructure/
│   ├── config/
│   │   └── config.go          # Configuration
│   └── persistence/
│       └── memory_repo.go     # In-memory repository implementation
└── cmd/
    └── main.go                # Application entry point
```

## Features

- **Extensible Sorting**: New sorting strategies can be added without modifying existing code
- **Configuration**: Sorters can be enabled/disabled via configuration
- **Pagination**: Support for paginating large result sets
- **Thread Safety**: All operations are thread-safe
- **Immutability**: Original data is never modified during sorting

## Available Sorters

- Price (ascending/descending)
- Creation Date (ascending/descending)
- Name (ascending/descending)
- Sales per View (ascending/descending)

## Usage

### Running the Application

```bash
go run cmd/main.go
```

### Adding a New Sorter

1. Create a new sorter in the `adapter/sorter` package:

```go
package sorter

import (
    "sort"
    "assessment/domain/model"
)

// MySorter sorts products by some criteria
type MySorter struct {
    ascending bool
}

// NewMySorter creates a new MySorter
func NewMySorter(ascending bool) *MySorter {
    return &MySorter{
        ascending: ascending,
    }
}

// Sort sorts products by some criteria
func (s *MySorter) Sort(products model.ProductList) model.ProductList {
    result := products.Clone()
    
    sort.Slice(result, func(i, j int) bool {
        // Your sorting logic here
        if s.ascending {
            return result[i].SomeField < result[j].SomeField
        }
        return result[i].SomeField > result[j].SomeField
    })
    
    return result
}

// Name returns the name of the sorter
func (s *MySorter) Name() string {
    if s.ascending {
        return "My Sorter (ascending)"
    }
    return "My Sorter (descending)"
}
```

2. Register the sorter in the registry:

```go
// In adapter/sorter/sorter_initializer.go
func InitializeDefaultSorters(registry service.SorterRegistry, cfg *config.Config) {
    // ... existing sorters ...
    
    // Register your new sorter
    registry.RegisterSorter(NewMySorter(true))
    registry.RegisterSorter(NewMySorter(false))
}
```

## Configuration

Configuration is stored in JSON format:

```json
{
  "disabled_sorters": [
    "Name (descending)"
  ],
  "default_page_size": 10
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## CI/CD Workflow

This project includes a CI/CD workflow using GitHub Actions that performs the following tasks:

### Build and Test
- Compiles the application
- Runs all tests
- Generates test coverage reports

### Linting
- Uses golangci-lint to check code quality
- Enforces code style and best practices

### Security Scanning
- Uses gosec to identify potential security issues
- Scans for common Go security anti-patterns

## Development

### Prerequisites
- Go 1.21 or higher
- Make (optional, for using the Makefile)

### Getting Started
1. Clone the repository
2. Install dependencies: `go mod download` or `make deps`
3. Build the application: `go build ./cmd/main.go` or `make build`
4. Run the application: `go run ./cmd/main.go` or `make run`

### Available Make Commands
```
make build     # Build the application
make test      # Run tests
make coverage  # Run tests with coverage
make lint      # Run linter
make run       # Run the application
make clean     # Clean build artifacts
make deps      # Install dependencies
make security  # Run security scan
make help      # Show help message
``` 