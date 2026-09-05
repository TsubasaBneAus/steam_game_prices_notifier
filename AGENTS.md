# Project Context: Steam Game Prices Notifier

## Overview

This project is a serverless application that monitors price drops for games on a user's Steam wishlist.

- **Functionality**: Syncs Steam wishlist items to a Notion Database and sends Discord notifications when the current price is lower than or equal to the recorded lowest price.
- **Architecture**: AWS Lambda (Go) triggered by an EventBridge schedule (daily at 18:00 JST).
- **Infrastructure**: Managed via AWS CDK (TypeScript).

## Prerequisites

- **Go**: Version 1.25.5 or later.
- **Node.js & npm**: Required for AWS CDK.
- **AWS CLI**: Configured with appropriate credentials.
- **External Services**:
  - Notion Integration (API Key & Database ID).
  - Discord Webhook (ID & Token).
  - Steam Account (User ID).

## Setup & Configuration

1. **Environment Variables**:
   Create a `.env` file in the root directory based on `.env.example`.

   ```env
   NOTION_API_KEY="..."
   NOTION_DATABASE_ID="..."
   DISCORD_WEBHOOK_ID="..."
   DISCORD_WEBHOOK_TOKEN="..."
   STEAM_USER_ID="..."
   ```

2. **Infrastructure (AWS CDK)**:
   Located in the `awscdk/` directory.
   - Install dependencies: `cd awscdk && npm ci`

## Build & Deploy

### Go Application (Lambda)

The Go application is compiled and zipped for Lambda deployment.

- **Build Command**: `./build.sh`
  - This script compiles the code for `linux/arm64` and creates `function.zip`.

### Infrastructure

- **Deploy**:

  ```bash
  cd awscdk
  npm run cdk deploy
  ```

## Testing

### Go Application

Run unit tests for the application logic.

```bash
go test -v -race -shuffle on ./app/...
```

### AWS CDK

Run tests for the infrastructure definitions.

```bash
cd awscdk
npm test
```

## Project Structure

- `app/`: Core application logic (Clean Architecture).
  - `external/`: External API clients (Discord, Notion, Steam).
  - `usecase/`, `interactor/`: Business logic.
  - `model/`: Domain models.
  - `service/`: Interface definitions.
- `awscdk/`: AWS CDK infrastructure code (TypeScript).
- `cmd/`: Application entry point (`main.go`) and dependency injection wiring (`wire.go`).
- `build.sh`: Build script for the Lambda function.
- `.github/workflows/`: CI/CD pipelines (Test, Lint, Deploy).

## Development Conventions

- **Dependency Injection**: Uses `google/wire`. Run `wire` in `cmd/` or other directories if `wire.go` is modified.
- **Mocking**: Uses `go.uber.org/mock`.
- **Linting**:
  - Go: `golangci-lint` (configured in `.golangci.yaml`).
  - CDK: `eslint` (via `npm run lint` in `awscdk` if available, check package.json).
