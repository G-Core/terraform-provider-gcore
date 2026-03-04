#!/bin/bash
#
# Manual cleanup script for mitmproxy processes
# Use this if test script crashes and leaves mitmproxy running
#

echo "=========================================="
echo "  mitmproxy Cleanup Script"
echo "=========================================="
echo ""

# Check for processes on port 9092
PIDS=$(lsof -ti :9092 2>/dev/null || true)

if [ -z "$PIDS" ]; then
    echo "✓ No mitmproxy processes found on port 9092"
    echo "  Port is free"
else
    echo "Found processes on port 9092:"
    lsof -i :9092 2>/dev/null | head -10
    echo ""
    echo "PIDs: $PIDS"
    echo ""
    read -p "Kill these processes? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Killing processes..."
        echo "$PIDS" | xargs kill -9 2>/dev/null || true
        sleep 1

        # Verify
        REMAINING=$(lsof -ti :9092 2>/dev/null || true)
        if [ -z "$REMAINING" ]; then
            echo "✓ Successfully killed all processes"
            echo "  Port 9092 is now free"
        else
            echo "⚠️  Some processes may still be running:"
            lsof -i :9092 2>/dev/null
        fi
    else
        echo "Cancelled"
    fi
fi

echo ""

# Also check for any mitmdump processes
MITM_PROCESSES=$(pgrep -f "mitmdump.*9092" 2>/dev/null || true)

if [ ! -z "$MITM_PROCESSES" ]; then
    echo "Found mitmdump processes:"
    ps -fp $MITM_PROCESSES 2>/dev/null || true
    echo ""
    read -p "Kill these mitmdump processes? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        pkill -9 -f "mitmdump.*9092" 2>/dev/null || true
        echo "✓ Killed mitmdump processes"
    fi
fi

echo ""
echo "Cleanup complete!"
echo ""
echo "To verify port is free:"
echo "  lsof -i :9092"
echo ""
