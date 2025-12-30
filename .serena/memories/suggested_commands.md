# Suggested Commands

## Build & Deploy

- **Build Go Binary**:

  ```bash
  ./build.sh
  ```

  _Compiles the Go app for Linux/ARM64 and packages it into `function.zip`._

- **Deploy Infrastructure**:

  ```bash
  cd awscdk
  npm run cdk deploy
  ```

## Testing

- **Run Go Tests**:

  ```bash
  go test -v -race -shuffle on ./app/...
  ```

- **Run CDK Tests**:

  ```bash
  cd awscdk
  npm test
  ```

## Development Utilities

- **Generate Wire Code** (Dependency Injection):

  ```bash
  cd cmd
  wire
  ```

- **Lint Go Code**:

  ```bash
  golangci-lint run
  ```

  _(Requires `golangci-lint` installed)_

- **Install CDK Dependencies**:

  ```bash
  cd awscdk
  npm ci
  ```
