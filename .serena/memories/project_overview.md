# Project Overview

## Purpose

Steam Game Prices Notifier is a serverless application that helps users buy games at the best time. It syncs a user's Steam wishlist to a Notion Database and sends Discord notifications when the current price of a game is lower than or equal to its historical lowest price recorded in Notion.

## Architecture

- **Type**: Serverless (AWS Lambda).
- **Schedule**: Triggered daily (e.g., 18:00 JST) via EventBridge.
- **Language**: Go (1.25.5) for the Lambda function.
- **Infrastructure**: AWS CDK (TypeScript) for resource provisioning.

## Tech Stack

- **Go**: Core application logic.
  - `github.com/aws/aws-lambda-go`: Lambda handler.
  - `github.com/google/wire`: Dependency Injection.
  - `go.uber.org/mock`: Mocking for tests.
- **TypeScript**: AWS CDK infrastructure definitions.
- **External APIs**:
  - Steam Store API
  - Notion API
  - Discord Webhook

## Key Directories

- `app/`: Core Go application logic (following Clean Architecture).
  - `external/`: Clients for Discord, Notion, Steam.
  - `interactor/`, `usecase/`: Business logic.
  - `model/`: Domain entities.
- `awscdk/`: CDK infrastructure code.
- `cmd/`: Application entry point (`main.go`) and DI wiring (`wire.go`).
- `.github/`: CI/CD workflows.
