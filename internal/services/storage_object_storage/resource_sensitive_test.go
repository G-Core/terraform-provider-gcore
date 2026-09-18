package storage_object_storage_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_object_storage"
)

func TestResourceSecretKeySensitive(t *testing.T) {
	t.Parallel()

	attrs := storage_object_storage.ResourceSchema(context.Background()).Attributes
	keys := attrs["access_keys"].(schema.ListNestedAttribute)
	secret := keys.NestedObject.Attributes["secret_key"].(schema.StringAttribute)
	if !secret.Sensitive {
		t.Error("access_keys.secret_key must be sensitive")
	}
	if !secret.Computed || secret.WriteOnly {
		t.Error("access_keys.secret_key must remain computed and stored in state")
	}
	if keys.Sensitive || keys.NestedObject.Attributes["access_key"].IsSensitive() {
		t.Error("access_keys and access_key must remain non-sensitive")
	}
}
