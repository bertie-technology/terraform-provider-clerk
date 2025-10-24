#!/bin/bash
# Setup git hooks for semantic commit messages and branch names

set -e

HOOKS_DIR=".githooks"
GIT_HOOKS_DIR=".git/hooks"

echo "🔧 Setting up git hooks for semantic checks..."

# Check if .git directory exists
if [ ! -d ".git" ]; then
    echo "❌ Error: Not a git repository"
    exit 1
fi

# Create .git/hooks directory if it doesn't exist
mkdir -p "$GIT_HOOKS_DIR"

# Copy commit-msg hook
if [ -f "$HOOKS_DIR/commit-msg" ]; then
    cp "$HOOKS_DIR/commit-msg" "$GIT_HOOKS_DIR/commit-msg"
    chmod +x "$GIT_HOOKS_DIR/commit-msg"
    echo "✅ Installed commit-msg hook"
else
    echo "⚠️  Warning: commit-msg hook not found in $HOOKS_DIR"
fi

# Alternatively, configure git to use the .githooks directory
git config core.hooksPath "$HOOKS_DIR"
echo "✅ Configured git to use $HOOKS_DIR directory"

echo ""
echo "✅ Git hooks setup complete!"
echo ""
echo "The following hooks are now active:"
echo "  - commit-msg: Validates commit messages follow Conventional Commits"
echo ""
echo "To bypass hooks (not recommended), use: git commit --no-verify"
