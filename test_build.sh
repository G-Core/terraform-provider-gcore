#!/bin/bash
cd /Users/user/repos/gcore-terraform
echo "=== Building provider ==="
go build -o terraform-provider-gcore . 2>&1
BUILD_EXIT=$?
if [ $BUILD_EXIT -eq 0 ]; then
    echo "Build SUCCESS"
    ls -la terraform-provider-gcore
else
    echo "Build FAILED with exit code $BUILD_EXIT"
fi
