---
name: server-linter
description: Running and fixing linting issues for the Go server using Makefile commands. Use this skill when the user mentions linting, CI failures related to server-lint, or when finishing a task to ensure code quality.
---

# Server Linter Skill

This skill provides a standardized workflow for maintaining Go code quality in the `server/` directory using `golangci-lint` via the project's `Makefile`.

## Workflow

### 1. Check for Issues
Always start by running the standard lint check to identify current issues:
```bash
make server-lint
```

### 2. Automatic Fixing
If issues are found, attempt to fix them automatically using the built-in fix command:
```bash
make server-lint-fix
```

### 3. Manual Resolution
If `make server-lint` still reports issues after an automatic fix, resolve them manually following these project standards:
- **funlen**: If a function is too long, extract logical blocks into smaller helper functions.
- **mnd (magic numbers)**: Replace raw numbers with descriptive constants.
- **wrapcheck**: For Fiber controller methods (`.JSON`, `.SendStatus`, etc.), use `//nolint:wrapcheck` if the error is simply being passed back to the Fiber context.
- **wsl (whitespace)**: Ensure proper spacing between declarations and error checks as per `golangci-lint` requirements.
- **nlreturn**: Ensure there is a blank line before `return` statements if required.

### 4. Final Validation
Always run `make server-lint` one last time to confirm all issues are resolved. A task is not complete until the linter returns zero issues.
