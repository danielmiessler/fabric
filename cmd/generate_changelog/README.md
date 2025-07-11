# Changelog Generator

A high-performance changelog generator for Git repositories that automatically creates comprehensive, well-formatted changelogs from your git history and GitHub pull requests.

## Features

- **One-pass git history walking**: Efficiently processes entire repository history in a single pass
- **Automatic PR detection**: Extracts pull request information from merge commits
- **GitHub API integration**: Fetches detailed PR information including commits, authors, and descriptions
- **Smart caching**: SQLite-based caching for instant incremental updates
- **Unreleased changes**: Tracks all commits since the last release
- **Concurrent processing**: Parallel GitHub API calls for improved performance
- **Flexible output**: Generate complete changelogs or target specific versions
- **Optimized PR fetching**: Batch fetches all merged PRs using GitHub Search API (drastically reduces API calls)
- **Intelligent sync**: Automatically syncs new PRs every 24 hours or when missing PRs are detected

## Installation

```bash
go install github.com/danielmiessler/fabric/cmd/generate_changelog@latest
```

## Usage

### Basic usage (generate complete changelog)
```bash
generate_changelog
```

### Save to file
```bash
generate_changelog -o CHANGELOG.md
```

### Generate for specific version
```bash
generate_changelog -v v1.4.244
```

### Limit to recent versions
```bash
generate_changelog -l 10
```

### Using GitHub token for private repos or higher rate limits
```bash
export GITHUB_TOKEN=your_token_here
generate_changelog

# Or pass directly
generate_changelog --token your_token_here
```

### Cache management
```bash
# Rebuild cache from scratch
generate_changelog --rebuild-cache

# Force a full PR sync from GitHub
generate_changelog --force-pr-sync

# Disable cache usage
generate_changelog --no-cache

# Use custom cache location
generate_changelog --cache /path/to/cache.db
```

## Command Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--repo` | `-r` | Repository path | `.` (current directory) |
| `--output` | `-o` | Output file | stdout |
| `--limit` | `-l` | Limit number of versions | 0 (all) |
| `--version` | `-v` | Generate for specific version | |
| `--save-data` | | Save version data to JSON | false |
| `--cache` | | Cache database file | `./cmd/generate_changelog/changelog.db` |
| `--no-cache` | | Disable cache usage | false |
| `--rebuild-cache` | | Rebuild cache from scratch | false |
| `--force-pr-sync` | | Force a full PR sync from GitHub | false |
| `--token` | | GitHub API token | `$GITHUB_TOKEN` |

## Output Format

The generated changelog follows this structure:

```markdown
# Changelog

## Unreleased

### PR [#1601](url) by [author](profile): PR Title
- Change description 1
- Change description 2

### Direct commits
- Direct commit message 1
- Direct commit message 2

## v1.4.244 (2025-07-09)

### PR [#1598](url) by [author](profile): PR Title
- Change description
...
```

## How It Works

1. **Git History Walking**: The tool walks through your git history from newest to oldest commits
2. **Version Detection**: Identifies version bump commits (pattern: "Update version to vX.Y.Z")
3. **PR Extraction**: Detects merge commits and extracts PR numbers
4. **GitHub API Calls**: Fetches detailed PR information in parallel batches
5. **Change Extraction**: Extracts changes from PR commit messages or PR body
6. **Formatting**: Generates clean, organized markdown output

## Performance

- **Native Go libraries**: Uses go-git and go-github for maximum performance
- **Concurrent API calls**: Processes up to 10 GitHub API requests in parallel
- **Smart caching**: SQLite cache eliminates redundant API calls
- **Incremental updates**: Only processes new commits on subsequent runs
- **Batch PR fetching**: Uses GitHub Search API to fetch all merged PRs in minimal API calls

### Major Optimization: Batch PR Fetching

The tool has been optimized to drastically reduce GitHub API calls:

**Before**: Individual API calls for each PR (2 API calls per PR - one for PR details, one for commits)
- For a repo with 500 PRs: 1,000 API calls

**After**: Batch fetching using GitHub Search API
- For a repo with 500 PRs: ~10 API calls (search) + 500 API calls (details) = ~510 API calls
- **50% reduction in API calls!**

The optimization includes:
1. **Batch Search**: Uses GitHub's Search API to find all merged PRs in paginated batches
2. **Smart Caching**: Stores complete PR data and tracks last sync timestamp
3. **Incremental Sync**: Only fetches PRs merged after the last sync
4. **Automatic Refresh**: PRs are synced every 24 hours or when missing PRs are detected
5. **Fallback Support**: If batch fetch fails, falls back to individual PR fetching

## Requirements

- Go 1.24+ (for installation from source)
- Git repository
- GitHub token (optional, for private repos or higher rate limits)

## Authentication

The tool supports GitHub authentication via:
1. Environment variable: `export GITHUB_TOKEN=your_token`
2. Command line flag: `--token your_token`

Without authentication, the tool is limited to 60 GitHub API requests per hour.

## Caching

The SQLite cache stores:
- Version information and commit associations
- Pull request details (title, body, commits, authors)
- Last processed commit SHA for incremental updates
- Last PR sync timestamp for intelligent refresh

Cache benefits:
- Instant changelog regeneration
- Drastically reduced GitHub API usage (50%+ reduction)
- Offline changelog generation (after initial cache build)
- Automatic PR data refresh every 24 hours
- Batch database transactions for better performance

## Contributing

This tool is part of the Fabric project. Contributions are welcome!

## License

Same as the Fabric project.