#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

ROOT_PKG="github.com/LevelFourAI/levelfour-go"
OUT="levelfour/aliases.go"

# Subpackages whose request/response types are part of the public surface and
# should be re-exported into the levelfour package for the single-import story.
# Service-client subpackages (apikeys, accounts, audit, auth, costs, health,
# providers, webhooks, recommendations/client, recommendations/audit) are
# accessed through the embedded fields on client.Client and intentionally
# excluded here.
SUBPKGS=(recommendations)

collect_types() {
    go doc -all "./$1" 2>/dev/null | grep -E '^type [A-Z]' | awk '{print $2}' | sort -u
}

collect_funcs() {
    go doc -all "./$1" 2>/dev/null | grep -E '^func [A-Z][a-zA-Z0-9]*\(' | sed 's/func //' | sed 's/(.*//' | sort -u
}

root_types=$(collect_types .)
root_funcs=$(collect_funcs .)

{
    echo "package levelfour"
    echo
    echo "import ("
    echo "	rootpkg \"$ROOT_PKG\""
    for pkg in "${SUBPKGS[@]}"; do
        alias_name="${pkg//\//_}pkg"
        echo "	$alias_name \"$ROOT_PKG/$pkg\""
    done
    echo ")"
    echo

    echo "// Type aliases re-exported so users only need a single levelfour import."
    echo "// Regenerate with: ./scripts/gen_aliases.sh"
    echo "type ("
    for t in $root_types; do
        echo "	$t = rootpkg.$t"
    done
    for pkg in "${SUBPKGS[@]}"; do
        alias_name="${pkg//\//_}pkg"
        for t in $(collect_types "$pkg"); do
            echo "	$t = $alias_name.$t"
        done
    done
    echo ")"
    echo

    echo "// Function aliases re-exported from the root and subpackages."
    echo "var ("
    for f in $root_funcs; do
        echo "	$f = rootpkg.$f"
    done
    for pkg in "${SUBPKGS[@]}"; do
        alias_name="${pkg//\//_}pkg"
        for f in $(collect_funcs "$pkg"); do
            echo "	$f = $alias_name.$f"
        done
    done
    echo ")"
} > "$OUT"

gofmt -w "$OUT"

total_types=$(echo "$root_types" | wc -w | tr -d ' ')
total_funcs=$(echo "$root_funcs" | wc -w | tr -d ' ')
for pkg in "${SUBPKGS[@]}"; do
    total_types=$((total_types + $(collect_types "$pkg" | wc -w | tr -d ' ')))
    total_funcs=$((total_funcs + $(collect_funcs "$pkg" | wc -w | tr -d ' ')))
done

echo "Generated $OUT with $total_types types and $total_funcs functions (root + ${#SUBPKGS[@]} subpackages)"
