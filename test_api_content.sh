#!/bin/bash
source /Users/user/repos/gcore-terraform/.env

echo "=== Test 1: Content without quotes ==="
curl -s -X POST \
  -H "Authorization: APIKey $GCORE_API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/dns/v2/zones/maxima.lt/tf-api-test1.maxima.lt/A" \
  -d '{"ttl": 300, "resource_records": [{"content": ["192.168.1.1"], "enabled": true}]}'
echo ""

echo "=== Test 2: Content with quotes ==="
curl -s -X POST \
  -H "Authorization: APIKey $GCORE_API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/dns/v2/zones/maxima.lt/tf-api-test2.maxima.lt/A" \
  -d '{"ttl": 300, "resource_records": [{"content": ["\"192.168.1.1\""], "enabled": true}]}'
echo ""

echo "=== Test 3: Get existing record to see format ==="
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/dns/v2/zones/maxima.lt/tf-new-provider-test.maxima.lt/A"
echo ""

echo "=== Cleanup ==="
curl -s -X DELETE -H "Authorization: APIKey $GCORE_API_KEY" "https://api.gcore.com/dns/v2/zones/maxima.lt/tf-api-test1.maxima.lt/A"
curl -s -X DELETE -H "Authorization: APIKey $GCORE_API_KEY" "https://api.gcore.com/dns/v2/zones/maxima.lt/tf-api-test2.maxima.lt/A"
echo "Done"
