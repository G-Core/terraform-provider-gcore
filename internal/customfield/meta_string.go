package customfield

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable                      = (*MetaStringType)(nil)
	_ basetypes.StringValuableWithSemanticEquals   = (*MetaStringValue)(nil)
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

// MetaStringValue is a custom string value for DNS meta fields.
// It wraps a basetypes.StringValue and stores plain string values.
// The encoder is responsible for serializing this to JSON strings.
type MetaStringValue struct {
	basetypes.StringValue
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

	// If both are null or both are unknown, they are equal
	if v.IsNull() && newValue.IsNull() {
		return true, diags
	}
	if v.IsUnknown() && newValue.IsUnknown() {
		return true, diags
	}

	// If one is null/unknown and the other isn't, they are not equal
	if v.IsNull() || v.IsUnknown() || newValue.IsNull() || newValue.IsUnknown() {
		return false, diags
	}

	// Compare string values
	return v.ValueString() == newValue.ValueString(), diags
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
