# FABRIC GitHub Workflow Guide

## Repository Setup

### Remote Repositories
- `origin` → jmd1010/FABRIC2 (your customized version)
- `upstream` → danielmiessler/fabric (original project)
- `fabricdm` → jmd1010/fabricDM (for your contributions)

### Branches
- `main` - Your working branch with all customizations
- `upstream-updates` - Tracks danielmiessler/fabric main branch

### Verify Setup
```bash
# Check remotes
git remote -v

# Should show:
origin    https://github.com/jmd1010/FABRIC2.git (fetch)
origin    https://github.com/jmd1010/FABRIC2.git (push)
upstream  https://github.com/danielmiessler/fabric.git (fetch)
upstream  https://github.com/danielmiessler/fabric.git (push)
fabricdm  https://github.com/jmd1010/fabricDM.git (fetch)
fabricdm  https://github.com/jmd1010/fabricDM.git (push)
```

## Daily Workflows

### 1. Trivial Changes
For small fixes or minor updates:
```bash
# Work directly on main branch
git add <files>
git commit -m "Description of changes"
git push origin main  # Backup to your repo
```

### 2. New Features/Risky Changes
For significant changes or experimental features:
```bash
# Create feature branch
git checkout -b feature/new-component main
# Work on changes...

# Regular commits
git add <files>
git commit -m "Feature progress"

# Backup branch to origin (optional)
git push origin feature/new-component

# When feature is stable:
git checkout main
git merge feature/new-component
git push origin main

# Clean up
git branch -d feature/new-component
git push origin --delete feature/new-component  # if you pushed it
```

### 3. Surgical Updates from Upstream

#### Get Specific File from Upstream
```bash
# Fetch latest from upstream
git fetch upstream

# Check what changed in specific file
git diff upstream/main -- path/to/file

# Option 1: Cherry-pick specific file
git checkout upstream/main -- path/to/file
# Then commit if you like the changes

# Option 2: Checkout file then cherry-pick specific changes
git checkout upstream/main -- path/to/file
git restore --staged path/to/file  # Unstage
# Edit file to keep only wanted changes
git add path/to/file
git commit -m "Cherry-picked specific changes from upstream"
```

#### Resolve Merge Conflicts
```bash
# If you get conflicts during merge
git status  # See conflicted files

# For each conflict, you have options:
# 1. Keep your version:
git checkout --ours path/to/file
git add path/to/file

# 2. Take their version:
git checkout --theirs path/to/file
git add path/to/file

# 3. Manually edit file to combine changes:
# Edit file, resolve markers, then:
git add path/to/file

# After resolving all conflicts:
git commit -m "Merged with conflict resolutions"
```

## Emergency Recovery

### Rebuild from Remote
If local repository becomes corrupted:

1. Backup your ENV file:
```bash
cp /Users/jmdb/.config/fabric/.env ~/env-backup
```

2. Remove corrupted local repo:
```bash
cd ..
mv FABRIC2 FABRIC2-corrupted
```

3. Clone fresh copy:
```bash
git clone https://github.com/jmd1010/FABRIC2.git
cd FABRIC2
```

4. Set up remotes:
```bash
git remote add upstream https://github.com/danielmiessler/fabric.git
git remote add fabricdm https://github.com/jmd1010/fabricDM.git
```

5. Create tracking branch:
```bash
git branch upstream-updates upstream/main
```

6. Restore ENV file:
```bash
mkdir -p /Users/jmdb/.config/fabric/
cp ~/env-backup /Users/jmdb/.config/fabric/.env
```

### Quick Recovery Steps
If you need to undo recent changes:

```bash
# Undo last commit but keep changes staged
git reset --soft HEAD^

# Undo last commit and remove changes
git reset --hard HEAD^

# Revert to specific commit
git reset --hard <commit-hash>

# Reset to remote state
git fetch origin
git reset --hard origin/main
```

## Contributing Back to danielmiessler/fabric

1. Create feature branch from upstream:
```bash
git fetch upstream
git checkout -b feature-name upstream/main
```

2. Make changes and commit:
```bash
git add <files>
git commit -m "Descriptive message"
```

3. Push to your fork:
```bash
git push fabricdm feature-name
```

4. Create PR:
- Go to https://github.com/danielmiessler/fabric
- Click "New Pull Request"
- Choose "compare across forks"
- Select your fork and branch
- Fill in PR description
- Submit

## Updating Existing PRs

### Sync fabricDM with Latest Changes
When you have new changes in FABRIC2 that you want to add to an existing PR:

1. Backup fabricDM's .gitignore:
```bash
# In fabricDM repo
cp fabric/.gitignore fabric/.gitignore.backup
```

2. Update fabricDM with your changes (one of two methods):

Method A - Complete sync:
```bash
git checkout feature/improvements  # or whatever your PR branch is
git fetch origin                  # get latest from FABRIC2
git reset --hard origin/main     # sync with your latest FABRIC2 state
cp fabric/.gitignore.backup fabric/.gitignore  # restore fabricDM's .gitignore
git add fabric/.gitignore
git commit -m "Restore fabricDM gitignore"
git push -f fabricdm feature/improvements
```

Method B - Merge approach:
```bash
git checkout feature/improvements
git pull origin main             # merge your latest FABRIC2 changes
git checkout --ours fabric/.gitignore  # keep fabricDM's version
git add fabric/.gitignore
git commit -m "Merge main while preserving fabricDM gitignore"
git push fabricdm feature/improvements
```

### Important Notes:
- Always protect fabricDM's .gitignore as it's specifically configured for contributions
- The backup step ensures you don't accidentally lose ENV file protection
- When in doubt, verify .gitignore content after any sync operation
- Consider Method B if you want to preserve PR commit history
- Use Method A for a clean slate that matches your FABRIC2 state exactly

## Best Practices

1. Always work in feature branches for significant changes
2. Regularly push to origin/main to backup your work
3. Keep commits atomic and well-described
4. Before starting new work:
   ```bash
   git status  # Ensure clean working directory
   git pull origin main  # Get latest changes
   ```
5. When in doubt, create a branch:
   ```bash
   git checkout -b backup/before-risky-change
   ```

## Additional Tips

1. Before major changes:
   - Create a backup branch
   - Push to origin
   - Document what you're about to do

2. When updating from upstream:
   - Always work in a branch first
   - Test thoroughly before merging to main
   - Keep notes of any manual conflict resolutions

3. Maintaining clean history:
   - Use meaningful commit messages
   - Group related changes in single commits
   - Push to origin main regularly

4. Emergency Situations:
   - Don't panic - your code is safe in origin/main
   - Create backup branches liberally
   - When in doubt, clone fresh and rebuild
