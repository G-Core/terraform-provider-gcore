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
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
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
	volume, err := r.client.Cloud.Volumes.NewAndPoll(
		ctx,
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	// Use raw JSON from the response to unmarshal the "computed" fields into the data model
	err = apijson.UnmarshalComputed([]byte(volume.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudVolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	// Handle volume size changes (resize) separately before other updates
	sizeChanged := !data.Size.Equal(state.Size)

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

		resizedVolume, err := r.client.Cloud.Volumes.ResizeAndPoll(
			ctx,
			data.ID.ValueString(),
			resizeParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to make http request", err.Error())
			return
		}
		err = apijson.UnmarshalComputed([]byte(resizedVolume.RawJSON()), &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
			return
		}
	}

	// Check if name or tags have changed before sending update request
	nameChanged := !data.Name.Equal(state.Name)
	tagsChanged := !custom.TagsEqual(data.Tags, state.Tags)

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
	}

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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudVolumeResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}
