---
name: terraform-linter
description: Checking and formatting Terraform code using Makefile commands. Use this skill after making changes to the server/terraform directory to ensure consistent code style.
---

# Terraform Linter Skill

This skill provides a standardized workflow for maintaining Terraform code quality using `make terraform-lint` and `make terraform-lint-fix` via the project's `Makefile`.

## Workflow

### 1. Check for Formatting Issues
Always run the lint check to identify files that are not properly formatted:
```bash
make terraform-lint
```

### 2. Automatic Formatting
If issues are found, fix them automatically using the format command:
```bash
make terraform-lint-fix
```

### 3. Final Validation
Run `make terraform-lint` again to confirm that all files are correctly formatted. A task involving Terraform is not complete until the linter returns zero issues.
