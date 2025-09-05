#!/bin/sh

set -e

# Parse arguments
dry_run=false
increment=patch

while [ $# -gt 0 ]; do
  case "$1" in
    --dry-run|-n)
      dry_run=true
      shift
      ;;
    patch|minor|major)
      increment="$1"
      shift
      ;;
    *)
      echo "Usage: $0 [--dry-run|-n] [patch|minor|major]"
      echo "  --dry-run, -n: Show what would be done without making changes"
      echo "  increment: patch (default), minor, or major"
      exit 1
      ;;
  esac
done

if ! command -v uvx >/dev/null 2>&1; then
    echo "Error: uvx is required but not installed"
    exit 1
fi

# Get new version using bump-my-version
version=$(uvx bump-my-version show --increment "$increment" new_version)
version=v$version

if git rev-parse "refs/tags/$version" >/dev/null 2>&1
then
  echo "tag $version exists"
  exit 1
fi

if [ "$dry_run" = true ]; then
  echo "[DRY RUN] Would create tag: $version"
  echo "[DRY RUN] Would run: git tag -am \"$version\" \"$version\""
  echo "[DRY RUN] Would run: git push origin \"$version\""
else
  echo "About to create and push tag: $version"
  printf "Proceed? (y/N): "
  read -r response
  case "$response" in
    [yY]|[yY][eE][sS])
      git tag -am "$version" "$version"
      echo "tag $version created"
      git push origin "$version"
      echo "tag $version pushed to origin"
      ;;
    *)
      echo "Aborted"
      exit 1
      ;;
  esac
fi
