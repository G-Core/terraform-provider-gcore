package waap_domain_advanced_rule

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestConfigValidatorsActionCount(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := ResourceSchema(ctx)
	typ := s.Type().TerraformType(ctx).(tftypes.Object)
	actionTyp := typ.AttributeTypes["action"].(tftypes.Object)

	actionWith := func(unknown []string, kinds ...string) tftypes.Value {
		vals := map[string]tftypes.Value{}
		for name, at := range actionTyp.AttributeTypes {
			vals[name] = tftypes.NewValue(at, nil)
		}
		for _, k := range unknown {
			vals[k] = tftypes.NewValue(actionTyp.AttributeTypes[k], tftypes.UnknownValue)
		}
		for _, k := range kinds {
			at := actionTyp.AttributeTypes[k]
			switch k {
			case "block", "tag":
				inner := map[string]tftypes.Value{}
				for n, it := range at.(tftypes.Object).AttributeTypes {
					inner[n] = tftypes.NewValue(it, nil)
				}
				vals[k] = tftypes.NewValue(at, inner)
			default:
				vals[k] = tftypes.NewValue(at, map[string]tftypes.Value{})
			}
		}
		return tftypes.NewValue(actionTyp, vals)
	}

	for name, tc := range map[string]struct {
		kinds   []string
		unknown []string
		wantErr bool
	}{
		"one":              {kinds: []string{"allow"}},
		"none":             {wantErr: true},
		"two":              {kinds: []string{"allow", "block"}, wantErr: true},
		"unknown only":     {unknown: []string{"tag"}},
		"one plus unknown": {kinds: []string{"allow"}, unknown: []string{"tag"}},
		"two plus unknown": {kinds: []string{"allow", "block"}, unknown: []string{"tag"}, wantErr: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			vals := map[string]tftypes.Value{}
			for n, at := range typ.AttributeTypes {
				vals[n] = tftypes.NewValue(at, nil)
			}
			vals["action"] = actionWith(tc.unknown, tc.kinds...)

			req := resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: s, Raw: tftypes.NewValue(typ, vals)}}
			resp := &resource.ValidateConfigResponse{}
			for _, v := range (&WaapDomainAdvancedRuleResource{}).ConfigValidators(ctx) {
				v.ValidateResource(ctx, req, resp)
			}
			if resp.Diagnostics.HasError() != tc.wantErr {
				t.Errorf("wantErr=%v, got %v", tc.wantErr, resp.Diagnostics)
			}
		})
	}
}

func TestTagsValidators(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	action := ResourceSchema(ctx).Attributes["action"].(schema.SingleNestedAttribute)
	tag := action.Attributes["tag"].(schema.SingleNestedAttribute)
	validators := tag.Attributes["tags"].(schema.SetAttribute).Validators

	str := types.StringValue
	for name, tc := range map[string]struct {
		elems   []types.String
		wantErr bool
	}{
		"valid":        {elems: []types.String{str("ab C9")}},
		"unknown tag":  {elems: []types.String{types.StringUnknown()}},
		"null tag":     {elems: []types.String{types.StringNull()}, wantErr: true},
		"empty set":    {elems: []types.String{}, wantErr: true},
		"digit first":  {elems: []types.String{str("1abc")}, wantErr: true},
		"hyphen":       {elems: []types.String{str("tf-test")}, wantErr: true},
		"empty string": {elems: []types.String{str("")}, wantErr: true},
		"six":          {elems: []types.String{str("a"), str("b"), str("c"), str("d"), str("e"), str("f")}, wantErr: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			set, diags := types.SetValueFrom(ctx, types.StringType, tc.elems)
			if diags.HasError() {
				t.Fatalf("building set: %v", diags)
			}
			req := validator.SetRequest{Path: path.Root("action").AtName("tag").AtName("tags"), ConfigValue: set}
			resp := &validator.SetResponse{}
			for _, v := range validators {
				v.ValidateSet(ctx, req, resp)
			}
			if resp.Diagnostics.HasError() != tc.wantErr {
				t.Errorf("wantErr=%v, got %v", tc.wantErr, resp.Diagnostics)
			}
		})
	}
}

func TestActionDurationValidators(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	action := ResourceSchema(ctx).Attributes["action"].(schema.SingleNestedAttribute)
	block := action.Attributes["block"].(schema.SingleNestedAttribute)
	validators := block.Attributes["action_duration"].(schema.StringAttribute).Validators

	for name, tc := range map[string]struct {
		value   types.String
		wantErr bool
	}{
		"seconds":  {value: types.StringValue("30")},
		"minutes":  {value: types.StringValue("5m")},
		"days":     {value: types.StringValue("7d")},
		"null":     {value: types.StringNull()},
		"unknown":  {value: types.StringUnknown()},
		"empty":    {value: types.StringValue(""), wantErr: true},
		"zero":     {value: types.StringValue("0s"), wantErr: true},
		"bad unit": {value: types.StringValue("5w"), wantErr: true},
		"too long": {value: types.StringValue("12345678901"), wantErr: true},
		"space":    {value: types.StringValue("5 m"), wantErr: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{Path: path.Root("action").AtName("block").AtName("action_duration"), ConfigValue: tc.value}
			resp := &validator.StringResponse{}
			for _, v := range validators {
				v.ValidateString(ctx, req, resp)
			}
			if resp.Diagnostics.HasError() != tc.wantErr {
				t.Errorf("wantErr=%v, got %v", tc.wantErr, resp.Diagnostics)
			}
		})
	}
}
