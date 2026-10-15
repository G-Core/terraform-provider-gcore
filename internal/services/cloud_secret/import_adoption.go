package cloud_secret

import (
	"github.com/G-Core/terraform-provider-gcore/internal/importadoption"
)

// secretImportAdoption marks a secret whose no_refresh attributes still have
// to be adopted from config into state after `terraform import`. See
// internal/importadoption for the lifecycle.
//
// Two attributes carry the marker-keyed modifier, and both are needed:
// "payload" is a *CloudSecretPayloadModel (a pointer), tagged no_refresh, so
// ImportState/Read never populate it and it stays a nil pointer - a real
// null-to-non-null transition when the config sets it, not a no-op masked by
// write-only nulling. "payload_wo_version" is the sibling
// version counter apijson never touches at all (no json tag): the API cannot
// echo it, so it is also null after import. Either one alone still forces a
// replacement on the imported plan - confirmed against real infra, not just
// read off the schema. They carry the marker-keyed modifier through the
// local wrappers in plan_modifiers.go, not the bare planmodifiers.* import.
//
// payload and payload_wo_version are Optional. State holds null for them
// after an import without material, after the import itself, and after an
// in-place removal. Neither the import-without-material case outliving its
// window nor the removal opens a window, so the marker (not Required-ness)
// is what keeps a later set on the replacement path. Removal is in place
// (see plan_modifiers.go), and the create/replace requirement lives in
// payload_rules.go. Known consequence: after an import without payload,
// setting payload before the window closes is adopted with a warning, not
// replaced.
var secretImportAdoption = importadoption.Marker{Key: "cloud_secret_import_adoption_pending"}

// createOnlyAttributes: names must equal the tfsdk tags on CloudSecretModel,
// since the shared warning prints them.
var createOnlyAttributes = importadoption.Table[CloudSecretModel]{
	"payload":            func(m CloudSecretModel) bool { return m.Payload == nil },
	"payload_wo_version": func(m CloudSecretModel) bool { return m.PayloadWoVersion.IsNull() },
}

// secretImportAdoptionPlan is the whole ModifyPlan body: this resource has no
// other plan-time adoption concern.
var secretImportAdoptionPlan = importadoption.PlanHook[CloudSecretModel]{
	Marker:  secretImportAdoption,
	Table:   createOnlyAttributes,
	Warning: importadoption.Warning{Noun: "secret", Transport: secretTransport},
}

// secretTransport replaces the default paragraph, which says the configured
// values are recorded: payload's fields are write-only, so only its presence
// and payload_wo_version reach state, and editing the material alone is not
// a change Terraform can see.
func secretTransport([]string) string {
	return "Nothing is sent to the API. The certificate and key in \"payload\" are write-only and are never " +
		"stored: only that a payload is set and its \"payload_wo_version\" are recorded, so editing the " +
		"material alone changes nothing. Bump \"payload_wo_version\" whenever the material changes."
}
