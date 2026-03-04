import json
from mitmproxy import http

def request(flow: http.HTTPFlow) -> None:
    if "api.gcore.com" in flow.request.host:
        method = flow.request.method
        url = flow.request.url
        body = flow.request.get_text() or ""
        print(f"\n>>> {method} {url}")
        if body:
            try:
                parsed = json.loads(body)
                print(f"BODY: {json.dumps(parsed, indent=2)}")
            except:
                print(f"BODY: {body}")

def response(flow: http.HTTPFlow) -> None:
    if "api.gcore.com" in flow.request.host:
        status = flow.response.status_code
        body = flow.response.get_text() or ""
        print(f"<<< {status}")
        if body:
            try:
                parsed = json.loads(body)
                print(f"RESPONSE: {json.dumps(parsed, indent=2)}")
            except:
                print(f"RESPONSE: {body[:500]}")
