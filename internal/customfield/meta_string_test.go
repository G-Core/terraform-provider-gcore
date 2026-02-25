package customfield

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestMetaStringType_ValueFromTerraform(t *testing.T) {
	ctx := context.Background()
	metaType := MetaStringType{}

	tests := []struct {
		name     string
		input    tftypes.Value
		expected string
		isNull   bool
	}{
		{
			name:     "plain string",
			input:    tftypes.NewValue(tftypes.String, "https://example.com"),
			expected: "https://example.com",
		},
		{
			name:     "string with special chars",
			input:    tftypes.NewValue(tftypes.String, `hello "world"`),
			expected: `hello "world"`,
		},
		{
			name:   "null value",
			input:  tftypes.NewValue(tftypes.String, nil),
			isNull: true,
		},
		{
			name:     "empty string",
			input:    tftypes.NewValue(tftypes.String, ""),
			expected: "",
		},
		{
			name:     "json-like string",
			input:    tftypes.NewValue(tftypes.String, `{"key": "value"}`),
			expected: `{"key": "value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := metaType.ValueFromTerraform(ctx, tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			metaValue, ok := result.(MetaStringValue)
			if !ok {
				t.Fatalf("expected MetaStringValue, got %T", result)
			}

			if tt.isNull {
				if !metaValue.IsNull() {
					t.Error("expected null value")
				}
				return
			}

			if metaValue.IsNull() {
				t.Error("unexpected null value")
				return
			}

			if metaValue.ValueString() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, metaValue.ValueString())
			}
		})
	}
}

func TestMetaStringValue_Type(t *testing.T) {
	ctx := context.Background()
	value := NewMetaStringValue("test")

	attrType := value.Type(ctx)
	_, ok := attrType.(MetaStringType)
	if !ok {
		t.Errorf("expected MetaStringType, got %T", attrType)
	}
}

func TestMetaStringValue_Equal(t *testing.T) {
	tests := []struct {
		name     string
		a        MetaStringValue
		b        attr.Value
		expected bool
	}{
		{
			name:     "equal strings",
			a:        NewMetaStringValue("test"),
			b:        NewMetaStringValue("test"),
			expected: true,
		},
		{
			name:     "unequal strings",
			a:        NewMetaStringValue("test1"),
			b:        NewMetaStringValue("test2"),
			expected: false,
		},
		{
			name:     "both null",
			a:        MetaStringNull(),
			b:        MetaStringNull(),
			expected: true,
		},
		{
			name:     "both unknown",
			a:        MetaStringUnknown(),
			b:        MetaStringUnknown(),
			expected: true,
		},
		{
			name:     "different types",
			a:        NewMetaStringValue("test"),
			b:        basetypes.NewStringValue("test"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Equal(tt.b)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMetaStringValue_StringSemanticEquals(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		a        MetaStringValue
		b        MetaStringValue
		expected bool
	}{
		{
			name:     "equal strings",
			a:        NewMetaStringValue("test"),
			b:        NewMetaStringValue("test"),
			expected: true,
		},
		{
			name:     "unequal strings",
			a:        NewMetaStringValue("test1"),
			b:        NewMetaStringValue("test2"),
			expected: false,
		},
		{
			name:     "both null",
			a:        MetaStringNull(),
			b:        MetaStringNull(),
			expected: true,
		},
		{
			name:     "both unknown",
			a:        MetaStringUnknown(),
			b:        MetaStringUnknown(),
			expected: true,
		},
		{
			name:     "one null one known",
			a:        MetaStringNull(),
			b:        NewMetaStringValue("test"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, diags := tt.a.StringSemanticEquals(ctx, tt.b)
			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMetaStringConstructors(t *testing.T) {
	t.Run("MetaStringNull", func(t *testing.T) {
		v := MetaStringNull()
		if !v.IsNull() {
			t.Error("expected null")
		}
	})

	t.Run("MetaStringUnknown", func(t *testing.T) {
		v := MetaStringUnknown()
		if !v.IsUnknown() {
			t.Error("expected unknown")
		}
	})

	t.Run("NewMetaStringValue", func(t *testing.T) {
		v := NewMetaStringValue("test")
		if v.IsNull() || v.IsUnknown() {
			t.Error("expected known value")
		}
		if v.ValueString() != "test" {
			t.Errorf("expected 'test', got %q", v.ValueString())
		}
	})

	t.Run("NewMetaStringPointerValue with value", func(t *testing.T) {
		s := "test"
		v := NewMetaStringPointerValue(&s)
		if v.IsNull() || v.IsUnknown() {
			t.Error("expected known value")
		}
		if v.ValueString() != "test" {
			t.Errorf("expected 'test', got %q", v.ValueString())
		}
	})

	t.Run("NewMetaStringPointerValue with nil", func(t *testing.T) {
		v := NewMetaStringPointerValue(nil)
		if !v.IsNull() {
			t.Error("expected null")
		}
	})
}

func TestMetaStringType_String(t *testing.T) {
	metaType := MetaStringType{}
	if metaType.String() != "MetaStringType" {
		t.Errorf("expected 'MetaStringType', got %q", metaType.String())
	}
}

func TestMetaStringType_Equal(t *testing.T) {
	t.Run("equal types", func(t *testing.T) {
		a := MetaStringType{}
		b := MetaStringType{}
		if !a.Equal(b) {
			t.Error("expected equal")
		}
	})

	t.Run("different types", func(t *testing.T) {
		a := MetaStringType{}
		b := basetypes.StringType{}
		if a.Equal(b) {
			t.Error("expected not equal")
		}
	})
}
