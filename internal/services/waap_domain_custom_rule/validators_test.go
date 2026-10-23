package waap_domain_custom_rule

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestCanonicalIP(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		value   types.String
		wantErr bool
	}{
		"ipv4":             {value: types.StringValue("192.0.2.10")},
		"ipv6":             {value: types.StringValue("2001:db8::1")},
		"ipv4-mapped ipv6": {value: types.StringValue("::ffff:192.0.2.1")},
		"null":             {value: types.StringNull()},
		"unknown":          {value: types.StringUnknown()},
		"uppercase ipv6":   {value: types.StringValue("2001:DB8::1"), wantErr: true},
		"expanded ipv6":    {value: types.StringValue("2001:db8:0::1"), wantErr: true},
		"cidr":             {value: types.StringValue("192.0.2.0/24"), wantErr: true},
		"zone":             {value: types.StringValue("fe80::1%eth0"), wantErr: true},
		"hostname":         {value: types.StringValue("example.com"), wantErr: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{Path: path.Root("ip"), ConfigValue: tc.value}
			resp := &validator.StringResponse{}
			canonicalIP{}.ValidateString(context.Background(), req, resp)
			if resp.Diagnostics.HasError() != tc.wantErr {
				t.Errorf("wantErr=%v, got %v", tc.wantErr, resp.Diagnostics)
			}
		})
	}
}

func TestNormalizedListValidators(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	conditions := ResourceSchema(ctx).Attributes["conditions"].(schema.ListNestedAttribute)
	// validate runs the validators of a string list or set on one value.
	validate := func(t *testing.T, kind, attribute, value string) diag.Diagnostics {
		var attr schema.Attribute
		if kind == "action.tag" {
			action := ResourceSchema(ctx).Attributes["action"].(schema.SingleNestedAttribute)
			attr = action.Attributes["tag"].(schema.SingleNestedAttribute).Attributes[attribute]
		} else {
			attr = conditions.NestedObject.Attributes[kind].(schema.SingleNestedAttribute).Attributes[attribute]
		}
		p := path.Root(kind).AtName(attribute)
		switch a := attr.(type) {
		case schema.ListAttribute:
			list, diags := types.ListValueFrom(ctx, types.StringType, []string{value})
			if diags.HasError() {
				t.Fatalf("building list: %v", diags)
			}
			resp := &validator.ListResponse{}
			for _, v := range a.Validators {
				v.ValidateList(ctx, validator.ListRequest{Path: p, ConfigValue: list}, resp)
			}
			return resp.Diagnostics
		case schema.SetAttribute:
			set, diags := types.SetValueFrom(ctx, types.StringType, []string{value})
			if diags.HasError() {
				t.Fatalf("building set: %v", diags)
			}
			resp := &validator.SetResponse{}
			for _, v := range a.Validators {
				v.ValidateSet(ctx, validator.SetRequest{Path: p, ConfigValue: set}, resp)
			}
			return resp.Diagnostics
		default:
			t.Fatalf("%s.%s is %T, not a list or set", kind, attribute, attr)
			return nil
		}
	}

	for name, tc := range map[string]struct {
		kind, attribute string
		value           string
		wantErr         bool
	}{
		"country upper":   {kind: "country", attribute: "country_code", value: "US"},
		"country lower":   {kind: "country", attribute: "country_code", value: "us", wantErr: true},
		"ja3 lower":       {kind: "ja3", attribute: "ja3_fingerprints", value: "e7d705a3286e19ea42f587b344ee6865"},
		"ja3 upper":       {kind: "ja3", attribute: "ja3_fingerprints", value: "E7D705A3286E19EA42F587B344EE6865", wantErr: true},
		"ja4 lower":       {kind: "ja4", attribute: "ja4_fingerprints", value: "t13d3113h2_e8f1e7e78f70_ce5650b735ce"},
		"ja4 upper":       {kind: "ja4", attribute: "ja4_fingerprints", value: "T13D3113H2_E8F1E7E78F70_CE5650B735CE", wantErr: true},
		"rate ip":         {kind: "request_rate", attribute: "ips", value: "2001:db8::1"},
		"rate ip uncanon": {kind: "request_rate", attribute: "ips", value: "2001:DB8::1", wantErr: true},
		// The API trims these values.
		"content type":       {kind: "content_type", attribute: "content_type", value: "text/html"},
		"content type space": {kind: "content_type", attribute: "content_type", value: "text/html ", wantErr: true},
		"tag":                {kind: "tags", attribute: "tags", value: "xss"},
		"tag space":          {kind: "tags", attribute: "tags", value: "xss ", wantErr: true},
		"user tag":           {kind: "user_defined_tags", attribute: "tags", value: "qa"},
		"user tag space":     {kind: "user_defined_tags", attribute: "tags", value: "qa ", wantErr: true},
		"action tag":         {kind: "action.tag", attribute: "tags", value: "lead"},
		"action tag space":   {kind: "action.tag", attribute: "tags", value: "lead ", wantErr: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if diags := validate(t, tc.kind, tc.attribute, tc.value); diags.HasError() != tc.wantErr {
				t.Errorf("wantErr=%v, got %v", tc.wantErr, diags)
			}
		})
	}
}

func TestConfigValidatorsActionCount(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := ResourceSchema(ctx)
	typ := s.Type().TerraformType(ctx).(tftypes.Object)
	actionTyp := typ.AttributeTypes["action"].(tftypes.Object)

	actionWith := func(kinds ...string) tftypes.Value {
		vals := map[string]tftypes.Value{}
		for name, at := range actionTyp.AttributeTypes {
			vals[name] = tftypes.NewValue(at, nil)
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
		wantErr bool
	}{
		"one":  {kinds: []string{"monitor"}},
		"none": {wantErr: true},
		"two":  {kinds: []string{"allow", "block"}, wantErr: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			vals := map[string]tftypes.Value{}
			for n, at := range typ.AttributeTypes {
				vals[n] = tftypes.NewValue(at, nil)
			}
			vals["action"] = actionWith(tc.kinds...)

			req := resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: s, Raw: tftypes.NewValue(typ, vals)}}
			resp := &resource.ValidateConfigResponse{}
			for _, v := range (&WaapDomainCustomRuleResource{}).ConfigValidators(ctx) {
				v.ValidateResource(ctx, req, resp)
			}
			if resp.Diagnostics.HasError() != tc.wantErr {
				t.Errorf("wantErr=%v, got %v", tc.wantErr, resp.Diagnostics)
			}
		})
	}
}

// Every condition must set exactly one condition type.
func TestOneConditionKind(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	conditions := ResourceSchema(ctx).Attributes["conditions"].(schema.ListNestedAttribute)
	objType := conditions.NestedObject.Type().(types.ObjectType)

	condition := func(kinds ...string) types.Object {
		vals := map[string]attr.Value{}
		for name, at := range objType.AttrTypes {
			vals[name] = types.ObjectNull(at.(types.ObjectType).AttrTypes)
		}
		for _, k := range kinds {
			kindType := objType.AttrTypes[k].(types.ObjectType)
			inner := map[string]attr.Value{}
			for n, it := range kindType.AttrTypes {
				v, err := it.ValueFromTerraform(ctx, tftypes.NewValue(it.TerraformType(ctx), nil))
				if err != nil {
					t.Fatalf("building %s.%s: %v", k, n, err)
				}
				inner[n] = v
			}
			vals[k] = types.ObjectValueMust(kindType.AttrTypes, inner)
		}
		return types.ObjectValueMust(objType.AttrTypes, vals)
	}

	for name, tc := range map[string]struct {
		kinds   []string
		wantErr bool
	}{
		"one":  {kinds: []string{"ip"}},
		"none": {wantErr: true},
		"two":  {kinds: []string{"ip", "url"}, wantErr: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.ObjectRequest{Path: path.Root("conditions").AtListIndex(0), ConfigValue: condition(tc.kinds...)}
			resp := &validator.ObjectResponse{}
			oneConditionKind{}.ValidateObject(ctx, req, resp)
			if resp.Diagnostics.HasError() != tc.wantErr {
				t.Errorf("wantErr=%v, got %v", tc.wantErr, resp.Diagnostics)
			}
		})
	}
}

func TestNoSurroundingSpace(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		value   types.String
		wantErr bool
	}{
		"plain":         {value: types.StringValue("Block admin")},
		"inner space":   {value: types.StringValue("a b")},
		"null":          {value: types.StringNull()},
		"unknown":       {value: types.StringUnknown()},
		"leading space": {value: types.StringValue(" rule"), wantErr: true},
		"trailing tab":  {value: types.StringValue("rule\t"), wantErr: true},
		"newline":       {value: types.StringValue("rule\n"), wantErr: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{Path: path.Root("name"), ConfigValue: tc.value}
			resp := &validator.StringResponse{}
			noSurroundingSpace{}.ValidateString(context.Background(), req, resp)
			if resp.Diagnostics.HasError() != tc.wantErr {
				t.Errorf("wantErr=%v, got %v", tc.wantErr, resp.Diagnostics)
			}
		})
	}
}

// The API matches enum values case-sensitively, so a lowercase value must fail
// at plan time instead of with a 422 on apply.
func TestEnumsAreCaseSensitive(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	conditions := ResourceSchema(ctx).Attributes["conditions"].(schema.ListNestedAttribute)
	method := conditions.NestedObject.Attributes["http_method"].(schema.SingleNestedAttribute)
	validators := method.Attributes["http_method"].(schema.StringAttribute).Validators

	for value, wantErr := range map[string]bool{"POST": false, "post": true} {
		req := validator.StringRequest{Path: path.Root("http_method"), ConfigValue: types.StringValue(value)}
		resp := &validator.StringResponse{}
		for _, v := range validators {
			v.ValidateString(ctx, req, resp)
		}
		if resp.Diagnostics.HasError() != wantErr {
			t.Errorf("%q: wantErr=%v, got %v", value, wantErr, resp.Diagnostics)
		}
	}
}
