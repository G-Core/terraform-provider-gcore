// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/custom"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
	"github.com/stainless-sdks/gcore-terraform/internal/importpath"
	"github.com/stainless-sdks/gcore-terraform/internal/logging"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudFileShareResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudFileShareResource)(nil)
var _ resource.ResourceWithImportState = (*CloudFileShareResource)(nil)

func NewResource() resource.Resource {
	return &CloudFileShareResource{}
}

// CloudFileShareResource defines the resource implementation.
type CloudFileShareResource struct {
	client *gcore.Client
}

func (r *CloudFileShareResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_file_share"
}

func (r *CloudFileShareResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudFileShareResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudFileShareModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.FileShareNewParams{}

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
	fileShare, err := r.client.Cloud.FileShares.NewAndPoll(
		ctx,
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	err = apijson.UnmarshalComputed([]byte(fileShare.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	// AccessRuleIDs is read-only, so set it to null
	data.AccessRuleIDs = customfield.NullList[types.String](ctx)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudFileShareResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudFileShareModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudFileShareModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check if fields that can be updated using the Update method have changed
	if !data.Name.Equal(state.Name) || !custom.TagsEqual(data.Tags, state.Tags) || !data.ShareSettings.Equal(state.ShareSettings) {
		params := cloud.FileShareUpdateParams{}

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
		fileShare, err := r.client.Cloud.FileShares.UpdateAndPoll(
			ctx,
			data.ID.ValueString(),
			params,
			option.WithRequestBody("application/json", dataBytes),
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to make http request", err.Error())
			return
		}
		err = apijson.UnmarshalComputed([]byte(fileShare.RawJSON()), &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
			return
		}
	}

	// Check if size has changed, which requires a resize operation
	if !data.Size.Equal(state.Size) {
		params := cloud.FileShareResizeParams{Size: data.Size.ValueInt64()}
		if !data.ProjectID.IsNull() {
			params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		fileShare, err := r.client.Cloud.FileShares.ResizeAndPoll(
			ctx,
			data.ID.ValueString(),
			params,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to make http request", err.Error())
			return
		}
		err = apijson.UnmarshalComputed([]byte(fileShare.RawJSON()), &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudFileShareResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudFileShareModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.FileShareGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.FileShares.Get(
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

	// Read access_rule_ids separately using the Access Rules list endpoint
	accessRuleListParams := cloud.FileShareAccessRuleListParams{}
	if !data.ProjectID.IsNull() {
		accessRuleListParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}
	if !data.RegionID.IsNull() {
		accessRuleListParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}
	accessRules, err := r.client.Cloud.FileShares.AccessRules.List(
		ctx,
		data.ID.ValueString(),
		accessRuleListParams,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		// don't fail the read if we can't get access rules, just log a warning
		resp.Diagnostics.AddWarning("failed to make http request", err.Error())
	}
	accessRuleIDs := make([]types.String, 0, len(accessRules.Results))
	for _, rule := range accessRules.Results {
		accessRuleIDs = append(accessRuleIDs, types.StringValue(rule.ID))
	}
	data.AccessRuleIDs, _ = customfield.NewList[types.String](ctx, accessRuleIDs)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudFileShareResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudFileShareModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.FileShareDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.FileShares.DeleteAndPoll(
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

func (r *CloudFileShareResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudFileShareModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_file_share_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<file_share_id>",
		&path_project_id,
		&path_region_id,
		&path_file_share_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.ID = types.StringValue(path_file_share_id)

	res := new(http.Response)
	_, err := r.client.Cloud.FileShares.Get(
		ctx,
		path_file_share_id,
		cloud.FileShareGetParams{
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

func (r *CloudFileShareResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}
