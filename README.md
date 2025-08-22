# Nexaa CLI

The Nexaa CLI is a tool to manage cloud resources on the Nexaa Serverless
Platform, which allows you to deploy containers very quickly.

## Setup your environment

Make sure you have golang (>=1.23.0) running (either locally or in a container).

## Building the binary

Building the binary can simply be done by `go build -o nexaa-cli .`

## Environment Configuration

By default, the CLI connects to the production environment. You can override specific endpoints using environment variables:

- `NEXAA_GRAPHQL_URL` - GraphQL API endpoint (default: production endpoint)
- `NEXAA_KEYCLOAK_URL` - Keycloak authentication URL (default: production endpoint)  
- `NEXAA_TOKEN_FILE` - Token storage file location (default: `./auth.json`)

### Custom Environment Setup

You can override the default endpoints in two ways:

1. **Using a .env file** (recommended for development):
   ```bash
   cp .env.example .env
   # Edit .env and customize the endpoints as needed
   ```

2. **Setting environment variables directly**:
   ```bash
   export NEXAA_GRAPHQL_URL="https://your-custom-endpoint/graphql/platform"
   export NEXAA_KEYCLOAK_URL="https://your-custom-auth-endpoint"
   go run . <args>
   ```

The CLI automatically loads `.env` files if present, making local development configuration easier.

## GraphQL Code Generation

To run the GraphQL code generation after making changes to the `operations` directory, you can run `GO111MODULE=on go run -mod=mod github.com/Khan/genqlient` to generate the `generated.go` file in the api directory.

## Autocomplete
To enable shell completion, run:

For Bash:
    `source <(nexaa-cli completion bash)`

For Zsh:
    `source <(nexaa-cli completion zsh)`

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