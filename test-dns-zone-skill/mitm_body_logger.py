"""mitmproxy addon that logs request/response bodies to a file."""
import json, os

LOG_FILE = os.path.join(os.path.dirname(__file__), "mitm_bodies.log")

def request(flow):
    if "api.gcore.com" not in (flow.request.host or ""):
        return
    body = flow.request.get_text() or ""
    with open(LOG_FILE, "a") as f:
        f.write(f"\n{'='*80}\n")
        f.write(f">>> {flow.request.method} {flow.request.url}\n")
        if body:
            try:
                parsed = json.loads(body)
                f.write(f"BODY: {json.dumps(parsed, indent=2)}\n")
            except:
                f.write(f"BODY: {body[:2000]}\n")

def response(flow):
    if "api.gcore.com" not in (flow.request.host or ""):
        return
    body = flow.response.get_text() or ""
    with open(LOG_FILE, "a") as f:
        f.write(f"<<< {flow.response.status_code} {flow.request.method} {flow.request.url}\n")
        if body:
            try:
                parsed = json.loads(body)
                f.write(f"RESPONSE: {json.dumps(parsed, indent=2)}\n")
            except:
                f.write(f"RESPONSE: {body[:2000]}\n")
