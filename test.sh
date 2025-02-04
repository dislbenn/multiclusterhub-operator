#!/bin/bash

# Get the branches and store them in a variable
BRANCHES=$(git branch -r | grep 'origin/release-' | sed 's/origin\///' | tr '\n' ',' | sed 's/,$//')

# Output the branches to a GitHub Actions environment file
echo "branches=$BRANCHES"
