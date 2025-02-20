# Create Coding Feature

Genererate code change to an existing coding project

## Usage

```bash
python3 patterns/create_coding_feature/code.py . "Create a simple Hello World C program in file main.c" | fabric --pattern create_coding_feature
```

## Example

### <u>Input</u>:
```bash
python3 patterns/create_coding_feature/code.py . "Ensure that all user input is validated and sanitized before being used in the program." | fabric --pattern create_coding_feature
git diff
make check;
git add <changed files>
git commit -s -m "Security fixes: Input validation"
```
### <u>Output</u>:
PROJECT:

Let the AI model apply code changes to the existing project. Before command is issues, its the responsibility of the user to have a clean git repo, so that the user can diff and review the changes.

SUMMARY:

AI Autonoumous code file edits and creation in code project.

STEPS:

1. Input directory and file contents to AI model.
2. Generate code changes.
3. Let fabric interpret AI output via the file manager API, which will directly create and modify files and directories.
4. User reviews code and approves changes in SCM of choice.

SUGGESTIONS:

- Enhance script generation with conditional logic.
- Include detailed logging for API responses.
- Consider adding a GUI for ease of use.