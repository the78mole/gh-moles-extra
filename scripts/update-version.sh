#!/bin/bash
# Script to update manifest.yml version from git tag

# Get the latest git tag, fallback to 0.0.0 if no tags exist
VERSION=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.0.0")

# Update manifest.yml
if [ -f "manifest.yml" ]; then
    sed -i "s/version: .*/version: ${VERSION}/" manifest.yml
    echo "Updated manifest.yml to version ${VERSION}"
else
    echo "manifest.yml not found"
    exit 1
fi