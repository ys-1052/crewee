#!/bin/bash

# Setup git hooks for the project

set -e

echo "Setting up git hooks..."

# Create scripts directory if it doesn't exist
mkdir -p .git/hooks

# Create pre-commit hook
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash

# Pre-commit hook for Go formatting and linting

set -e

echo "Running pre-commit checks..."

# Check if we're in the backend directory context
if [ -d "backend" ]; then
    # Format Go code
    echo "Formatting Go code..."
    (cd backend && gofmt -w .)
    (cd backend && goimports -w -local github.com/ytakahashi/crewee .)
    
    # Run linter
    echo "Running golangci-lint..."
    (cd backend && golangci-lint run)
    
    # Add formatted files back to staging
    git add backend/
fi

# Check if we're in frontend directory context
if [ -d "frontend" ]; then
    echo "Checking frontend code..."
    # Frontend checks will be added later
fi

echo "Pre-commit checks completed successfully!"
EOF

# Make the hook executable
chmod +x .git/hooks/pre-commit

echo "Git hooks setup completed!"
echo "Pre-commit hook will now:"
echo "  - Format Go code with gofmt and goimports"
echo "  - Run golangci-lint"
echo "  - Add formatted files back to staging"