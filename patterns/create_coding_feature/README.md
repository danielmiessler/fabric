# Create Coding Feature

Generate code change to an existing coding project

## Usage

After installing the `fabric_code` binary

```bash
go install github.com/danielmiessler/fabric/plugins/tools/fabric_code@latest
```

We can use it like this:

```bash
fabric_code . "Create a simple Hello World C program in file main.c" | fabric --pattern create_coding_feature
```

## Example Workflow

### Input

```bash
fabric_code . . "Ensure that all user input is validated and sanitized before being used in the program." | fabric --pattern create_coding_feature
git diff
make check;
git add <changed files>
git commit -s -m "Security fixes: Input validation"
```

### Expected outcome

This feature enables the AI model to apply code changes to the existing project.

Before using this feature, it iss the responsibility of the user to have a clean git repo, so that the user can diff and review the changes.

### Summary of the process

The AI working with the `fabric` tool automates file edits and creation in our project.

1. Input directory and file contents are given to AI model using a standardized json format.
2. AI to Generate code changes and produce file change/create annotations.
3. The `fabric` CLI will interpret AI output via the file manager API, which will directly create and modify files and directories.
4. User reviews code and approves changes in SCM of choice.

### Suggestions for future improvements

- Enhance script generation with conditional logic.
- Include detailed logging for API responses.
- Consider adding a GUI for ease of use.
