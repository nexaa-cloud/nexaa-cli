#!/bin/bash
set -e

echo "🔍 Running local linting checks..."

# Check go fmt
echo "📝 Checking go fmt..."
unformatted=$(find . -name "*.go" -not -path "./vendor/*" | xargs gofmt -l)
if [ -z "$unformatted" ]; then
    echo "✅ go fmt passed"
else
    echo "❌ go fmt failed. These files need formatting:"
    echo "$unformatted"
    echo "Run: go fmt ./..."
    exit 1
fi

# Check go mod tidy
echo "📦 Checking go mod tidy..."
# Check if go.mod and go.sum are clean
if ! git diff --exit-code go.mod go.sum >/dev/null 2>&1; then
    echo "❌ go.mod/go.sum have uncommitted changes. Commit them first."
    exit 1
fi

# Save current state and run go mod tidy
cp go.mod go.mod.orig
cp go.sum go.sum.orig 2>/dev/null || touch go.sum.orig
go mod tidy

# Check if files changed
if ! cmp -s go.mod go.mod.orig || ! cmp -s go.sum go.sum.orig; then
    echo "❌ go mod tidy needed. These files need updating:"
    echo "Run: go mod tidy && git add go.mod go.sum"
    # Restore original files
    mv go.mod.orig go.mod
    mv go.sum.orig go.sum
    exit 1
fi

# Clean up
rm go.mod.orig go.sum.orig
echo "✅ go mod tidy passed"

# Run golangci-lint via Docker
echo "🔧 Running golangci-lint..."
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.63.1 golangci-lint run --timeout=3m

echo "🎉 All linting checks passed!"