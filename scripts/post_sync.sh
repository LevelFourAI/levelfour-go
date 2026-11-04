#!/usr/bin/env bash
# Everything a regenerated tree needs to pass CI. The sync runs this once after it replaces the
# generated files, so the pull request it opens is already final and nothing pushes to it later.
set -euo pipefail
cd "$(dirname "$0")/.."

./scripts/fix_headers.sh
./scripts/gen_aliases.sh
