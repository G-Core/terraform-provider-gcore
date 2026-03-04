#!/bin/bash
#
# Create certificate bundle with system CAs + mitmproxy CA
#

set -e

echo "Creating certificate bundle for Go..."

# Export system certificates
security find-certificate -a -p /System/Library/Keychains/SystemRootCertificates.keychain > /tmp/system_certs.pem
security find-certificate -a -p /Library/Keychains/System.keychain >> /tmp/system_certs.pem

# Add mitmproxy certificate
cat ~/.mitmproxy/mitmproxy-ca-cert.pem >> /tmp/system_certs.pem

# Save to project directory
cp /tmp/system_certs.pem ./ca-bundle.pem

echo "✓ Certificate bundle created: $(pwd)/ca-bundle.pem"
echo "✓ Contains $(grep -c 'BEGIN CERTIFICATE' ./ca-bundle.pem) certificates"
