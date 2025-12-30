# Style and Conventions

## Architecture

- **Clean Architecture**: The project is structured into layers:
  - `model`: Domain entities.
  - `usecase`: Business logic interfaces.
  - `interactor`: Business logic implementation.
  - `external`: Adapters for external services (Steam, Notion, Discord).
  - `cmd`: Application entry point.

## Go Conventions

- **Dependency Injection**: Uses `google/wire` to manage dependencies. Update `cmd/wire.go` and run `wire` when adding new dependencies.
- **Mocking**: Uses `go.uber.org/mock` for generating mocks in tests.
- **Linting**: strictly enforced via `.golangci.yaml`.
  - Linters enabled: `revive`, `staticcheck`, `gosec`, `govet`, etc.
- **Formatting**: Uses `gofumpt` and `gci` (imports organization).

## Infrastructure (CDK)

- **Language**: TypeScript.
- **Testing**: Jest snapshots and unit tests.
