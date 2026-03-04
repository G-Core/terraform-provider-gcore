#!/bin/bash
#
# Simple mitmproxy runner for manual testing
# Usage: ./run_mitm.sh
# Press Ctrl+C to stop
#

set -e

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
FLOW_FILE="$PROJECT_ROOT/flow.mitm"
LOG_FILE="$PROJECT_ROOT/mitm_output.log"
PORT=9092

# Cleanup function
cleanup() {
    echo ""
    echo "Stopping mitmproxy..."

    if [ ! -z "$MITM_PID" ]; then
        kill $MITM_PID 2>/dev/null || true
        sleep 1
        # Force kill if still running
        if ps -p $MITM_PID > /dev/null 2>&1; then
            kill -9 $MITM_PID 2>/dev/null || true
        fi
    fi

    # Kill any remaining processes on port
    lsof -ti :$PORT 2>/dev/null | xargs kill -9 2>/dev/null || true

    echo "✓ mitmproxy stopped"

    if [ -f "$FLOW_FILE" ] && [ -s "$FLOW_FILE" ]; then
        SIZE=$(ls -lh "$FLOW_FILE" | awk '{print $5}')
        echo ""
        echo "Captured traffic saved to: $FLOW_FILE ($SIZE)"
        echo "View with: mitmproxy -r flow.mitm"
    fi

    exit 0
}

# Register cleanup on Ctrl+C and exit
trap cleanup INT TERM EXIT

echo "========================================"
echo "  mitmproxy Manual Testing"
echo "========================================"
echo ""

# Check if mitmproxy is installed
if ! command -v mitmdump &> /dev/null; then
    echo "✗ ERROR: mitmproxy not installed"
    echo "  Install with: brew install mitmproxy"
    exit 1
fi

# Check if port is already in use
if lsof -ti :$PORT > /dev/null 2>&1; then
    echo "⚠️  Port $PORT is already in use"
    read -p "Kill existing process? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        lsof -ti :$PORT | xargs kill -9 2>/dev/null || true
        sleep 1
        echo "✓ Port freed"
    else
        echo "✗ Cannot start mitmproxy"
        exit 1
    fi
fi

# Clean up old capture file
if [ -f "$FLOW_FILE" ]; then
    echo "Cleaning previous capture..."
    rm -f "$FLOW_FILE"
    echo "✓ Removed old flow.mitm"
fi

# Clean up old log
rm -f "$LOG_FILE"

echo ""
echo "Starting mitmproxy..."
echo "  Port: $PORT"
echo "  Output: $FLOW_FILE"
echo "  Log: $LOG_FILE"
echo ""

# Start mitmproxy in background
mitmdump -p $PORT \
    -w "$FLOW_FILE" \
    --set stream_large_bodies=1 \
    --ssl-insecure \
    > "$LOG_FILE" 2>&1 &

MITM_PID=$!

# Wait for mitmproxy to start
sleep 2

# Check if mitmproxy is running
if ! ps -p $MITM_PID > /dev/null 2>&1; then
    echo "✗ ERROR: mitmproxy failed to start"
    echo "Check log: tail $LOG_FILE"
    tail -20 "$LOG_FILE"
    exit 1
fi

echo "✓ mitmproxy started (PID: $MITM_PID)"
echo ""
echo "========================================"
echo "  Ready for manual testing!"
echo "========================================"
echo ""
echo "1. In another terminal, run: source ./set_env.sh"
echo "2. Run your terraform commands"
echo "3. Press Ctrl+C here to stop mitmproxy"
echo ""
echo "Proxy: http://127.0.0.1:$PORT"
echo ""
echo "Waiting... (Ctrl+C to stop)"
echo ""

# Wait forever (until Ctrl+C)
while true; do
    sleep 1

    # Check if mitmproxy is still running
    if ! ps -p $MITM_PID > /dev/null 2>&1; then
        echo ""
        echo "⚠️  mitmproxy process died unexpectedly"
        echo "Check log: tail $LOG_FILE"
        tail -20 "$LOG_FILE"
        exit 1
    fi
done
