#!/usr/bin/env python3
"""
mitmproxy addon to log HTTP requests and responses
"""

import logging
from mitmproxy import http

class HTTPLogger:
    def __init__(self):
        self.log_file = "http_requests.log"
        # Clear the log file
        with open(self.log_file, "w") as f:
            f.write("=" * 80 + "\n")
            f.write("HTTP Traffic Log\n")
            f.write("=" * 80 + "\n\n")

    def request(self, flow: http.HTTPFlow) -> None:
        """Log HTTP request"""
        req = flow.request

        # Only log api.gcore.com requests
        if "api.gcore.com" not in req.host and "gcorelabs.com" not in req.host:
            return

        with open(self.log_file, "a") as f:
            f.write("\n" + "=" * 80 + "\n")
            f.write(f"REQUEST: {req.method} {req.scheme}://{req.host}{req.path}\n")
            f.write("-" * 80 + "\n")
            f.write(f"Method: {req.method}\n")
            f.write(f"URL: {req.url}\n")
            f.write(f"HTTP Version: {req.http_version}\n")

            # Headers
            f.write("\nRequest Headers:\n")
            for k, v in req.headers.items():
                # Don't log sensitive auth headers in full
                if k.lower() in ["authorization", "cookie"]:
                    f.write(f"  {k}: [REDACTED]\n")
                else:
                    f.write(f"  {k}: {v}\n")

            # Body (if present and not too large)
            if req.content:
                try:
                    body = req.text
                    if len(body) < 5000:  # Only log bodies under 5KB
                        f.write(f"\nRequest Body ({len(body)} bytes):\n")
                        f.write(body[:5000] + "\n")
                    else:
                        f.write(f"\nRequest Body: [Large body: {len(body)} bytes]\n")
                except:
                    f.write(f"\nRequest Body: [Binary data: {len(req.content)} bytes]\n")

    def response(self, flow: http.HTTPFlow) -> None:
        """Log HTTP response"""
        req = flow.request
        resp = flow.response

        # Only log api.gcore.com requests
        if "api.gcore.com" not in req.host and "gcorelabs.com" not in req.host:
            return

        if resp is None:
            return

        with open(self.log_file, "a") as f:
            f.write("\nRESPONSE:\n")
            f.write("-" * 80 + "\n")
            f.write(f"Status Code: {resp.status_code} {resp.reason}\n")
            f.write(f"HTTP Version: {resp.http_version}\n")

            # Headers
            f.write("\nResponse Headers:\n")
            for k, v in resp.headers.items():
                f.write(f"  {k}: {v}\n")

            # Body (if present and not too large)
            if resp.content:
                try:
                    body = resp.text
                    if len(body) < 5000:  # Only log bodies under 5KB
                        f.write(f"\nResponse Body ({len(body)} bytes):\n")
                        f.write(body[:5000] + "\n")
                    else:
                        f.write(f"\nResponse Body: [Large body: {len(body)} bytes]\n")
                        # For large bodies, show first 500 chars
                        f.write(body[:500] + "...\n")
                except:
                    f.write(f"\nResponse Body: [Binary data: {len(resp.content)} bytes]\n")

            f.write("\n")

addons = [HTTPLogger()]
