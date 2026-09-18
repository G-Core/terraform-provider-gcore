package storage_access_key_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_access_key"
)

func TestResourceSecretKeySensitive(t *testing.T) {
	t.Parallel()

	attrs := storage_access_key.ResourceSchema(context.Background()).Attributes
	secret := attrs["secret_key"].(schema.StringAttribute)
	if !secret.Sensitive {
		t.Error("secret_key must be sensitive")
	}
	if !secret.Computed || secret.WriteOnly {
		t.Error("secret_key must remain computed and stored in state")
	}
	if attrs["access_key"].IsSensitive() {
		t.Error("access_key must remain non-sensitive")
	}
}
