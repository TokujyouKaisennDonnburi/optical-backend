#!/usr/bin/env bash

set -e

if [ ! -d "tmp" ]; then
    mkdir "tmp"
fi

opencode run --agent=review "review $(git diff $(git merge-base HEAD origin/main))" > tmp/ai-review.md

less tmp/ai-review.md
