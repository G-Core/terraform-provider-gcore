#!/bin/bash
#
# Environment setup for manual testing with mitmproxy
# Usage: source ./set_env.sh
#

# Check if script is being sourced (not executed)
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    echo "❌ ERROR: This script must be sourced, not executed!"
    echo ""
    echo "Run this instead:"
    echo "  source ./set_env.sh"
    echo ""
    echo "Or:"
    echo "  . ./set_env.sh"
    echo ""
    exit 1
fi

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "========================================"
echo "  Setting up environment for mitmproxy"
echo "========================================"
echo ""

# Load Gcore credentials from .env
if [ -f "$PROJECT_ROOT/.env" ]; then
    set -o allexport
    source "$PROJECT_ROOT/.env"
    set +o allexport
    echo "✓ Loaded credentials from .env"
else
    echo "⚠️  WARNING: .env file not found"
    echo "  Create .env with your Gcore credentials:"
    echo "    GCORE_CLOUD_API_KEY=your_key"
    echo "    GCORE_CLOUD_PROJECT_ID=379987"
    echo "    GCORE_CLOUD_REGION_ID=76"
fi

# Terraform configuration
export TF_CLI_CONFIG_FILE="$PROJECT_ROOT/.terraformrc"
echo "✓ TF_CLI_CONFIG_FILE=$TF_CLI_CONFIG_FILE"

# Terraform logging
export TF_LOG=TRACE
export TF_LOG_PATH="$PROJECT_ROOT/terraform_manual.log"
echo "✓ TF_LOG=TRACE (output: terraform_manual.log)"

# Allow plugin cache to break dependency lock file
export TF_PLUGIN_CACHE_MAY_BREAK_DEPENDENCY_LOCK_FILE="1"

# Proxy configuration
export HTTP_PROXY="http://127.0.0.1:9092"
export HTTPS_PROXY="http://127.0.0.1:9092"
echo "✓ Proxy: http://127.0.0.1:9092"

# Bypass proxy for terraform registry (avoid SSL issues during init)
export NO_PROXY="registry.terraform.io,releases.hashicorp.com"
echo "✓ NO_PROXY for Terraform registry"

# SSL Certificate configuration
MITM_CERT="$HOME/.mitmproxy/mitmproxy-ca-cert.pem"
CA_BUNDLE="$PROJECT_ROOT/ca-bundle.pem"

if [ ! -f "$MITM_CERT" ]; then
    echo ""
    echo "⚠️  WARNING: mitmproxy CA certificate not found"
    echo "  Start mitmproxy first: ./run_mitm.sh"
    echo "  Then make a test request to generate the certificate"
else
    echo "✓ mitmproxy CA certificate exists"

    # Create certificate bundle if needed
    if [ ! -f "$CA_BUNDLE" ] || [ "$MITM_CERT" -nt "$CA_BUNDLE" ]; then
        echo "  Creating certificate bundle..."

        # Export system certificates
        security find-certificate -a -p /System/Library/Keychains/SystemRootCertificates.keychain > "$CA_BUNDLE" 2>/dev/null || true
        security find-certificate -a -p /Library/Keychains/System.keychain >> "$CA_BUNDLE" 2>/dev/null || true

        # Add mitmproxy certificate
        cat "$MITM_CERT" >> "$CA_BUNDLE"

        CERT_COUNT=$(grep -c 'BEGIN CERTIFICATE' "$CA_BUNDLE" || echo "0")
        echo "  ✓ Certificate bundle created with $CERT_COUNT certificates"
    fi

    export SSL_CERT_FILE="$CA_BUNDLE"
    export REQUESTS_CA_BUNDLE="$CA_BUNDLE"
    echo "✓ SSL_CERT_FILE=$CA_BUNDLE"
fi

echo ""
echo "========================================"
echo "  Environment ready!"
echo "========================================"
echo ""
echo "Verify environment variables are set:"
echo "  echo \$HTTP_PROXY           # Should show: http://127.0.0.1:9092"
echo "  echo \$TF_CLI_CONFIG_FILE    # Should show: $TF_CLI_CONFIG_FILE"
echo ""
echo "You can now run terraform commands:"
echo "  cd test-router-manual"
echo "  terraform init"
echo "  terraform plan"
echo "  terraform apply"
echo ""
echo "HTTP traffic will be captured to: flow.mitm"
echo "View with: mitmproxy -r flow.mitm"
echo ""
echo "To reset environment, close terminal or run:"
echo "  unset HTTP_PROXY HTTPS_PROXY SSL_CERT_FILE REQUESTS_CA_BUNDLE"
echo ""
