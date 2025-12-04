#!/usr/bin/env bash

set -e

if [ ! -d "tmp" ]; then
    mkdir "tmp"
fi

git fetch origin main
opencode run --agent=review "$(git diff origin/main)" > tmp/ai-review.md

less tmp/ai-review.md
