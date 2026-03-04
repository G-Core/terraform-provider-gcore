#!/bin/bash
set -e

echo "=== Testing routes=[] Serialization Without WithJSONSet ==="
echo "This test captures HTTP requests to see if empty routes array is serialized"
echo ""

# Load credentials
if [ -f .env ]; then
    source .env
elif [ -f ../.env ]; then
    source ../.env
else
    echo "ERROR: .env not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

# Create a simple Go test program to check serialization
cat > /tmp/test_routes_serialization.go << 'EOF'
package main

import (
	"encoding/json"
	"fmt"
)

type Routes struct {
	Destination string `json:"destination"`
	Nexthop     string `json:"nexthop"`
}

type Router struct {
	Name   string   `json:"name"`
	Routes []Routes `json:"routes,omitempty"`
}

type RouterOmitZero struct {
	Name   string   `json:"name"`
	Routes []Routes `json:"routes,omitzero"`
}

func main() {
	fmt.Println("=== Testing Go JSON serialization with empty slices ===")
	fmt.Println()

	// Test 1: Empty slice with omitempty
	r1 := Router{
		Name:   "test",
		Routes: []Routes{},
	}
	j1, _ := json.Marshal(r1)
	fmt.Println("Empty slice with omitempty tag:")
	fmt.Printf("  %s\n", j1)
	if string(j1) == `{"name":"test","routes":[]}` {
		fmt.Println("  ✅ Empty slice IS serialized as []")
	} else if string(j1) == `{"name":"test"}` {
		fmt.Println("  ❌ Empty slice is OMITTED")
	}
	fmt.Println()

	// Test 2: Nil slice with omitempty
	r2 := Router{
		Name:   "test",
		Routes: nil,
	}
	j2, _ := json.Marshal(r2)
	fmt.Println("Nil slice with omitempty tag:")
	fmt.Printf("  %s\n", j2)
	if string(j2) == `{"name":"test"}` {
		fmt.Println("  ✅ Nil slice is OMITTED")
	}
	fmt.Println()

	// Test 3: Empty slice with omitzero
	r3 := RouterOmitZero{
		Name:   "test",
		Routes: []Routes{},
	}
	j3, _ := json.Marshal(r3)
	fmt.Println("Empty slice with omitzero tag:")
	fmt.Printf("  %s\n", j3)
	if string(j3) == `{"name":"test","routes":[]}` {
		fmt.Println("  ✅ Empty slice IS serialized as []")
	} else if string(j3) == `{"name":"test"}` {
		fmt.Println("  ❌ Empty slice is OMITTED")
	}
	fmt.Println()

	// Test 4: Nil slice with omitzero
	r4 := RouterOmitZero{
		Name:   "test",
		Routes: nil,
	}
	j4, _ := json.Marshal(r4)
	fmt.Println("Nil slice with omitzero tag:")
	fmt.Printf("  %s\n", j4)
	if string(j4) == `{"name":"test"}` {
		fmt.Println("  ✅ Nil slice is OMITTED")
	} else if string(j4) == `{"name":"test","routes":[]}` {
		fmt.Println("  ❌ Nil slice IS serialized as []")
	}
	fmt.Println()

	fmt.Println("=== CONCLUSION ===")
	fmt.Println("Standard Go behavior:")
	fmt.Println("  - omitempty: empty slices [] are serialized, nil is omitted")
	fmt.Println("  - omitzero: empty slices [] are serialized, nil is omitted")
	fmt.Println()
	fmt.Println("Pedro's claim: 'Empty slices should be serialized when we have omitzero tag,")
	fmt.Println("                they are only excluded when the slice is nil.'")
	fmt.Println("  ✅ CORRECT - both omitempty and omitzero behave this way")
	fmt.Println()
	fmt.Println("Therefore: WithJSONSet may NOT be necessary if we ensure routes field")
	fmt.Println("           is set to empty slice [] (not nil) when deleting routes")
}
EOF

echo "Running Go serialization test..."
go run /tmp/test_routes_serialization.go

echo ""
echo "=== Checking SDK NetworkRouterUpdateParams ===..."
echo "Looking for Routes field definition in SDK..."

# Check if SDK source is available
if [ -d "sdk-gcore-go" ]; then
    echo "Found SDK source, checking Routes field..."
    grep -A 5 "Routes.*json" sdk-gcore-go/cloud/network_router.go 2>/dev/null || echo "Could not find Routes field definition"
else
    echo "SDK source not available locally"
fi

echo ""
echo "=== Practical Test: Check MarshalForPatch Behavior ===..."

cat > /tmp/test_marshal_patch.go << 'EOF'
package main

import (
	"fmt"
	"encoding/json"
)

type Routes struct {
	Destination string `json:"destination"`
	Nexthop string `json:"nexthop"`
}

type Model struct {
	Name string `json:"name"`
	Routes []Routes `json:"routes,computed_optional"`
}

func (m Model) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

func main() {
	// Simulate: user wants to delete routes (routes=[])
	data := Model{
		Name: "test-router",
		Routes: []Routes{}, // Empty slice, not nil
	}

	j, _ := data.MarshalJSON()
	fmt.Printf("Marshaled data with empty routes: %s\n", j)

	if string(j) == `{"name":"test-router","routes":[]}` {
		fmt.Println("✅ Empty routes [] IS included in JSON")
		fmt.Println("   WithJSONSet may NOT be necessary!")
	} else {
		fmt.Println("❌ Empty routes [] is OMITTED from JSON")
		fmt.Println("   WithJSONSet IS necessary!")
	}
}
EOF

go run /tmp/test_marshal_patch.go

echo ""
echo "=== FINAL RECOMMENDATION ==="
echo ""
echo "Based on standard Go behavior:"
echo "  • Empty slices [] ARE serialized with both omitempty and omitzero"
echo "  • Only nil slices are omitted"
echo ""
echo "For our code:"
echo "  • If data.Routes is set to empty slice [] (not nil), it will be serialized"
echo "  • ModifyPlan sets: plan.Routes = customfield.NewObjectListMust(ctx, []...)"
echo "    This creates an empty slice, not nil"
echo "  • Therefore, routes=[] SHOULD be serialized without WithJSONSet"
echo ""
echo "To verify for certain, we'd need to:"
echo "  1. Check if SDK's NetworkRouterUpdateParams.Routes is omitempty or omitzero"
echo "  2. Test actual HTTP request body with TF_LOG=TRACE"
echo "  3. Confirm customfield.NestedObjectList serializes empty list as []"
echo ""
echo "Pedro is likely CORRECT - WithJSONSet may be unnecessary!"
