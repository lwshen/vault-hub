#!/bin/sh

set -e

version=$(git cliff --bumped-version)
if git rev-parse "refs/tags/$version" >/dev/null 2>&1
then
  echo "tag $version exists"
  exit 1
fi

git tag -am "$version" "$version"
echo "tag $version created"
git push origin "$version"
