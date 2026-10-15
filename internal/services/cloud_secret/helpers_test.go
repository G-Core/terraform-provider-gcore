package cloud_secret_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_secret"
)

// secretValues builds a resource object whose attributes are null except for
// the overrides, so cases only state what they care about.
func secretValues(ctx context.Context, t *testing.T, overrides map[string]tftypes.Value) tftypes.Value {
	t.Helper()

	objType := cloud_secret.ResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
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

func secretNull(ctx context.Context, t *testing.T) tftypes.Value {
	t.Helper()
	return tftypes.NewValue(cloud_secret.ResourceSchema(ctx).Type().TerraformType(ctx), nil)
}

// payloadValue builds the payload tftypes object with placeholder strings -
// no real PEM material, this only exercises wiring.
func payloadValue(t *testing.T) tftypes.Value {
	t.Helper()

	ctx := context.Background()
	objType := cloud_secret.ResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
	payloadType := objType.AttributeTypes["payload"].(tftypes.Object)

	return tftypes.NewValue(payloadType, map[string]tftypes.Value{
		"certificate_wo":       tfStr("unit-test-cert"),
		"certificate_chain_wo": tfStr("unit-test-chain"),
		"private_key_wo":       tfStr("unit-test-key"),
	})
}

func tfStr(s string) tftypes.Value { return tftypes.NewValue(tftypes.String, s) }
func tfNum(n int64) tftypes.Value  { return tftypes.NewValue(tftypes.Number, n) }
