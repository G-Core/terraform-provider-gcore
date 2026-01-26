// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
	"github.com/stainless-sdks/gcore-terraform/internal/importpath"
	"github.com/stainless-sdks/gcore-terraform/internal/logging"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudGPUVirtualClusterResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudGPUVirtualClusterResource)(nil)
var _ resource.ResourceWithImportState = (*CloudGPUVirtualClusterResource)(nil)

func NewResource() resource.Resource {
	return &CloudGPUVirtualClusterResource{}
}

// CloudGPUVirtualClusterResource defines the resource implementation.
type CloudGPUVirtualClusterResource struct {
	client *gcore.Client
}

func (r *CloudGPUVirtualClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_gpu_virtual_cluster"
}

func (r *CloudGPUVirtualClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudGPUVirtualClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudGPUVirtualClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("servers_settings"), &data.ServersSettings)...)

	params := cloud.GPUVirtualClusterNewParams{}

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
	cluster, err := r.client.Cloud.GPUVirtual.Clusters.NewAndPoll(
		ctx,
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	err = apijson.UnmarshalComputed([]byte(cluster.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudGPUVirtualClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudGPUVirtualClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudGPUVirtualClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stateHasChanged := false

	// Check for a cluster name change
	if !data.Name.IsNull() && data.Name.ValueString() != state.Name.ValueString() {
		params := cloud.GPUVirtualClusterUpdateParams{}

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
		_, err = r.client.Cloud.GPUVirtual.Clusters.Update(
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
		stateHasChanged = true
	}

	// Check if tags have changed
	if !data.Tags.IsNull() && !data.Tags.Equal(state.Tags) {
		params := cloud.GPUVirtualClusterActionParams{}
		if !data.ProjectID.IsNull() {
			params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		tagsMap := make(map[string]string)
		tagsValue, _ := data.Tags.Value(ctx)
		for k, v := range tagsValue {
			tagsMap[k] = v.ValueString()
		}
		params.OfUpdateTags = &cloud.GPUVirtualClusterActionParamsBodyUpdateTags{
			Tags: tagsMap,
		}
		cluster, err := r.client.Cloud.GPUVirtual.Clusters.ActionAndPoll(
			ctx,
			data.ID.ValueString(),
			params,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to make http request", err.Error())
			return
		}
		err = apijson.UnmarshalComputed([]byte(cluster.RawJSON()), &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
			return
		}
		stateHasChanged = true
	}

	// Check if servers count has changed
	if !data.ServersCount.IsNull() && data.ServersCount.ValueInt64() != state.ServersCount.ValueInt64() {
		params := cloud.GPUVirtualClusterActionParams{}

		if !data.ProjectID.IsNull() {
			params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		params.OfResize = &cloud.GPUVirtualClusterActionParamsBodyResize{
			ServersCount: data.ServersCount.ValueInt64(),
		}
		cluster, err := r.client.Cloud.GPUVirtual.Clusters.ActionAndPoll(
			ctx,
			data.ID.ValueString(),
			params,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to make http request", err.Error())
			return
		}
		err = apijson.UnmarshalComputed([]byte(cluster.RawJSON()), &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
			return
		}
		stateHasChanged = true
	}

	if stateHasChanged {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *CloudGPUVirtualClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudGPUVirtualClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.GPUVirtualClusterGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.GPUVirtual.Clusters.Get(
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

func (r *CloudGPUVirtualClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudGPUVirtualClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.GPUVirtualClusterDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.GPUVirtual.Clusters.DeleteAndPoll(
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

func (r *CloudGPUVirtualClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudGPUVirtualClusterModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_cluster_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<cluster_id>",
		&path_project_id,
		&path_region_id,
		&path_cluster_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.ID = types.StringValue(path_cluster_id)

	res := new(http.Response)
	_, err := r.client.Cloud.GPUVirtual.Clusters.Get(
		ctx,
		path_cluster_id,
		cloud.GPUVirtualClusterGetParams{
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

	// Workaround: API doesn't return source field in volumes, infer it from image_id
	if data.ServersSettings != nil && data.ServersSettings.Volumes != nil {
		for _, vol := range *data.ServersSettings.Volumes {
			if vol.ImageID.IsNull() || vol.ImageID.ValueString() == "" {
				vol.Source = types.StringValue("new")
			} else {
				vol.Source = types.StringValue("image")
			}
		}
	}

	// Workaround: Fields with no_refresh tag are skipped during Unmarshal.
	// For import, we need to manually extract these from the raw API response.
	var rawResponse struct {
		Tags []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"tags"`
		ServersSettings struct {
			UserData   string `json:"user_data"`
			SSHKeyName string `json:"ssh_key_name"`
			Username   string `json:"username"`
		} `json:"servers_settings"`
	}
	if err := json.Unmarshal(bytes, &rawResponse); err == nil {
		// tags: API returns array of {key, value} objects, but model expects map[string]string.
		if len(rawResponse.Tags) > 0 {
			tagsMap := make(map[string]types.String, len(rawResponse.Tags))
			for _, tag := range rawResponse.Tags {
				tagsMap[tag.Key] = types.StringValue(tag.Value)
			}
			data.Tags = customfield.NewMapMust[types.String](ctx, tagsMap)
		}

		// user_data: API returns plain text, but users send base64 encoded.
		// Encode the API response so state matches what config would produce.
		if rawResponse.ServersSettings.UserData != "" && data.ServersSettings != nil {
			encoded := base64.StdEncoding.EncodeToString([]byte(rawResponse.ServersSettings.UserData))
			data.ServersSettings.UserData = types.StringValue(encoded)
		}

		// credentials: API returns ssh_key_name and username at servers_settings level,
		// but model expects them under credentials. password is write-only and not returned.
		if (rawResponse.ServersSettings.SSHKeyName != "" || rawResponse.ServersSettings.Username != "") && data.ServersSettings != nil {
			data.ServersSettings.Credentials = &CloudGPUVirtualClusterServersSettingsCredentialsModel{}
			if rawResponse.ServersSettings.SSHKeyName != "" {
				data.ServersSettings.Credentials.SSHKeyName = types.StringValue(rawResponse.ServersSettings.SSHKeyName)
			}
			if rawResponse.ServersSettings.Username != "" {
				data.ServersSettings.Credentials.Username = types.StringValue(rawResponse.ServersSettings.Username)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudGPUVirtualClusterResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}
