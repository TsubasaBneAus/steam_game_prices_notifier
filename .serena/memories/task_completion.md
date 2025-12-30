# Task Completion Checklist

After completing a task (feature or bug fix), ensure the following steps are performed:

1. **Run Go Tests**:

   ```bash
   go test -v -race -shuffle on ./app/...
   ```

   _Ensure all application logic tests pass._

2. **Lint Go Code**:

   ```bash
   golangci-lint run
   ```

   _Fix any linter errors._

3. **Run CDK Tests** (if infrastructure was modified):

   ```bash
   cd awscdk
   npm test
   ```

4. **Verify Build**:

   ```bash
   ./build.sh
   ```

   _Ensure the project compiles successfully._
