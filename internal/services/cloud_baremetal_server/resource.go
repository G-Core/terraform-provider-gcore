// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudBaremetalServerResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudBaremetalServerResource)(nil)
var _ resource.ResourceWithImportState = (*CloudBaremetalServerResource)(nil)

func NewResource() resource.Resource {
	return &CloudBaremetalServerResource{}
}

// CloudBaremetalServerResource defines the resource implementation.
type CloudBaremetalServerResource struct {
	client *gcore.Client
}

func (r *CloudBaremetalServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_baremetal_server"
}

func (r *CloudBaremetalServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudBaremetalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudBaremetalServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Write-only attributes are not available from the plan; read from config.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo"), &data.Password)...)

	params := cloud.BaremetalServerNewParams{}

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
	baremetalServer, err := r.client.Cloud.Baremetal.Servers.NewAndPoll(
		ctx,
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	err = apijson.UnmarshalComputed([]byte(baremetalServer.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	// The API returns tags as an array of objects [{key, value, read_only}] which
	// cannot be deserialized into a Map[string]. Resolve any remaining unknown
	// tags to null so Terraform doesn't report unknown values after apply.
	if data.Tags.IsUnknown() {
		data.Tags = customfield.NullMap[types.String](ctx)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudBaremetalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudBaremetalServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudBaremetalServerModel

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
		params := cloud.BaremetalServerUpdateParams{}

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
		_, err = r.client.Cloud.Baremetal.Servers.Update(
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

	// Check if image_id or user_data have changed, if so perform a server rebuild operation.
	imageChanged := !data.ImageID.IsNull() && data.ImageID.ValueString() != state.ImageID.ValueString()
	userDataChanged := !data.UserData.IsNull() && data.UserData.ValueString() != state.UserData.ValueString()

	if imageChanged || userDataChanged {
		params := cloud.BaremetalServerRebuildParams{}
		if !data.ProjectID.IsNull() {
			params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		if !data.ImageID.IsNull() {
			params.ImageID = param.NewOpt(data.ImageID.ValueString())
		}
		if !data.UserData.IsNull() {
			params.UserData = param.NewOpt(data.UserData.ValueString())
		}
		server, err := r.client.Cloud.Baremetal.Servers.RebuildAndPoll(
			ctx,
			data.ID.ValueString(),
			params,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to make http request", err.Error())
			return
		}
		err = apijson.UnmarshalComputed([]byte(server.RawJSON()), &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
			return
		}
		stateHasChanged = true
	}

	if stateHasChanged {
		// The API returns tags as an array of objects [{key, value, read_only}] which
		// cannot be deserialized into a Map[string]. Resolve any remaining unknown
		// tags to null so Terraform doesn't report unknown values after apply.
		if data.Tags.IsUnknown() {
			data.Tags = customfield.NullMap[types.String](ctx)
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	} else {
		// No API changes were made, but no_refresh fields (like interfaces) may
		// need to be populated in state from the plan. This handles the post-import
		// scenario: after import, interfaces is null in state because the API doesn't
		// return it. The one-time update-in-place persists the config value to state
		// so that future changes trigger replacement correctly.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("interfaces"), data.Interfaces)...)
	}
}

func (r *CloudBaremetalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudBaremetalServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.BaremetalServerGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.Baremetal.Servers.Get(
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

func (r *CloudBaremetalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudBaremetalServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.BaremetalServerDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.Baremetal.Servers.DeleteAndPoll(
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

func (r *CloudBaremetalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudBaremetalServerModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_server_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<server_id>",
		&path_project_id,
		&path_region_id,
		&path_server_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.ID = types.StringValue(path_server_id)

	res := new(http.Response)
	_, err := r.client.Cloud.Baremetal.Servers.Get(
		ctx,
		path_server_id,
		cloud.BaremetalServerGetParams{
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
		Tags []struct {
			Key      string `json:"key"`
			Value    string `json:"value"`
			ReadOnly bool   `json:"read_only"`
		} `json:"tags"`
		Flavor struct {
			FlavorName string `json:"flavor_name"`
		} `json:"flavor"`
		Metadata struct {
			ImageID string `json:"image_id"`
		} `json:"metadata"`
		SSHKeyName string `json:"ssh_key_name"`
	}
	if err := json.Unmarshal(bytes, &rawResponse); err == nil {
		// tags: API returns array of {key, value, read_only} objects, but model expects map[string]string.
		// Filter out read-only tags (system-managed) so they don't cause drift.
		tagsMap := make(map[string]types.String)
		for _, tag := range rawResponse.Tags {
			if !tag.ReadOnly {
				tagsMap[tag.Key] = types.StringValue(tag.Value)
			}
		}
		if len(tagsMap) > 0 {
			data.Tags = customfield.NewMapMust[types.String](ctx, tagsMap)
		} else {
			data.Tags = customfield.NullMap[types.String](ctx)
		}
		// flavor: API returns nested object, model expects string.
		if rawResponse.Flavor.FlavorName != "" {
			data.Flavor = types.StringValue(rawResponse.Flavor.FlavorName)
		}
		// image_id: stored in metadata, not as a top-level field in the response.
		if rawResponse.Metadata.ImageID != "" {
			data.ImageID = types.StringValue(rawResponse.Metadata.ImageID)
		}
		// ssh_key_name: API returns keypair UUID, look up the human-readable name.
		if rawResponse.SSHKeyName != "" {
			sshKey, sshErr := r.client.Cloud.SSHKeys.Get(
				ctx,
				rawResponse.SSHKeyName,
				cloud.SSHKeyGetParams{
					ProjectID: param.NewOpt(path_project_id),
				},
			)
			if sshErr == nil {
				data.SSHKeyName = types.StringValue(sshKey.Name)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudBaremetalServerResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}
