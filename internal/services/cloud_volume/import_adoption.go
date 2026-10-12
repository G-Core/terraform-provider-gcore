package cloud_volume

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/G-Core/terraform-provider-gcore/internal/importadoption"
)

// importAdoption marks a volume whose no_refresh create-only attributes still
// have to be adopted after `terraform import`. See internal/importadoption for
// the lifecycle. source is Required, so a configuration that does not ignore
// it adopts at least source on the first apply after import and that apply
// closes the window; under lifecycle.ignore_changes = [source] the planned
// source stays null, that apply is a no-op, and a no-op `apply -refresh=false`
// leaves the window open (Read is the only hook that advances the marker) -
// refuseAdoptionWithoutSource covers that shape.
var importAdoption = importadoption.Marker{Key: "volume_import_adoption_pending"}

// createOnlyAttributes: names must equal the tfsdk tags in model.go (the
// warning prints them). type_name is deliberately absent: the GET returns it
// as volume_type and backfillTypeName copies it, so it is compared against the
// truth, never adopted.
var createOnlyAttributes = importadoption.Table[CloudVolumeModel]{
	"source":                   func(m CloudVolumeModel) bool { return m.Source.IsNull() },
	"attachment_tag":           func(m CloudVolumeModel) bool { return m.AttachmentTag.IsNull() },
	"image_id":                 func(m CloudVolumeModel) bool { return m.ImageID.IsNull() },
	"instance_id_to_attach_to": func(m CloudVolumeModel) bool { return m.InstanceIDToAttachTo.IsNull() },
	"snapshot_id":              func(m CloudVolumeModel) bool { return m.SnapshotID.IsNull() },
}

var importAdoptionPlan = importadoption.PlanHook[CloudVolumeModel]{
	Marker:  importAdoption,
	Table:   createOnlyAttributes,
	Warning: importadoption.Warning{Noun: "volume", Transport: volumeTransport},
}

// volumeTransport: the volume PATCH is reduced to name+tags, so nothing adopted
// reaches the API; two of the values are checked against what the API reports.
func volumeTransport([]string) string {
	return "Nothing is sent to the API for them. \"image_id\" is checked against the volume's " +
		"image metadata and \"instance_id_to_attach_to\" against its attachments, and a " +
		"contradiction fails the apply; \"source\", \"snapshot_id\" and \"attachment_tag\" " +
		"cannot be checked and are trusted as-is."
}

// refuseAdoptionWithoutSource is the volume's fail-closed half of the plan-time
// adoption hook. The volume cannot verify every adopted value in Update
// (source, snapshot_id and attachment_tag are unverifiable, see
// volumeTransport), so the one update shape in which an adopting apply can
// record the other create-only attributes without ever recording source - the
// attribute every import must adopt - is refused at plan time rather than
// warned about. The shape is reachable only through lifecycle.ignore_changes
// = [source]: source is Required, so a configured value is never null, and
// core substitutes the (null, imported) prior value for an ignored attribute.
// Nothing is refused while nothing is adopted: a no-op plan with source
// ignored is the user's choice and writes nothing wrong. Never writes private
// state.
func refuseAdoptionWithoutSource(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() { // create / destroy
		return
	}
	adopting, diags := importAdoption.Pending(ctx, req.Private)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() || !adopting {
		return
	}
	var plan, state *CloudVolumeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() || plan == nil || state == nil {
		return
	}
	if adopted := adoptionWithoutSource(adopting, false, *plan, *state); len(adopted) > 0 {
		resp.Diagnostics.Append(adoptionWithoutSourceError(adopted))
	}
}

// adoptionWithoutSource is the pure decision behind refuseAdoptionWithoutSource:
// the create-only attributes an adopting update plan would record while its
// planned source is null. Empty when the marker is not live, the plan creates
// the volume, source is planned (known or unknown), or nothing is adopted.
func adoptionWithoutSource(adopting, creating bool, plan, state CloudVolumeModel) []string {
	if !adopting || creating || !plan.Source.IsNull() {
		return nil
	}
	return importadoption.Adopted(plan, state, createOnlyAttributes)
}

// adoptionWithoutSourceError names the attributes the plan would record and
// the two ways out: make the adoption explicit, or replace the volume.
func adoptionWithoutSourceError(adopted []string) diag.Diagnostic {
	quoted := make([]string, 0, len(adopted))
	for _, a := range adopted {
		quoted = append(quoted, fmt.Sprintf("%q", a))
	}
	return diag.NewAttributeErrorDiagnostic(path.Root("source"),
		"source is ignored while the imported volume is still adopting create-only attributes",
		fmt.Sprintf("\"source\" is excluded by lifecycle.ignore_changes while this imported volume's "+
			"create-only attributes are still being adopted, so this apply would record %s from "+
			"configuration without ever recording \"source\". Remove \"source\" from ignore_changes so the "+
			"adopting apply records \"source\" together with %s and the adoption is explicit, or run "+
			"`terraform apply -replace=` to recreate the volume from the configuration.",
			strings.Join(quoted, ", "), strings.Join(quoted, ", ")))
}

// verifyAdoption runs at the top of an adopting Update, BEFORE any mutation,
// against the prior state's computed fields (fresh from the refresh that
// preceded this apply; under -refresh=false they are as old as the last one).
// Exact where the API always reports the fact, mismatch-only where it is a hint.
func verifyAdoption(ctx context.Context, plan, state CloudVolumeModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if state.ImageID.IsNull() && !plan.ImageID.IsNull() && !plan.ImageID.IsUnknown() &&
		!state.VolumeImageMetadata.IsNull() && !state.VolumeImageMetadata.IsUnknown() {
		// Hint only: cinder propagates volume_image_metadata through snapshot
		// lineage, so absence proves nothing and only a DIFFERENT id is evidence.
		if meta, d := state.VolumeImageMetadata.Value(ctx); !d.HasError() {
			if got, ok := meta["image_id"]; ok && !got.IsNull() && !got.IsUnknown() && got.ValueString() != plan.ImageID.ValueString() {
				diags.AddAttributeError(path.Root("image_id"),
					"Configured image_id contradicts the imported volume",
					fmt.Sprintf("The volume's image metadata names image %q, not %q. Fix the "+
						"configuration, or run `terraform apply -replace=` to recreate the volume "+
						"from the configured image.", got.ValueString(), plan.ImageID.ValueString()))
			}
		}
	}

	if state.InstanceIDToAttachTo.IsNull() && !plan.InstanceIDToAttachTo.IsNull() && !plan.InstanceIDToAttachTo.IsUnknown() {
		attached := false
		if !state.Attachments.IsNull() && !state.Attachments.IsUnknown() {
			if atts, d := state.Attachments.AsStructSliceT(ctx); !d.HasError() {
				for _, a := range atts {
					attached = attached || a.ServerID.ValueString() == plan.InstanceIDToAttachTo.ValueString()
				}
			}
		}
		if !attached {
			diags.AddAttributeError(path.Root("instance_id_to_attach_to"),
				"Configured instance_id_to_attach_to contradicts the imported volume",
				fmt.Sprintf("The volume is not attached to instance %q. Remove the attribute, "+
					"attach the volume out of band, or run `terraform apply -replace=`.",
					plan.InstanceIDToAttachTo.ValueString()))
		}
	}
	return diags
}

// backfillTypeName: the request field type_name and the response field
// volume_type are the same value under two names. Null-only, so a value already
// in state (including one recorded by 2.0.0-rc.1 without being applied) is
// never overwritten and state == config keeps planning no changes. Call after
// every decode of a GET/PATCH body (Create, Update, Read, ImportState); a
// planned unknown is resolved here because state cannot hold one.
func backfillTypeName(data *CloudVolumeModel) {
	if data == nil || !(data.TypeName.IsNull() || data.TypeName.IsUnknown()) {
		return
	}
	if data.VolumeType.IsNull() || data.VolumeType.IsUnknown() || data.VolumeType.ValueString() == "" {
		data.TypeName = types.StringNull() // resolve a planned unknown
		return
	}
	data.TypeName = data.VolumeType
}

// normaliseTypeNameCase gives a type_name recorded in another case (accepted by
// earlier versions) the API's case. Read only: in Create/Update the planned
// value must come back unchanged.
func normaliseTypeNameCase(data *CloudVolumeModel) {
	if data == nil || data.TypeName.IsNull() || data.TypeName.IsUnknown() ||
		data.VolumeType.IsNull() || data.VolumeType.IsUnknown() {
		return
	}
	if strings.EqualFold(data.TypeName.ValueString(), data.VolumeType.ValueString()) {
		data.TypeName = data.VolumeType
	}
}

// warnRecordedTypeNameDrift is marker-independent: it names a recorded type_name
// that the refreshed volume_type contradicts (a 2.0.0-rc.1 adoption that was
// never applied). It plans nothing: type_name is Computed, so omitting it keeps
// the recorded value. Setting it to the reported volume_type corrects the
// record in place without an API call (typeNameReplaces compares against
// volume_type). Fires on every plan until then, except the correcting one.
func warnRecordedTypeNameDrift(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var state *CloudVolumeModel
	if req.State.Get(ctx, &state).HasError() || state == nil {
		return
	}
	if state.TypeName.IsNull() || state.TypeName.IsUnknown() || state.VolumeType.IsNull() || state.VolumeType.IsUnknown() {
		return
	}
	if strings.EqualFold(state.TypeName.ValueString(), state.VolumeType.ValueString()) {
		return
	}
	// This plan already corrects the record.
	var planned types.String
	if !req.Plan.GetAttribute(ctx, path.Root("type_name"), &planned).HasError() &&
		!planned.IsUnknown() && strings.EqualFold(planned.ValueString(), state.VolumeType.ValueString()) {
		return
	}
	resp.Diagnostics.AddAttributeWarning(path.Root("type_name"),
		"Recorded type_name does not match the volume",
		fmt.Sprintf("State records type_name %[1]q but the API reports volume_type %[2]q. The recorded "+
			"value was accepted by an earlier provider version without being applied. Set type_name "+
			"to %[2]q to correct the record in place; omitting it keeps the recorded value and this "+
			"warning. Any other value is planned against the real type %[2]q.",
			state.TypeName.ValueString(), state.VolumeType.ValueString()))
}
