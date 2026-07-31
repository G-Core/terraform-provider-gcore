package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var attrChangeTestSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Required: true,
		},
		"tags": schema.MapAttribute{
			Computed:    true,
			Optional:    true,
			ElementType: types.StringType,
		},
		"servers_count": schema.Int64Attribute{
			Required: true,
		},
		"status": schema.StringAttribute{
			Computed: true,
		},
		"has_pending_changes": schema.BoolAttribute{
			Computed: true,
		},
	},
}

var attrChangeSchemaType = tftypes.Object{
	AttributeTypes: map[string]tftypes.Type{
		"id":                  tftypes.String,
		"name":                tftypes.String,
		"tags":                tftypes.Map{ElementType: tftypes.String},
		"servers_count":       tftypes.Number,
		"status":              tftypes.String,
		"has_pending_changes": tftypes.Bool,
	},
}

type attrChangeValues struct {
	name   any // string or nil (null); tftypes.UnknownValue for unknown
	tags   any // map[string]string, nil, or tftypes.UnknownValue
	count  any // int, nil, or tftypes.UnknownValue
	status any // string, nil, or tftypes.UnknownValue
}

func attrChangeRaw(v attrChangeValues) tftypes.Value {
	tagsType := tftypes.Map{ElementType: tftypes.String}
	var tagsVal tftypes.Value
	switch tv := v.tags.(type) {
	case map[string]string:
		elems := map[string]tftypes.Value{}
		for k, s := range tv {
			elems[k] = tftypes.NewValue(tftypes.String, s)
		}
		tagsVal = tftypes.NewValue(tagsType, elems)
	default:
		tagsVal = tftypes.NewValue(tagsType, v.tags)
	}

	var countVal tftypes.Value
	switch cv := v.count.(type) {
	case int:
		countVal = tftypes.NewValue(tftypes.Number, cv)
	default:
		countVal = tftypes.NewValue(tftypes.Number, v.count)
	}

	return tftypes.NewValue(attrChangeSchemaType, map[string]tftypes.Value{
		"id":                  tftypes.NewValue(tftypes.String, "cluster-1"),
		"name":                tftypes.NewValue(tftypes.String, v.name),
		"tags":                tagsVal,
		"servers_count":       countVal,
		"status":              tftypes.NewValue(tftypes.String, v.status),
		"has_pending_changes": tftypes.NewValue(tftypes.Bool, false),
	})
}

func TestStringUseStateUnlessAttributesChange(t *testing.T) {
	t.Parallel()

	watched := []string{"name", "tags", "servers_count"}
	stateVals := attrChangeValues{name: "cluster", tags: map[string]string{"env": "test"}, count: 2, status: "active"}

	makeReq := func(configVals attrChangeValues) planmodifier.StringRequest {
		return planmodifier.StringRequest{
			StateValue: types.StringValue("active"),
			PlanValue:  types.StringUnknown(),
			Config:     tfsdk.Config{Raw: attrChangeRaw(configVals), Schema: attrChangeTestSchema},
			State:      tfsdk.State{Raw: attrChangeRaw(stateVals), Schema: attrChangeTestSchema},
			Path:       path.Root("status"),
		}
	}

	run := func(req planmodifier.StringRequest) *planmodifier.StringResponse {
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
		planmodifiers.StringUseStateUnlessAttributesChange(watched...).PlanModifyString(context.Background(), req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
		}
		return resp
	}

	t.Run("no state (create) - stays unknown", func(t *testing.T) {
		t.Parallel()
		req := makeReq(stateVals)
		req.State = tfsdk.State{Raw: tftypes.NewValue(attrChangeSchemaType, nil), Schema: attrChangeTestSchema}
		req.StateValue = types.StringNull()
		if resp := run(req); !resp.PlanValue.IsUnknown() {
			t.Errorf("expected unknown, got %v", resp.PlanValue)
		}
	})

	t.Run("watched attrs unchanged - preserves state", func(t *testing.T) {
		t.Parallel()
		// config: same name/tags/count as state; computed attrs null in config
		if resp := run(makeReq(attrChangeValues{name: "cluster", tags: map[string]string{"env": "test"}, count: 2, status: nil})); !resp.PlanValue.Equal(types.StringValue("active")) {
			t.Errorf("expected preserved state, got %v", resp.PlanValue)
		}
	})

	t.Run("watched attr null in config - treated as unchanged", func(t *testing.T) {
		t.Parallel()
		// tags omitted from config (computed_optional keeps prior state)
		if resp := run(makeReq(attrChangeValues{name: "cluster", tags: nil, count: 2, status: nil})); !resp.PlanValue.Equal(types.StringValue("active")) {
			t.Errorf("expected preserved state, got %v", resp.PlanValue)
		}
	})

	t.Run("watched attr changed - stays unknown", func(t *testing.T) {
		t.Parallel()
		if resp := run(makeReq(attrChangeValues{name: "renamed", tags: map[string]string{"env": "test"}, count: 2, status: nil})); !resp.PlanValue.IsUnknown() {
			t.Errorf("expected unknown, got %v", resp.PlanValue)
		}
	})

	t.Run("watched map attr changed - stays unknown", func(t *testing.T) {
		t.Parallel()
		if resp := run(makeReq(attrChangeValues{name: "cluster", tags: map[string]string{"env": "prod"}, count: 2, status: nil})); !resp.PlanValue.IsUnknown() {
			t.Errorf("expected unknown, got %v", resp.PlanValue)
		}
	})

	t.Run("watched attr unknown in config - stays unknown", func(t *testing.T) {
		t.Parallel()
		if resp := run(makeReq(attrChangeValues{name: tftypes.UnknownValue, tags: map[string]string{"env": "test"}, count: 2, status: nil})); !resp.PlanValue.IsUnknown() {
			t.Errorf("expected unknown, got %v", resp.PlanValue)
		}
	})

	t.Run("plan value already known - not overridden", func(t *testing.T) {
		t.Parallel()
		req := makeReq(attrChangeValues{name: "cluster", tags: nil, count: 2, status: nil})
		req.PlanValue = types.StringValue("degraded")
		if resp := run(req); !resp.PlanValue.Equal(types.StringValue("degraded")) {
			t.Errorf("expected known plan value kept, got %v", resp.PlanValue)
		}
	})

	t.Run("null state value - stays unknown", func(t *testing.T) {
		t.Parallel()
		req := makeReq(attrChangeValues{name: "cluster", tags: nil, count: 2, status: nil})
		req.StateValue = types.StringNull()
		if resp := run(req); !resp.PlanValue.IsUnknown() {
			t.Errorf("expected unknown, got %v", resp.PlanValue)
		}
	})
}

func TestBoolUseStateUnlessAttributesChange(t *testing.T) {
	t.Parallel()

	stateVals := attrChangeValues{name: "cluster", tags: nil, count: 2, status: "active"}

	makeReq := func(configVals attrChangeValues) planmodifier.BoolRequest {
		return planmodifier.BoolRequest{
			StateValue: types.BoolValue(false),
			PlanValue:  types.BoolUnknown(),
			Config:     tfsdk.Config{Raw: attrChangeRaw(configVals), Schema: attrChangeTestSchema},
			State:      tfsdk.State{Raw: attrChangeRaw(stateVals), Schema: attrChangeTestSchema},
			Path:       path.Root("has_pending_changes"),
		}
	}

	run := func(req planmodifier.BoolRequest) *planmodifier.BoolResponse {
		resp := &planmodifier.BoolResponse{PlanValue: req.PlanValue}
		planmodifiers.BoolUseStateUnlessAttributesChange("name", "servers_count").PlanModifyBool(context.Background(), req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
		}
		return resp
	}

	t.Run("watched attrs unchanged - preserves state", func(t *testing.T) {
		t.Parallel()
		if resp := run(makeReq(attrChangeValues{name: "cluster", tags: nil, count: 2, status: nil})); !resp.PlanValue.Equal(types.BoolValue(false)) {
			t.Errorf("expected preserved state, got %v", resp.PlanValue)
		}
	})

	t.Run("watched attr changed - stays unknown", func(t *testing.T) {
		t.Parallel()
		if resp := run(makeReq(attrChangeValues{name: "cluster", tags: nil, count: 3, status: nil})); !resp.PlanValue.IsUnknown() {
			t.Errorf("expected unknown, got %v", resp.PlanValue)
		}
	})
}
