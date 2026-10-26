## Description
<!-- Provide a brief description of the changes introduced by this PR. -->
<!-- Explain the problem you are solving or the feature you are adding. -->

## Related Issue(s)
<!-- Link any related issues (e.g., "Closes #123"). -->

## Type of Change
<!-- Check the appropriate box below. -->
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Refactoring (no functional changes)
- [ ] CI / Dependency update
- [ ] Documentation update

## Infrastructure Changes
<!-- Briefly describe any AWS CDK / infrastructure resources added, modified, or deleted (if applicable). -->

## Checklist
<!-- Verify that all the following items have been completed. -->
- [ ] I have run `golangci-lint run` to ensure the Go code passes linting (if applicable).
- [ ] I have run `go test -v -race -shuffle on ./app/...` to verify unit tests pass (if applicable).
- [ ] I have regenerated `wire_gen.go` (`wire`) and mocks (`mockgen`) if DI or interfaces changed (if applicable).
- [ ] I have run `npm test` in `awscdk/` to verify AWS CDK tests pass (if applicable).
- [ ] I have updated any relevant documentation (if applicable).
