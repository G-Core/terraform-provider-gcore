#!/bin/bash
#
# Compare two provider binaries
# Usage: ./compare_provider.sh <provider1> <provider2>
#

set -e

PROVIDER1="${1}"
PROVIDER2="${2}"

if [ -z "$PROVIDER1" ] || [ -z "$PROVIDER2" ]; then
    echo "Usage: $0 <provider1> <provider2>"
    echo ""
    echo "Example:"
    echo "  $0 ./terraform-provider-gcore ~/backup/terraform-provider-gcore"
    echo "  $0 ./terraform-provider-gcore /tmp/terraform-provider-gcore-machine2"
    exit 1
fi

echo "========================================"
echo "  Provider Binary Comparison"
echo "========================================"
echo ""

# Check if files exist
if [ ! -f "$PROVIDER1" ]; then
    echo "✗ ERROR: Provider 1 not found: $PROVIDER1"
    exit 1
fi

if [ ! -f "$PROVIDER2" ]; then
    echo "✗ ERROR: Provider 2 not found: $PROVIDER2"
    exit 1
fi

echo "Provider 1: $PROVIDER1"
echo "Provider 2: $PROVIDER2"
echo ""

# ============================================================================
# File Sizes
# ============================================================================
echo "File Sizes:"
echo "----------------------------------------"

SIZE1=$(ls -lh "$PROVIDER1" | awk '{print $5}')
SIZE2=$(ls -lh "$PROVIDER2" | awk '{print $5}')

SIZE1_BYTES=$(stat -f%z "$PROVIDER1" 2>/dev/null || stat -c%s "$PROVIDER1")
SIZE2_BYTES=$(stat -f%z "$PROVIDER2" 2>/dev/null || stat -c%s "$PROVIDER2")

echo "  Provider 1: $SIZE1 ($SIZE1_BYTES bytes)"
echo "  Provider 2: $SIZE2 ($SIZE2_BYTES bytes)"

if [ "$SIZE1_BYTES" -eq "$SIZE2_BYTES" ]; then
    echo "  ✓ Sizes match"
else
    echo "  ✗ Sizes differ!"
    DIFF_BYTES=$((SIZE1_BYTES - SIZE2_BYTES))
    echo "    Difference: $DIFF_BYTES bytes"
fi

echo ""

# ============================================================================
# SHA256 Hashes
# ============================================================================
echo "SHA256 Hashes:"
echo "----------------------------------------"

SHA256_1=$(shasum -a 256 "$PROVIDER1" | awk '{print $1}')
SHA256_2=$(shasum -a 256 "$PROVIDER2" | awk '{print $1}')

echo "  Provider 1: $SHA256_1"
echo "  Provider 2: $SHA256_2"

if [ "$SHA256_1" = "$SHA256_2" ]; then
    echo ""
    echo "  ✓✓✓ HASHES MATCH - BINARIES ARE IDENTICAL ✓✓✓"
    IDENTICAL=true
else
    echo ""
    echo "  ✗✗✗ HASHES DIFFER - BINARIES ARE DIFFERENT ✗✗✗"
    IDENTICAL=false
fi

echo ""

# ============================================================================
# Build Information
# ============================================================================
if [ "$IDENTICAL" = false ]; then
    echo "Build Information Comparison:"
    echo "----------------------------------------"

    if command -v go &> /dev/null; then
        echo ""
        echo "Provider 1:"
        go version -m "$PROVIDER1" 2>&1 | head -10 | sed 's/^/  /'

        echo ""
        echo "Provider 2:"
        go version -m "$PROVIDER2" 2>&1 | head -10 | sed 's/^/  /'
    else
        echo "  ⚠️  Go not installed, cannot compare build info"
    fi

    echo ""
fi

# ============================================================================
# Summary
# ============================================================================
echo "========================================"
echo "  Summary"
echo "========================================"
echo ""

if [ "$IDENTICAL" = true ]; then
    echo "✓ The providers are IDENTICAL"
    echo "  SHA256: $SHA256_1"
    echo ""
    echo "You can safely use either binary - they are exactly the same."
else
    echo "✗ The providers are DIFFERENT"
    echo ""
    echo "Possible reasons:"
    echo "  1. Built from different commits/branches"
    echo "  2. Built with different Go versions"
    echo "  3. Built with different build flags"
    echo "  4. Built at different times (if using build timestamps)"
    echo ""
    echo "Recommendation:"
    echo "  - Use the same binary on all machines"
    echo "  - Rebuild from the same commit with same Go version"
    echo "  - Distribute the binary instead of rebuilding"
fi

echo ""
