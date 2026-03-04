#!/bin/bash
#
# Analyze mitmproxy capture for router interface operations
# Shows detailed request/response information
#

if [ ! -f "http_requests.log" ] && [ ! -f "router_flows.mitm" ]; then
    echo "ERROR: No mitmproxy logs found"
    echo "Run test_router_mitmproxy.sh first"
    exit 1
fi

# Prefer http_requests.log (human-readable)
USE_TEXT_LOG=false
if [ -f "http_requests.log" ]; then
    USE_TEXT_LOG=true
    echo "Using http_requests.log (human-readable format)"
else
    echo "Using router_flows.mitm (binary format)"
fi

echo "======================================================================"
echo "  mitmproxy Capture Analysis - Detailed View"
echo "======================================================================"
echo ""

# ============================================================================
# Show all HTTP requests to routers endpoint
# ============================================================================
echo "----------------------------------------------------------------------"
echo "ALL Router API Requests"
echo "----------------------------------------------------------------------"
echo ""

if [ "$USE_TEXT_LOG" = true ]; then
    grep -E "REQUEST:.*routers" http_requests.log | head -50
else
    mitmdump -nr router_flows.mitm "~u /routers/" 2>/dev/null | \
        grep -E "POST|PATCH|GET|DELETE|PUT" | \
        head -50
fi

echo ""

# ============================================================================
# Analyze POST /attach requests
# ============================================================================
echo "----------------------------------------------------------------------"
echo "POST /attach Requests (Adding Interface)"
echo "----------------------------------------------------------------------"
echo ""

if [ "$USE_TEXT_LOG" = true ]; then
    ATTACH_REQUESTS=$(grep -A 20 "REQUEST: POST.*attach" http_requests.log 2>/dev/null)
    ATTACH_COUNT=$(grep -c "REQUEST: POST.*attach" http_requests.log 2>/dev/null || echo "0")
else
    ATTACH_REQUESTS=$(mitmdump -nr router_flows.mitm "~m POST & ~u /attach" 2>/dev/null)
    ATTACH_COUNT=$(echo "$ATTACH_REQUESTS" | grep -c "POST" 2>/dev/null || echo "0")
fi

if [ -z "$ATTACH_REQUESTS" ]; then
    echo "No POST /attach requests found"
else
    echo "$ATTACH_REQUESTS" | head -50
    echo ""
    echo "Count: $ATTACH_COUNT"
fi

echo ""

# ============================================================================
# Analyze POST /detach requests
# ============================================================================
echo "----------------------------------------------------------------------"
echo "POST /detach Requests (Removing Interface)"
echo "----------------------------------------------------------------------"
echo ""

if [ "$USE_TEXT_LOG" = true ]; then
    DETACH_REQUESTS=$(grep -A 20 "REQUEST: POST.*detach" http_requests.log 2>/dev/null)
    DETACH_COUNT=$(grep -c "REQUEST: POST.*detach" http_requests.log 2>/dev/null || echo "0")
else
    DETACH_REQUESTS=$(mitmdump -nr router_flows.mitm "~m POST & ~u /detach" 2>/dev/null)
    DETACH_COUNT=$(echo "$DETACH_REQUESTS" | grep -c "POST" 2>/dev/null || echo "0")
fi

if [ -z "$DETACH_REQUESTS" ]; then
    echo "No POST /detach requests found"
else
    echo "$DETACH_REQUESTS" | head -50
    echo ""
    echo "Count: $DETACH_COUNT"
fi

echo ""

# ============================================================================
# Analyze PATCH requests
# ============================================================================
echo "----------------------------------------------------------------------"
echo "PATCH /routers Requests (Should be minimal for interface-only changes)"
echo "----------------------------------------------------------------------"
echo ""

if [ "$USE_TEXT_LOG" = true ]; then
    PATCH_REQUESTS=$(grep -A 20 "REQUEST: PATCH.*routers" http_requests.log 2>/dev/null)
    PATCH_COUNT=$(grep -c "REQUEST: PATCH.*routers" http_requests.log 2>/dev/null || echo "0")

    if [ -z "$PATCH_REQUESTS" ]; then
        echo "✅ No PATCH requests found - GOOD!"
        echo "This means interface changes are handled via attach/detach only"
    else
        echo "PATCH requests found:"
        echo "$PATCH_REQUESTS" | head -50
        echo ""
        echo "Count: $PATCH_COUNT"
        echo ""
        echo "Checking if PATCH contains 'interfaces' field in body..."

        if grep -i "interfaces" http_requests.log | grep -A 5 "Request Body" | grep -q "interfaces"; then
            echo "❌ FAIL: Found 'interfaces' in PATCH payload"
        else
            echo "✅ PASS: No 'interfaces' field in PATCH payload"
        fi
    fi
else
    PATCH_REQUESTS=$(mitmdump -nr router_flows.mitm "~m PATCH & ~u /routers/" 2>/dev/null)
    PATCH_COUNT=$(echo "$PATCH_REQUESTS" | grep -c "PATCH" 2>/dev/null || echo "0")

    if [ -z "$PATCH_REQUESTS" ]; then
        echo "✅ No PATCH requests found - GOOD!"
        echo "This means interface changes are handled via attach/detach only"
    else
        echo "PATCH requests found:"
        echo "$PATCH_REQUESTS"
        echo ""
        echo "Count: $PATCH_COUNT"
        echo ""
        echo "Checking if PATCH contains 'interfaces' field in body..."

        mitmdump -nr router_flows.mitm "~m PATCH & ~u /routers/" -w - 2>/dev/null | \
            strings | grep -i "interfaces" && \
            echo "❌ FAIL: Found 'interfaces' in PATCH payload" || \
            echo "✅ PASS: No 'interfaces' field in PATCH payload"
    fi
fi

echo ""

# ============================================================================
# Summary statistics
# ============================================================================
echo "======================================================================"
echo "  Summary Statistics"
echo "======================================================================"
echo ""

TOTAL_REQUESTS=$(mitmdump -nr router_flows.mitm 2>/dev/null | grep -c "^<" || echo "0")
ROUTER_REQUESTS=$(mitmdump -nr router_flows.mitm "~u /routers/" 2>/dev/null | grep -c "^\s*POST\|^\s*PATCH\|^\s*GET" || echo "0")
POST_ATTACH=$(mitmdump -nr router_flows.mitm "~m POST & ~u /attach" 2>/dev/null | grep -c "POST" || echo "0")
POST_DETACH=$(mitmdump -nr router_flows.mitm "~m POST & ~u /detach" 2>/dev/null | grep -c "POST" || echo "0")
PATCH_ROUTER=$(mitmdump -nr router_flows.mitm "~m PATCH & ~u /routers/" 2>/dev/null | grep -c "PATCH" || echo "0")

echo "Total HTTP requests captured:     $TOTAL_REQUESTS"
echo "Router-related requests:          $ROUTER_REQUESTS"
echo ""
echo "  POST /attach requests:          $POST_ATTACH"
echo "  POST /detach requests:          $POST_DETACH"
echo "  PATCH /routers requests:        $PATCH_ROUTER"
echo ""

# ============================================================================
# Verification against Kirill's issues
# ============================================================================
echo "======================================================================"
echo "  Issue Verification"
echo "======================================================================"
echo ""

echo "Issue #1: Adding interface sends both POST attach AND PATCH"
echo "  Expected: Only POST attach (count >= 1), NO PATCH"
echo "  Actual:   POST attach = $POST_ATTACH, PATCH = $PATCH_ROUTER"

if [ "$POST_ATTACH" -ge 1 ] && [ "$PATCH_ROUTER" -eq 0 ]; then
    echo "  ✅ PASS - Only POST attach used"
elif [ "$POST_ATTACH" -ge 1 ] && [ "$PATCH_ROUTER" -gt 0 ]; then
    echo "  ⚠️  WARNING - Both POST attach and PATCH detected"
    echo "  Need to verify PATCH is for other fields, not interfaces"
else
    echo "  ❌ FAIL - No POST attach found"
fi
echo ""

echo "Issue #2: Removing interface sends PATCH with empty payload"
echo "  Expected: POST detach (count >= 1)"
echo "  Actual:   POST detach = $POST_DETACH"

if [ "$POST_DETACH" -ge 1 ]; then
    echo "  ✅ PASS - POST detach used correctly"
else
    echo "  ❌ FAIL - No POST detach found"
fi
echo ""

# ============================================================================
# Show sample request details
# ============================================================================
echo "======================================================================"
echo "  Sample Request Details"
echo "======================================================================"
echo ""

echo "To view detailed request/response bodies, run:"
echo "  mitmproxy -r router_flows.mitm"
echo ""
echo "Or export to readable text:"
echo "  mitmdump -nr router_flows.mitm > flows.txt"
echo ""

echo "Useful filters:"
echo "  Show only attach:  mitmdump -nr router_flows.mitm \"~u /attach\""
echo "  Show only detach:  mitmdump -nr router_flows.mitm \"~u /detach\""
echo "  Show only PATCH:   mitmdump -nr router_flows.mitm \"~m PATCH\""
echo "  Show router ops:   mitmdump -nr router_flows.mitm \"~u /routers/\""
echo ""
