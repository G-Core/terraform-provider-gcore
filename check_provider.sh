#!/bin/bash
#
# Check provider version and hash
# Usage: ./check_provider.sh [path/to/provider]
#

set -e

PROVIDER_PATH="${1:-./terraform-provider-gcore}"

echo "========================================"
echo "  Provider Binary Information"
echo "========================================"
echo ""

# Check if provider exists
if [ ! -f "$PROVIDER_PATH" ]; then
    echo "✗ ERROR: Provider not found at: $PROVIDER_PATH"
    echo ""
    echo "Usage: $0 [path/to/provider]"
    echo "Example: $0 ./terraform-provider-gcore"
    exit 1
fi

echo "Provider path: $PROVIDER_PATH"
echo ""

# ============================================================================
# File Information
# ============================================================================
echo "File Information:"
echo "----------------------------------------"

# File size
SIZE=$(ls -lh "$PROVIDER_PATH" | awk '{print $5}')
echo "  Size: $SIZE"

# File permissions
PERMS=$(ls -l "$PROVIDER_PATH" | awk '{print $1}')
echo "  Permissions: $PERMS"

# Modification date
MOD_DATE=$(ls -l "$PROVIDER_PATH" | awk '{print $6, $7, $8}')
echo "  Modified: $MOD_DATE"

echo ""

# ============================================================================
# Hashes
# ============================================================================
echo "Cryptographic Hashes:"
echo "----------------------------------------"

# SHA256
SHA256=$(shasum -a 256 "$PROVIDER_PATH" | awk '{print $1}')
echo "  SHA256: $SHA256"

# MD5 (for quick comparison)
MD5=$(md5 -q "$PROVIDER_PATH")
echo "  MD5:    $MD5"

echo ""

# ============================================================================
# Go Build Information
# ============================================================================
echo "Go Build Information:"
echo "----------------------------------------"

if command -v go &> /dev/null; then
    # Get build info
    BUILD_INFO=$(go version -m "$PROVIDER_PATH" 2>&1)

    if echo "$BUILD_INFO" | grep -q "go version"; then
        # Extract Go version
        GO_VERSION=$(echo "$BUILD_INFO" | grep "go version" | head -1)
        echo "  $GO_VERSION"

        # Extract build settings
        if echo "$BUILD_INFO" | grep -q "build"; then
            echo ""
            echo "  Build Settings:"
            echo "$BUILD_INFO" | grep "build" | sed 's/^/    /'
        fi

        # Extract dependencies (first few)
        if echo "$BUILD_INFO" | grep -q "dep"; then
            echo ""
            echo "  Key Dependencies:"
            echo "$BUILD_INFO" | grep -E "dep.*(github.com/G-Core/gcore-go|github.com/hashicorp/terraform-plugin-framework)" | head -5 | sed 's/^/    /'
        fi

        # Extract VCS info if available
        if echo "$BUILD_INFO" | grep -q "vcs"; then
            echo ""
            echo "  VCS Information:"
            echo "$BUILD_INFO" | grep "vcs" | sed 's/^/    /'
        fi
    else
        echo "  Unable to read build info (not a Go binary or stripped)"
    fi
else
    echo "  ⚠️  Go not installed, skipping build info"
fi

echo ""

# ============================================================================
# Provider Version (from Terraform)
# ============================================================================
echo "Provider Version (Terraform):"
echo "----------------------------------------"

# Create a temporary test directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

cd "$TEMP_DIR"

# Create minimal terraform config
cat > test.tf << 'EOF'
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}
EOF

# Create temporary .terraformrc pointing to our provider
TEMP_RC=$(mktemp)
trap "rm -f $TEMP_RC" EXIT

cat > "$TEMP_RC" << EORC
provider_installation {
  dev_overrides {
    "gcore/gcore" = "$(dirname "$(cd "$(dirname "$PROVIDER_PATH")" && pwd)/$(basename "$PROVIDER_PATH")")"
  }
  direct {}
}
EORC

# Try to get version from Terraform
export TF_CLI_CONFIG_FILE="$TEMP_RC"

if command -v terraform &> /dev/null; then
    # Initialize (will show provider info)
    INIT_OUTPUT=$(terraform init 2>&1 || true)

    # Try to extract version from any output
    if echo "$INIT_OUTPUT" | grep -q "gcore/gcore"; then
        echo "$INIT_OUTPUT" | grep -i "gcore" | head -3 | sed 's/^/  /'
    else
        echo "  Using dev_overrides (version from binary)"
    fi

    # Try terraform version
    TF_VERSION=$(terraform version 2>/dev/null | head -1 || echo "  Terraform: not available")
    echo ""
    echo "  $TF_VERSION"
else
    echo "  ⚠️  Terraform not installed, cannot check version"
fi

cd - > /dev/null 2>&1

echo ""

# ============================================================================
# Summary for Distribution
# ============================================================================
echo "========================================"
echo "  Distribution Checksum"
echo "========================================"
echo ""
echo "To verify identical binary on another machine:"
echo ""
echo "  1. Copy this checksum:"
echo "     $SHA256"
echo ""
echo "  2. On the other machine, run:"
echo "     shasum -a 256 terraform-provider-gcore"
echo ""
echo "  3. Compare the output - they should match exactly"
echo ""
echo "Quick verification command:"
echo "  echo \"$SHA256  terraform-provider-gcore\" | shasum -a 256 -c"
echo ""

# Create checksum file
CHECKSUM_FILE="$(dirname "$PROVIDER_PATH")/terraform-provider-gcore.sha256"
echo "$SHA256  $(basename "$PROVIDER_PATH")" > "$CHECKSUM_FILE"
echo "✓ Checksum saved to: $CHECKSUM_FILE"
echo ""
