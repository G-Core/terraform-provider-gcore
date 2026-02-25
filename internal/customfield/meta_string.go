package customfield

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable                    = (*MetaStringType)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*MetaStringValue)(nil)
	_ json.Unmarshaler                           = (*MetaStringValue)(nil)
)

// MetaStringType is a custom string type for DNS meta fields that accepts plain HCL strings
// without requiring jsonencode(). This improves UX by allowing users to write:
//
//	meta = {
//	  webhook = "https://example.com"
//	}
//
// Instead of:
//
//	meta = {
//	  webhook = jsonencode("https://example.com")
//	}
//
// The type stores the string value as-is. The encoder handles JSON serialization
// by outputting the value as a proper JSON string (with quotes).
type MetaStringType struct {
	basetypes.StringType
}

func (t MetaStringType) Equal(o attr.Type) bool {
	other, ok := o.(MetaStringType)
	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t MetaStringType) String() string {
	return "MetaStringType"
}

func (t MetaStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return MetaStringNull(), diags
	}
	if in.IsUnknown() {
		return MetaStringUnknown(), diags
	}

	return MetaStringValue{StringValue: in}, diags
}

func (t MetaStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t MetaStringType) ValueType(ctx context.Context) attr.Value {
	return MetaStringValue{}
}

// metaStringMarker is a zero-size type that breaks structural convertibility
// between MetaStringValue and timetypes.RFC3339. Without this marker, Go considers
// MetaStringValue convertible to timetypes.RFC3339 (both embed basetypes.StringValue),
// causing apijson's decoder to route MetaStringValue through the custom time decoder
// instead of the json.Unmarshaler path.
type metaStringMarker struct{}

// MetaStringValue is a custom string value for DNS meta fields.
// It wraps a basetypes.StringValue and stores plain string values.
//
// The marker field breaks type convertibility with timetypes.RFC3339, ensuring that:
//   - Encoding uses the CustomMarshaler interface (checked before time types in apijson)
//   - Decoding uses the json.Unmarshaler interface (now reachable since time check doesn't match)
type MetaStringValue struct {
	basetypes.StringValue
	_ metaStringMarker // breaks ConvertibleTo(timetypes.RFC3339{})
}

func (v MetaStringValue) Type(ctx context.Context) attr.Type {
	return MetaStringType{}
}

func (v MetaStringValue) Equal(o attr.Value) bool {
	other, ok := o.(MetaStringValue)
	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals implements semantic equality for MetaStringValue.
// Two values are semantically equal if their string contents are equal.
func (v MetaStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(MetaStringValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: MetaStringValue\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)
		return false, diags
	}

	if v.IsNull() && newValue.IsNull() {
		return true, diags
	}
	if v.IsUnknown() && newValue.IsUnknown() {
		return true, diags
	}

	if v.IsNull() || v.IsUnknown() || newValue.IsNull() || newValue.IsUnknown() {
		return false, diags
	}

	return normalizeJSON(v.ValueString()) == normalizeJSON(newValue.ValueString()), diags
}

// MarshalJSONWithState implements apijson.CustomMarshaler.
// This is checked BEFORE time type checks in apijson's encoder, so no encoder changes needed.
func (v MetaStringValue) MarshalJSONWithState(plan any, state any) ([]byte, error) {
	p := plan.(MetaStringValue)
	if p.IsNull() || p.IsUnknown() {
		return nil, nil
	}
	val := p.ValueString()
	raw := []byte(val)
	if json.Valid(raw) {
		return raw, nil // pass through valid JSON (numbers, booleans, objects, arrays)
	}
	return json.Marshal(val) // encode as JSON string
}

// UnmarshalJSON implements json.Unmarshaler.
// With the marker field breaking timetypes.RFC3339 convertibility, the apijson decoder
// reaches the json.Unmarshaler check instead of routing through the custom time decoder.
func (v *MetaStringValue) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*v = MetaStringNull()
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*v = NewMetaStringValue(s)
		return nil
	}
	// Non-string JSON (numbers, arrays, objects) → store raw
	*v = NewMetaStringValue(string(data))
	return nil
}

// normalizeJSON round-trips a string through Go's JSON encoder to produce
// canonical form. For example, "[50.0,30.0]" becomes "[50,30]".
// Non-JSON strings are returned as-is.
func normalizeJSON(s string) string {
	raw := []byte(s)
	if !json.Valid(raw) {
		return s
	}
	var parsed interface{}
	if err := json.Unmarshal(raw, &parsed); err == nil {
		if normalized, err := json.Marshal(parsed); err == nil {
			return string(normalized)
		}
	}
	return s
}

// MetaStringNull creates a null MetaStringValue.
func MetaStringNull() MetaStringValue {
	return MetaStringValue{StringValue: basetypes.NewStringNull()}
}

// MetaStringUnknown creates an unknown MetaStringValue.
func MetaStringUnknown() MetaStringValue {
	return MetaStringValue{StringValue: basetypes.NewStringUnknown()}
}

// NewMetaStringValue creates a known MetaStringValue from a string.
func NewMetaStringValue(value string) MetaStringValue {
	return MetaStringValue{StringValue: basetypes.NewStringValue(value)}
}

// NewMetaStringPointerValue creates a MetaStringValue from a string pointer.
// If the pointer is nil, it returns a null value.
func NewMetaStringPointerValue(value *string) MetaStringValue {
	if value == nil {
		return MetaStringNull()
	}
	return NewMetaStringValue(*value)
}
