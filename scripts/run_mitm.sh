#!/bin/bash

# Terminate any existing mitmproxy on port 9092
lsof -ti :9092 | xargs kill -9 2>/dev/null || true

# Remove old capture file
rm -f flow.mitm

echo "=== Starting mitmproxy on port 9092 ==="
echo "Capture will be saved to: flow.mitm"
echo "Request bodies will be logged to: mitm_requests.log"
echo ""
echo "Press Ctrl+C to stop capturing"
echo ""

# Create addon script to log request/response bodies
cat > /tmp/mitm_body_logger.py << 'EOF'
import sys
import json
from mitmproxy import http

class BodyLogger:
    def __init__(self):
        self.log_file = open("mitm_requests.log", "a", buffering=1)

    def request(self, flow: http.HTTPFlow) -> None:
        """Log request details including body"""
        req = flow.request

        # Log request header
        self.log_file.write(f"\n{'='*80}\n")
        self.log_file.write(f"REQUEST: {req.method} {req.pretty_url}\n")
        self.log_file.write(f"{'='*80}\n")

        # Log request headers
        self.log_file.write("Headers:\n")
        for key, value in req.headers.items():
            if key.lower() not in ['authorization', 'cookie']:  # Hide sensitive headers
                self.log_file.write(f"  {key}: {value}\n")

        # Log request body
        if req.content:
            self.log_file.write(f"\nRequest Body ({len(req.content)} bytes):\n")
            try:
                # Try to parse as JSON for pretty printing
                body_json = json.loads(req.content.decode('utf-8', errors='replace'))
                self.log_file.write(json.dumps(body_json, indent=2))
            except:
                # If not JSON, write raw content
                self.log_file.write(req.content.decode('utf-8', errors='replace'))
            self.log_file.write("\n")
        else:
            self.log_file.write("\nNo request body\n")

    def response(self, flow: http.HTTPFlow) -> None:
        """Log response details including body"""
        resp = flow.response

        self.log_file.write(f"\n{'-'*80}\n")
        self.log_file.write(f"RESPONSE: {resp.status_code} {resp.reason}\n")
        self.log_file.write(f"{'-'*80}\n")

        # Log response body for API calls
        if resp.content and len(resp.content) < 50000:  # Limit size
            self.log_file.write(f"\nResponse Body ({len(resp.content)} bytes):\n")
            try:
                body_json = json.loads(resp.content.decode('utf-8', errors='replace'))
                self.log_file.write(json.dumps(body_json, indent=2))
            except:
                self.log_file.write(resp.content.decode('utf-8', errors='replace')[:1000])
            self.log_file.write("\n")

        self.log_file.write("\n")

addons = [BodyLogger()]
EOF

# Start mitmproxy with body logging addon
mitmdump -p 9092 \
  -w flow.mitm \
  --set flow_detail=2 \
  --set console_eventlog_verbosity=error \
  -s /tmp/mitm_body_logger.py

echo ""
echo "=== Mitmproxy stopped ==="
echo "Capture saved to: flow.mitm"
echo "Request bodies saved to: mitm_requests.log"
echo ""
echo "To view capture:"
echo "  mitmproxy -r flow.mitm"
echo ""
echo "To view request bodies:"
echo "  cat mitm_requests.log"
