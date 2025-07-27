#!/bin/bash
cd /Users/viacheslavezhkin/projects/gh-cp-2
git checkout initial
git branch -D cherry-pick/add-readme-documentation 2>/dev/null || true
exec /Users/viacheslavezhkin/projects/gh-cp/gh-cp 1 initial --dry-run