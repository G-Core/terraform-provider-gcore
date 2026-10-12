// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/custom"
	"github.com/G-Core/terraform-provider-gcore/internal/importadoption"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudVolumeResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudVolumeResource)(nil)
var _ resource.ResourceWithImportState = (*CloudVolumeResource)(nil)

func NewResource() resource.Resource {
	return &CloudVolumeResource{}
}

// CloudVolumeResource defines the resource implementation.
type CloudVolumeResource struct {
	client *gcore.Client
}

func (r *CloudVolumeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_volume"
}

func (r *CloudVolumeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*gcore.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"unexpected resource configure type",
			fmt.Sprintf("Expected *gcore.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *CloudVolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudVolumeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.VolumeNewParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	dataBytes, err := data.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}
	res := new(http.Response)
	_, err = r.client.Cloud.Volumes.NewAndPoll(
		ctx,
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	// Use raw JSON from the response to unmarshal the "computed" fields into the data model
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
		data.Tags = tags
	}
	backfillTypeName(data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudVolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Import-adoption window (internal/importadoption): settle re-arms the
	// marker on an errored apply and clears it on a successful one. Installed
	// before the first fallible step so every error raised in this method
	// re-arms. Every non-error exit below writes plan state.
	adopting, settle := importAdoption.BeginUpdate(ctx, req.Private, resp.Private, &resp.Diagnostics)
	defer settle()
	if resp.Diagnostics.HasError() {
		return
	}

	var data *CloudVolumeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudVolumeModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if adopting {
		// Verify before the resize/PATCH branches so a mismatch never follows a real change.
		resp.Diagnostics.Append(verifyAdoption(ctx, *data, *state)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// type_name changes that reach Update are retypable (see
	// TypeNameReplaceUnlessRetypable); finish the retype before a resize.
	retyped, diags := r.retypeIfChanged(ctx, data, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if retyped {
		// resp.State still holds the prior state, which an error below would
		// save as is. Record the new type now so a failed resize or PATCH
		// cannot hide the retype.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("type_name"), data.TypeName)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("volume_type"), data.TypeName)...)
	}

	// Handle volume size changes (resize) separately before other updates
	sizeChanged := !data.Size.IsUnknown() && !data.Size.Equal(state.Size)

	if sizeChanged {
		resizeParams := cloud.VolumeResizeParams{
			Size: data.Size.ValueInt64(),
		}

		if !data.ProjectID.IsNull() {
			resizeParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}

		if !data.RegionID.IsNull() {
			resizeParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		res := new(http.Response)
		_, err := r.client.Cloud.Volumes.ResizeAndPoll(
			ctx,
			data.ID.ValueString(),
			resizeParams,
			option.WithResponseBodyInto(&res),
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to make http request", err.Error())
			return
		}
		bytes, _ := io.ReadAll(res.Body)
		err = apijson.UnmarshalComputed(bytes, &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
			return
		}
		if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
			data.Tags = tags
		}
	}

	// Check if name or tags have changed before sending update request
	nameChanged := !data.Name.Equal(state.Name)
	tagsChanged := !data.Tags.IsUnknown() && !data.Tags.Equal(state.Tags)

	// Only send update request if name or tags changed
	if nameChanged || tagsChanged {
		params := cloud.VolumeUpdateParams{}

		if !data.ProjectID.IsNull() {
			params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}

		if !data.RegionID.IsNull() {
			params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		dataBytes, err := data.MarshalJSONForUpdate(*state)
		if err == nil {
			// PATCH /v1/volumes accepts name and tags only; keep adopted create-only
			// attributes (state null, plan set) off the wire.
			dataBytes, err = importadoption.KeepKeys(dataBytes, "name", "tags")
		}
		if err != nil {
			resp.Diagnostics.AddError("failed to serialize http request", err.Error())
			return
		}
		res := new(http.Response)
		_, err = r.client.Cloud.Volumes.Update(
			ctx,
			data.ID.ValueString(),
			params,
			option.WithRequestBody("application/json", dataBytes),
			option.WithResponseBodyInto(&res),
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to make http request", err.Error())
			return
		}
		bytes, _ := io.ReadAll(res.Body)
		err = apijson.UnmarshalComputed(bytes, &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
			return
		}
		if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
			data.Tags = tags
		}
	}

	// Any change entering Update marks computed plan values as unknown, and
	// an adopting apply after import may send nothing at all. Always refresh
	// from the API so every computed field is known before saving state;
	// type_name follows volume_type where it is still unset.
	getParams := cloud.VolumeGetParams{}

	if !data.ProjectID.IsNull() {
		getParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		getParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.Volumes.Get(
		ctx,
		data.ID.ValueString(),
		getParams,
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
		data.Tags = tags
	}
	backfillTypeName(data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudVolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudVolumeModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.VolumeGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.Volumes.Get(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if res != nil && res.StatusCode == 404 {
		resp.Diagnostics.AddWarning("Resource not found", "The resource was not found on the server and will be removed from state.")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.Unmarshal(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
		data.Tags = tags
	}
	backfillTypeName(data)
	normaliseTypeNameCase(data)

	// Advance the import-adoption marker (internal/importadoption). Must run on
	// every successful Read: Terraform persists this refresh in `terraform apply`
	// even when the plan is a no-op - the only hook that retires a marker whose
	// import adopted nothing.
	resp.Diagnostics.Append(importAdoption.Advance(ctx, resp.Private)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudVolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudVolumeModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.VolumeDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.Volumes.DeleteAndPoll(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudVolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudVolumeModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_volume_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<volume_id>",
		&path_project_id,
		&path_region_id,
		&path_volume_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.ID = types.StringValue(path_volume_id)

	res := new(http.Response)
	_, err := r.client.Cloud.Volumes.Get(
		ctx,
		path_volume_id,
		cloud.VolumeGetParams{
			ProjectID: param.NewOpt(path_project_id),
			RegionID:  param.NewOpt(path_region_id),
		},
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.Unmarshal(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
		data.Tags = tags
	}
	backfillTypeName(data)

	// source, attachment_tag, image_id, instance_id_to_attach_to and snapshot_id
	// are no_refresh and stay null here; mark the adoption window so their plan
	// modifiers know this null came from an import.
	resp.Diagnostics.Append(importAdoption.Mark(ctx, resp.Private)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudVolumeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Never writes private state (PlannedPrivate is dropped on no-op plans);
	// retirement is Read's Advance. The guard errors, the hook only warns, and
	// the hook stays silent once the guard has erred.
	refuseAdoptionWithoutSource(ctx, req, resp)
	importAdoptionPlan.ModifyPlan(ctx, req, resp)
	warnRecordedTypeNameDrift(ctx, req, resp)
	checkCreatePrerequisites(ctx, req, resp)
}
