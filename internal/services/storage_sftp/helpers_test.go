package storage_sftp_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_sftp"
)

// storageSftpValues builds a resource object whose attributes are null except
// for the overrides, so cases only state what they care about
func storageSftpValues(ctx context.Context, t *testing.T, overrides map[string]tftypes.Value) tftypes.Value {
	t.Helper()

	objType := storage_sftp.ResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
	values := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		values[name] = tftypes.NewValue(attrType, nil)
	}
	for name, value := range overrides {
		if _, ok := objType.AttributeTypes[name]; !ok {
			t.Fatalf("no attribute %q in the resource schema", name)
		}
		values[name] = value
	}

	return tftypes.NewValue(objType, values)
}

func storageSftpNull(ctx context.Context, t *testing.T) tftypes.Value {
	t.Helper()
	return tftypes.NewValue(storage_sftp.ResourceSchema(ctx).Type().TerraformType(ctx), nil)
}

func str(s string) tftypes.Value   { return tftypes.NewValue(tftypes.String, s) }
func num(n int64) tftypes.Value    { return tftypes.NewValue(tftypes.Number, n) }
func boolean(b bool) tftypes.Value { return tftypes.NewValue(tftypes.Bool, b) }
func ptr(b bool) *bool             { return &b }
