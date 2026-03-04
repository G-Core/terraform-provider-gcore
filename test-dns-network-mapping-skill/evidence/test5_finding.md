## Test 5 Finding: Name field is immutable

API PUT /dns/v2/network-mappings/{id} returns:
500 Internal Server Error {"error":"field 'name' should not change"}

Action needed: Add RequiresReplace() plan modifier to name field in schema.go
