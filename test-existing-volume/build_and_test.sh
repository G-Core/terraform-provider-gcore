#!/bin/bash
set -e

# Navigate to provider root
cd /Users/user/repos/gcore-terraform

# Build provider
echo "Building provider..."
go build -o terraform-provider-gcore
echo "Build complete"

# Check binary
ls -la terraform-provider-gcore
