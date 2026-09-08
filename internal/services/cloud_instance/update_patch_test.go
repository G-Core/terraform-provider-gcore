package cloud_instance

import (
	"strings"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestConfigurationSurvivesComputedRefresh covers the refresh every write path
// runs after its API call. configuration is create-only and never echoed, so
// the response carries no key for it and the configured value has to survive
// untouched - anything else is a state that contradicts the plan.
func TestConfigurationSurvivesComputedRefresh(t *testing.T) {
	t.Parallel()

	configuration := map[string]customfield.MetaStringValue{
		"qa_key": customfield.NewMetaStringValue("qa_value"),
	}
	model := &CloudInstanceModel{Configuration: &configuration}

	if err := apijson.UnmarshalComputed([]byte(`{"id":"abc","status":"ACTIVE"}`), &model); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := (*model.Configuration)["qa_key"].ValueString(); got != "qa_value" {
		t.Errorf("expected %q to survive the refresh, got %q", "qa_value", got)
	}
}

// TestUpdatePatchOmitsUnconfiguredName covers an instance created from
// name_template: config has no name, state carries the one the API assigned.
// Serializing that pair sends "name": null, which the API rejects.
func TestUpdatePatchOmitsUnconfiguredName(t *testing.T) {
	t.Parallel()

	state := CloudInstanceModel{Name: types.StringValue("tf-nt-10-0-0-1")}
	plan := state
	plan.Name = types.StringNull()

	if !plan.Name.IsNull() && !plan.Name.Equal(state.Name) {
		t.Error("an unconfigured name must not count as a rename")
	}

	patch := plan
	if patch.Name.IsNull() {
		patch.Name = state.Name
	}

	body, err := patch.MarshalJSONForUpdate(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(string(body), `"name"`) {
		t.Errorf("expected no name in the patch body, got %s", body)
	}
}
