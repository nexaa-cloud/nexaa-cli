# Tilaa CLI

The Tilaa CLI is a tool to manage cloud resources on the Tilaa Serverless
Platform, which allows you to deploy containers very quickly.

## Setup your environment

Make sure you have golang (>=1.23.0) running (either locally or in a container).

## Building the binary

Building the binary can simply be done by `go build -o tilaa .`

## Running development

If you're developing and want to talk to staging instead of production, use the environment variable TILAA_ENV=dev to talk to the dev environment

`go run . <args>`

This will ensure the code will talk to staging-auth and staging-graphql. Your
access token, refresh token, and ExpiresAt will be stored in JSON in
`./auth.json`. In the production version, this is stored at "~/.tilaa/auth.json"
which is better.

## Autocomplete
To enable shell completion, run:

For Bash:
    `source <(tilaa-cli completion bash)`

For Zsh:
    `source <(tilaa-cli completion zsh)`

Or to persist it, save the output to a file and source it in your shell config.

## Version Management

This project uses automated version management with semantic versioning. Versions are automatically incremented on pushes to main.

### Automatic Versioning
- **Patch** (default): Automatically increments on every push to main (e.g., `v1.0.0` → `v1.0.1`)
- **Minor**: Include `[minor]`, `[MINOR]`, `feat:`, or `feature:` in your commit message
- **Major**: Include `[major]`, `[MAJOR]`, or `BREAKING CHANGE` in your commit message

### Examples

**Patch version bump (automatic):**
```bash
git commit -m "Fix container status display bug"
git push origin main
# Creates v1.0.1
```

**Minor version bump:**
```bash
git commit -m "Add new registry management commands [minor]"
git push origin main
# Creates v1.1.0
```

**Major version bump:**
```bash
git commit -m "BREAKING CHANGE: Redesign CLI command structure"
git push origin main
# Creates v2.0.0
```

### Manual Version Control
You can also manually trigger version bumps via GitHub Actions:
1. Go to **Actions** → **Version Management**
2. Click **Run workflow**
3. Select version type (patch/minor/major)
4. Click **Run workflow**

Or via GitHub CLI:
```bash
gh workflow run version.yml -f version_type=major
```

### Releases
Once a version tag is created, the release workflow automatically:
- Builds binaries for Windows, Linux, and macOS (amd64 + arm64)
- Creates a GitHub release with installation instructions
- Generates checksums for verification