# action-semantic-versioning

A GitHub Action that automatically bumps semantic versions based on [Conventional Commits](https://www.conventionalcommits.org/) and creates git tags and releases.

## How It Works

1. Finds the latest semver tag in the repository
2. Parses all commits since that tag using Conventional Commits format
3. Determines the version bump:
   - `fix:` or `perf:` → **patch** (1.2.3 → 1.2.4)
   - `feat:` → **minor** (1.2.3 → 1.3.0)
   - `BREAKING CHANGE:` footer or `!` after type → **major** (1.2.3 → 2.0.0)
4. Creates a new git tag and optionally a GitHub release with changelog

## Usage

### Basic (tag only)

```yaml
name: Version
on:
  push:
    branches: [main]

jobs:
  version:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Required: full history for tag detection

      - uses: netwarlan/action-semantic-versioning@v1
        id: version

      - name: Print version
        if: steps.version.outputs.skipped != 'true'
        run: echo "New version: ${{ steps.version.outputs.new-version }}"
```

### With GitHub Release

```yaml
      - uses: netwarlan/action-semantic-versioning@v1
        id: version
        with:
          create-release: 'true'
```

### Dry Run

```yaml
      - uses: netwarlan/action-semantic-versioning@v1
        id: version
        with:
          dry-run: 'true'

      - run: echo "Would bump to ${{ steps.version.outputs.new-version }}"
```

### Gate Downstream Jobs

```yaml
      - uses: netwarlan/action-semantic-versioning@v1
        id: version

      - name: Build and publish
        if: steps.version.outputs.skipped != 'true'
        run: ./build.sh
```

## Inputs

| Input | Default | Description |
|-------|---------|-------------|
| `token` | `${{ github.token }}` | GitHub token for pushing tags and creating releases |
| `default-version` | `v0.1.0` | Starting version when no existing tags are found |
| `tag-prefix` | `v` | Tag prefix |
| `create-release` | `false` | Create a GitHub release with changelog |
| `release-draft` | `false` | Create the release as a draft |
| `release-prerelease` | `false` | Mark the release as a prerelease |
| `bump-patch-on-unknown` | `false` | Bump patch for non-conventional commits (docs, chore, etc.) |
| `dry-run` | `false` | Calculate version without creating tag or release |

## Outputs

| Output | Description |
|--------|-------------|
| `previous-version` | The previous semver tag found |
| `new-version` | The new calculated version |
| `bump-type` | The bump type applied: `major`, `minor`, `patch`, or `none` |
| `changelog` | Generated changelog markdown |
| `skipped` | `true` if no version bump occurred, `false` otherwise |

## Commit Message Format

This action follows [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope][!]: <description>

[optional body]

[optional footer(s)]
```

**Examples:**

```
fix: correct timezone handling in scheduler
feat(api): add user endpoint
feat!: redesign authentication flow
docs: update API reference

fix: handle null pointer in parser

The parser previously crashed when encountering null values
in the configuration file.

BREAKING CHANGE: config format changed from YAML to TOML
```

## Requirements

- **`fetch-depth: 0`** on `actions/checkout` — the action needs full git history to find tags and read commits
- **`permissions: contents: write`** on the job — required to push tags and create releases
