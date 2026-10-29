package cloud_security_group_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_security_group"
)

const (
	testGroupID = "11111111-1111-4111-8111-111111111111"
	testTaskID  = "22222222-2222-4222-8222-222222222222"
	testRuleID  = "33333333-3333-4333-8333-333333333333"
)

// recordedRequest is one write request that the fake API received.
type recordedRequest struct {
	method string
	path   string
	body   map[string]any
}

// fakeSecurityGroupAPI serves the endpoints that Create and Update call. It
// records every POST and PATCH, and returns groupBody for every group read.
type fakeSecurityGroupAPI struct {
	t         *testing.T
	groupBody map[string]any

	mu     sync.Mutex
	writes []recordedRequest
}

func (f *fakeSecurityGroupAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	switch {
	case r.Method == http.MethodPost || r.Method == http.MethodPatch:
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			f.t.Errorf("%s body is not a JSON object: %s: %q", r.Method, err, raw)
		}
		f.writes = append(f.writes, recordedRequest{method: r.Method, path: r.URL.Path, body: body})
		tasks := []string{}
		if r.Method == http.MethodPost {
			tasks = []string{testTaskID}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"tasks": tasks})
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/cloud/v1/tasks/"):
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":                testTaskID,
			"state":             "FINISHED",
			"created_resources": map[string]any{"security_groups": []string{testGroupID}},
		})
	case r.Method == http.MethodGet:
		_ = json.NewEncoder(w).Encode(f.groupBody)
	default:
		f.t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (f *fakeSecurityGroupAPI) recordedWrites() []recordedRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]recordedRequest(nil), f.writes...)
}

// TestCloudSecurityGroupCreateSendsEmptyRules checks that Create always sends
// `rules: []`. The API adds a default ingress + egress template when `rules` is
// omitted, so the new group must start without rules.
func TestCloudSecurityGroupCreateSendsEmptyRules(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	schema := cloud_security_group.ResourceSchema(ctx)

	api := &fakeSecurityGroupAPI{t: t, groupBody: securityGroupResponse("tf-test-sg", "created", nil)}
	srv := httptest.NewServer(api)
	defer srv.Close()

	r := configuredSecurityGroupResource(ctx, t, srv.URL)

	values := map[string]tftypes.Value{
		"project_id":  tftypes.NewValue(tftypes.Number, 1),
		"region_id":   tftypes.NewValue(tftypes.Number, 2),
		"name":        tftypes.NewValue(tftypes.String, "tf-test-sg"),
		"description": tftypes.NewValue(tftypes.String, "created"),
	}
	req := resource.CreateRequest{
		Config: tfsdk.Config{Raw: securityGroupValue(ctx, t, values), Schema: schema},
		Plan:   tfsdk.Plan{Raw: securityGroupValue(ctx, t, values), Schema: schema},
	}
	resp := &resource.CreateResponse{State: tfsdk.State{Raw: securityGroupNull(ctx), Schema: schema}}
	r.Create(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", resp.Diagnostics.Errors())
	}

	writes := api.recordedWrites()
	if len(writes) != 1 || writes[0].method != http.MethodPost {
		t.Fatalf("got writes %+v, want exactly one POST", writes)
	}
	if want := "/cloud/v2/security_groups/1/2"; writes[0].path != want {
		t.Fatalf("POST path %s, want %s", writes[0].path, want)
	}
	want := map[string]any{"name": "tf-test-sg", "description": "created", "rules": []any{}}
	if !reflect.DeepEqual(writes[0].body, want) {
		t.Fatalf("POST body:\n got %v\nwant %v", writes[0].body, want)
	}
}

// TestCloudSecurityGroupUpdateSendsNoRules checks that a metadata update of a
// group with existing rules sends no rule field. The read-only
// `security_group_rules` output must still reflect the latest API response.
func TestCloudSecurityGroupUpdateSendsNoRules(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	schema := cloud_security_group.ResourceSchema(ctx)

	liveRule := map[string]any{
		"id":                testRuleID,
		"security_group_id": testGroupID,
		"direction":         "ingress",
		"ethertype":         "IPv4",
		"protocol":          "tcp",
		"port_range_min":    22,
		"port_range_max":    22,
		"remote_ip_prefix":  "10.0.0.0/8",
	}
	api := &fakeSecurityGroupAPI{
		t:         t,
		groupBody: securityGroupResponse("tf-test-sg-renamed", "renamed", []any{liveRule}),
	}
	srv := httptest.NewServer(api)
	defer srv.Close()

	r := configuredSecurityGroupResource(ctx, t, srv.URL)

	stateRules := securityGroupRulesValue(ctx, t, map[string]tftypes.Value{
		"id":                tftypes.NewValue(tftypes.String, testRuleID),
		"security_group_id": tftypes.NewValue(tftypes.String, testGroupID),
		"direction":         tftypes.NewValue(tftypes.String, "ingress"),
	})
	common := map[string]tftypes.Value{
		"id":                   tftypes.NewValue(tftypes.String, testGroupID),
		"project_id":           tftypes.NewValue(tftypes.Number, 1),
		"region_id":            tftypes.NewValue(tftypes.Number, 2),
		"security_group_rules": stateRules,
	}
	state := merge(common, map[string]tftypes.Value{
		"name":        tftypes.NewValue(tftypes.String, "tf-test-sg"),
		"description": tftypes.NewValue(tftypes.String, "created"),
	})
	plan := merge(common, map[string]tftypes.Value{
		"name":        tftypes.NewValue(tftypes.String, "tf-test-sg-renamed"),
		"description": tftypes.NewValue(tftypes.String, "renamed"),
	})

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Raw: securityGroupValue(ctx, t, plan), Schema: schema},
		Plan:   tfsdk.Plan{Raw: securityGroupValue(ctx, t, plan), Schema: schema},
		State:  tfsdk.State{Raw: securityGroupValue(ctx, t, state), Schema: schema},
	}
	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: securityGroupValue(ctx, t, state), Schema: schema}}
	r.Update(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics.Errors())
	}

	writes := api.recordedWrites()
	if len(writes) != 1 || writes[0].method != http.MethodPatch {
		t.Fatalf("got writes %+v, want exactly one PATCH", writes)
	}
	if want := "/cloud/v2/security_groups/1/2/" + testGroupID; writes[0].path != want {
		t.Fatalf("PATCH path %s, want %s", writes[0].path, want)
	}
	for _, key := range []string{"rules", "security_group_rules"} {
		if _, found := writes[0].body[key]; found {
			t.Errorf("PATCH body contains %q: %v", key, writes[0].body)
		}
	}
	want := map[string]any{"name": "tf-test-sg-renamed", "description": "renamed"}
	if !reflect.DeepEqual(writes[0].body, want) {
		t.Fatalf("PATCH body:\n got %v\nwant %v", writes[0].body, want)
	}

	var rulesValue customfield.NestedObjectList[cloud_security_group.CloudSecurityGroupSecurityGroupRulesModel]
	if d := resp.State.GetAttribute(ctx, path.Root("security_group_rules"), &rulesValue); d.HasError() {
		t.Fatalf("reading security_group_rules from state: %v", d)
	}
	rules, d := rulesValue.AsStructSliceT(ctx)
	if d.HasError() {
		t.Fatalf("decoding security_group_rules: %v", d)
	}
	if len(rules) != 1 ||
		!rules[0].Protocol.Equal(types.StringValue("tcp")) ||
		!rules[0].PortRangeMin.Equal(types.Int64Value(22)) {
		t.Fatalf("state security_group_rules = %+v, want the one rule from the API response", rules)
	}
}

func configuredSecurityGroupResource(ctx context.Context, t *testing.T, baseURL string) resource.Resource {
	t.Helper()

	client := gcore.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("unit-test"),
		option.WithMaxRetries(0),
	)

	r := cloud_security_group.NewResource()
	configurable, ok := r.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatal("the cloud_security_group resource no longer implements Configure")
	}
	configureResp := &resource.ConfigureResponse{}
	configurable.Configure(ctx, resource.ConfigureRequest{ProviderData: &client}, configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", configureResp.Diagnostics)
	}
	return r
}

func securityGroupResponse(name, description string, rules []any) map[string]any {
	if rules == nil {
		rules = []any{}
	}
	return map[string]any{
		"id":                   testGroupID,
		"name":                 name,
		"description":          description,
		"region":               "Luxembourg",
		"region_id":            2,
		"project_id":           1,
		"revision_number":      1,
		"created_at":           "2026-01-01T00:00:00Z",
		"updated_at":           "2026-01-01T00:00:00Z",
		"security_group_rules": rules,
		"tags_v2":              []any{},
	}
}

func securityGroupType(ctx context.Context) tftypes.Object {
	return cloud_security_group.ResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
}

// securityGroupValue builds a resource object with every attribute null except
// the given overrides.
func securityGroupValue(ctx context.Context, t *testing.T, overrides map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	return objectValue(t, securityGroupType(ctx), overrides)
}

func securityGroupNull(ctx context.Context) tftypes.Value {
	return tftypes.NewValue(securityGroupType(ctx), nil)
}

// securityGroupRulesValue builds a one-element security_group_rules list.
func securityGroupRulesValue(ctx context.Context, t *testing.T, overrides map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	listType := securityGroupType(ctx).AttributeTypes["security_group_rules"].(tftypes.List)
	element := objectValue(t, listType.ElementType.(tftypes.Object), overrides)
	return tftypes.NewValue(listType, []tftypes.Value{element})
}

func objectValue(t *testing.T, objType tftypes.Object, overrides map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	values := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		values[name] = tftypes.NewValue(attrType, nil)
	}
	for name, value := range overrides {
		if _, ok := objType.AttributeTypes[name]; !ok {
			t.Fatalf("no attribute %q in the object type", name)
		}
		values[name] = value
	}
	return tftypes.NewValue(objType, values)
}

func merge(maps ...map[string]tftypes.Value) map[string]tftypes.Value {
	merged := map[string]tftypes.Value{}
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}
