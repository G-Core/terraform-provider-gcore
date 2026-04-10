// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster

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
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/custom"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudGPUBaremetalClusterResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudGPUBaremetalClusterResource)(nil)
var _ resource.ResourceWithImportState = (*CloudGPUBaremetalClusterResource)(nil)

func NewResource() resource.Resource {
	return &CloudGPUBaremetalClusterResource{}
}

// CloudGPUBaremetalClusterResource defines the resource implementation.
type CloudGPUBaremetalClusterResource struct {
	client *gcore.Client
}

func (r *CloudGPUBaremetalClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_gpu_baremetal_cluster"
}

func (r *CloudGPUBaremetalClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudGPUBaremetalClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudGPUBaremetalClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve servers_settings from config (for no_refresh fields like credentials)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("servers_settings"), &data.ServersSettings)...)

	params := cloud.GPUBaremetalClusterNewParams{}

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
	cluster, err := r.client.Cloud.GPUBaremetal.Clusters.NewAndPoll(
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
	if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, []byte(cluster.RawJSON())); ok {
		data.Tags = tags
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudGPUBaremetalClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudGPUBaremetalClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudGPUBaremetalClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stateHasChanged := false

	// Check if name or tags have changed.
	// Skip unknown values — they indicate computed fields, not user changes.
	nameChanged := !data.Name.IsUnknown() && !data.Name.Equal(state.Name)
	tagsChanged := !data.Tags.IsUnknown() && !data.Tags.Equal(state.Tags)

	if nameChanged || tagsChanged {
		params := cloud.GPUBaremetalClusterUpdateParams{}
		if !data.ProjectID.IsNull() {
			params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		// Build a model that only differs from state in name/tags, so
		// MarshalJSONForUpdate produces a patch with only those fields.
		updateData := *state
		updateData.Name = data.Name
		updateData.Tags = data.Tags
		dataBytes, err := updateData.MarshalJSONForUpdate(*state)
		if err != nil {
			resp.Diagnostics.AddError("failed to serialize http request", err.Error())
			return
		}
		res := new(http.Response)
		_, err = r.client.Cloud.GPUBaremetal.Clusters.Update(
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
		stateHasChanged = true
	}

	// Check if servers count has changed
	if !data.ServersCount.IsNull() && data.ServersCount.ValueInt64() != state.ServersCount.ValueInt64() {
		params := cloud.GPUBaremetalClusterResizeParams{
			InstancesCount: data.ServersCount.ValueInt64(),
		}
		if !data.ProjectID.IsNull() {
			params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		cluster, err := r.client.Cloud.GPUBaremetal.Clusters.ResizeAndPoll(
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
		if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, []byte(cluster.RawJSON())); ok {
			data.Tags = tags
		}
		stateHasChanged = true
	}

	// Check if server settings that require UpdateServersSettings + Rebuild have changed
	// (image_id, credentials, user_data)
	imageChanged := !data.ImageID.IsNull() && data.ImageID.ValueString() != state.ImageID.ValueString()
	credentialsChanged := credentialsHaveChanged(data.ServersSettings, state.ServersSettings)
	userDataChanged := false
	if data.ServersSettings != nil && state.ServersSettings != nil {
		userDataChanged = !data.ServersSettings.UserData.Equal(state.ServersSettings.UserData)
	}

	if imageChanged || credentialsChanged || userDataChanged {
		updateParams := cloud.GPUBaremetalClusterUpdateServersSettingsParams{}
		if !data.ProjectID.IsNull() {
			updateParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			updateParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		if imageChanged {
			updateParams.ImageID = param.NewOpt(data.ImageID.ValueString())
		}
		if data.ServersSettings != nil && data.ServersSettings.Credentials != nil {
			creds := data.ServersSettings.Credentials
			if !creds.SSHKeyName.IsNull() {
				updateParams.ServersSettings.Credentials.SSHKeyName = param.NewOpt(creds.SSHKeyName.ValueString())
			}
		}
		if data.ServersSettings != nil && !data.ServersSettings.UserData.IsNull() && !data.ServersSettings.UserData.IsUnknown() {
			updateParams.ServersSettings.UserData = param.NewOpt(data.ServersSettings.UserData.ValueString())
		}

		_, err := r.client.Cloud.GPUBaremetal.Clusters.UpdateServersSettings(
			ctx,
			data.ID.ValueString(),
			updateParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to make http request", err.Error())
			return
		}

		// Rebuild to apply changes to existing servers
		rebuildParams := cloud.GPUBaremetalClusterRebuildParams{}
		if !data.ProjectID.IsNull() {
			rebuildParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			rebuildParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		cluster, err := r.client.Cloud.GPUBaremetal.Clusters.RebuildAndPoll(
			ctx,
			data.ID.ValueString(),
			rebuildParams,
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
		if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, []byte(cluster.RawJSON())); ok {
			data.Tags = tags
		}
		stateHasChanged = true
	}

	if stateHasChanged {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *CloudGPUBaremetalClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudGPUBaremetalClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.GPUBaremetalClusterGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.GPUBaremetal.Clusters.Get(
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudGPUBaremetalClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudGPUBaremetalClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.GPUBaremetalClusterDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.GPUBaremetal.Clusters.DeleteAndPoll(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
}

func (r *CloudGPUBaremetalClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudGPUBaremetalClusterModel)

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
	_, err := r.client.Cloud.GPUBaremetal.Clusters.Get(
		ctx,
		path_cluster_id,
		cloud.GPUBaremetalClusterGetParams{
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

	// Workaround: Fields with no_refresh tag are skipped during Unmarshal.
	// For import, we need to manually extract these from the raw API response.
	var rawResponse struct {
		ServersSettings struct {
			UserData   string `json:"user_data"`
			SSHKeyName string `json:"ssh_key_name"`
			Username   string `json:"username"`
		} `json:"servers_settings"`
	}
	if err := json.Unmarshal(bytes, &rawResponse); err == nil {
		if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
			data.Tags = tags
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
			data.ServersSettings.Credentials = &CloudGPUBaremetalClusterServersSettingsCredentialsModel{}
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

func (r *CloudGPUBaremetalClusterResource) ModifyPlan(_ context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// No validation needed on create or destroy.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan, state CloudGPUBaremetalClusterModel
	resp.Diagnostics.Append(req.Plan.Get(context.Background(), &plan)...)
	resp.Diagnostics.Append(req.State.Get(context.Background(), &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The UpdateServersSettings API only supports updating ssh_key_name.
	// Reject username/password changes at plan time with a clear error.
	if credentialsHaveChanged(plan.ServersSettings, state.ServersSettings) &&
		usernameOrPasswordChanged(plan.ServersSettings, state.ServersSettings) {
		resp.Diagnostics.AddError(
			"Unsupported credential update",
			"Updating username or password is not supported. These credentials can only be set during resource creation. To change them, the resource must be recreated.",
		)
	}
}

// credentialsHaveChanged returns true if any credential field has changed between plan and state.
//
// Compared fields: SSHKeyName, Username, PasswordWoVersion (which tracks password_wo changes).
// If new credential fields are added to CloudGPUBaremetalClusterServersSettingsCredentialsModel,
// this function must be updated to include them.
func credentialsHaveChanged(plan, state *CloudGPUBaremetalClusterServersSettingsModel) bool {
	if plan == nil || plan.Credentials == nil {
		return false
	}
	if state == nil || state.Credentials == nil {
		return true
	}
	return !plan.Credentials.SSHKeyName.Equal(state.Credentials.SSHKeyName) ||
		!plan.Credentials.Username.Equal(state.Credentials.Username) ||
		!plan.Credentials.PasswordWoVersion.Equal(state.Credentials.PasswordWoVersion)
}

// usernameOrPasswordChanged returns true if specifically the username or password
// credential fields have changed (as opposed to ssh_key_name).
// The UpdateServersSettings API only supports updating ssh_key_name, so changes to
// username or password must be rejected.
func usernameOrPasswordChanged(plan, state *CloudGPUBaremetalClusterServersSettingsModel) bool {
	if plan == nil || plan.Credentials == nil {
		return false
	}
	if state == nil || state.Credentials == nil {
		// New credentials being set — check if they include username/password
		return !plan.Credentials.Username.IsNull() || !plan.Credentials.PasswordWoVersion.IsNull()
	}
	return !plan.Credentials.Username.Equal(state.Credentials.Username) ||
		!plan.Credentials.PasswordWoVersion.Equal(state.Credentials.PasswordWoVersion)
}
