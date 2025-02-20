# Fabric Web UI Enhancement Contribution Setup

### One-Time Setup in VSCode

Before adding web UI enhancements:

1. Verify working branch:
```bash
git status

Ensures active branch is feature/web-ui-enhancements for isolated development

Check repository connections:
git fetch upstream
git status

Validates:

Clean working tree
Latest upstream code
Branch synchronization status
Pending changes if any

Verify remote configurations:
git remote -v

Shows all configured remotes:

origin: Your fork (jmd1010/fabricDM)
upstream: Original repo (danielmiessler/fabric)
backup: Full environment (jmd1010/fabric-backup)

This verification process:

Establishes clean development foundation
Confirms proper repository relationships
Enables isolated feature development
Maintains clear upgrade path
Preserves complete backup access


This documentation provides clear steps for initial VSCode setup while maintaining proper Git workflow and repository relationships.



## Repository Structure
- Fork: jmd1010/fabric-contrib
- Original: danielmiessler/fabric
- Branch: feature/web-ui-enhancements

## Remote Configuration
- origin: git@github.com:jmd1010/fabricDM.git
- upstream: https://github.com/danielmiessler/fabric.git

## Development Environment
### Shell Configurations
- Updated .zshrc and .bashrc
- Corrected paths for new repository structure
- Go environment variables
- Pattern directory mappings

### Web Development Setup
- Backend: fabric --serve
- Frontend: npm run dev in web directory
- Development URLs and ports

## Pull Request Workflow
1. Sync with upstream
2. Feature branch development
3. Code review preparation
4. PR submission guidelines

## Backup Management
### Repository Setup
- Backup: git@github.com:jmd1010/fabric-backup.git
- Contains complete development environment
- Excludes only sensitive files (.env)
- Private repository for safety

### Backup Workflow
1. Regular backup pushes:
```bash
git push backup feature/web-ui-enhancements

2. After significant changes:

git add .
git commit -m "Development checkpoint - [description]"
git push backup feature/web-ui-enhancements

3. Before major refactoring:

git tag backup-pre-refactor-[date]
git push backup --tags


Development vs Backup

Backup repo: Complete environment preservation

Feature branch: Clean commits for PR

Separate commit messages for backup vs PR

Tag significant development points

Regular synchronization with both repos

Recovery Procedures

1. From backup:

git clone git@github.com:jmd1010/fabric-backup.git
git checkout feature/web-ui-enhancements

2. Environment restoration:

Copy .env.example to .env
Configure local paths
Install dependencies
Verify development setup


This structure provides:

1. Clear separation of concerns
2. Detailed workflow procedures
3. Recovery instructions
4. Best practices for both backup and development
5. Easy reference for ongoing work

