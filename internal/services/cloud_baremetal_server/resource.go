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
	"github.com/G-Core/terraform-provider-gcore/internal/custom"
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
	res := new(http.Response)
	_, err = r.client.Cloud.Baremetal.Servers.NewAndPoll(
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
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
		data.Tags = tags
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

	// Read the adoption marker before anything branches on it: it gates the
	// rebuild below and is retired by the single state write at the end.
	adopting, adoptionDiags := importAdoptionPending(ctx, req.Private)
	resp.Diagnostics.Append(adoptionDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Whether this apply rebuilds the server is decided before anything is sent,
	// so the one answer that cannot be acted on stops the apply while it is
	// still a no-op. Terraform has resolved every planned value by the time it
	// calls Update, so an unknown image_id or user_data arriving here is a
	// contradiction rather than a change: rebuilding on it would reinstall the
	// operating system while sending the API an empty value, and quietly not
	// rebuilding would break the promise that the plan-time warning and the
	// apply agree.
	rebuild := rebuildsOnApply(*data, *state, adopting)

	if rebuild == rebuildUnknown {
		resp.Diagnostics.Append(rebuildAtUnknownValueError())
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
		if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
			data.Tags = tags
		}
		stateHasChanged = true
	}

	// A change to image_id or user_data is applied by rebuilding the server,
	// which reinstalls the operating system and erases the disks. ModifyPlan
	// asks rebuildsOnApply the same question to decide whether to warn first.
	if rebuild == rebuildYes {
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
		res := new(http.Response)
		_, err := r.client.Cloud.Baremetal.Servers.RebuildAndPoll(
			ctx,
			data.ID.ValueString(),
			params,
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

	if stateHasChanged {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	} else {
		// No API changes were made, but no_refresh fields may need to be
		// populated in state from the plan. This handles the post-import
		// scenario: after import those fields are null in state because the API
		// does not return them. The one-time update-in-place persists the
		// config values to state so that future changes trigger replacement
		// correctly - and without it Terraform rejects the apply outright with
		// "Provider produced inconsistent result after apply", because the
		// planned adoption never reached state.
		//
		// The whole model is deliberately NOT written here: no request was
		// made, so computed attributes are still unknown and writing them would
		// produce an invalid state object.
		for attribute, value := range map[string]any{
			"interfaces":     data.Interfaces,
			"apptemplate_id": data.ApptemplateID,
			"app_config":     data.AppConfig,
			"flavor":         data.Flavor,
			"image_id":       data.ImageID,
			"name_template":  data.NameTemplate,
			"ssh_key_name":   data.SSHKeyName,
			"user_data":      data.UserData,
			"username":       data.Username,
		} {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attribute), value)...)
		}
	}

	// Clear only after the write succeeded: a marker cleared before success
	// would turn a transient failure into a destroy+recreate on the retry,
	// because the next plan would no longer recognise the import.
	if adopting && !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(clearImportAdoptionPending(ctx, resp.Private)...)
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

	if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
		data.Tags = tags
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

	if tags, ok := custom.ConvertAPITagsToCustomfieldMap(ctx, bytes); ok {
		data.Tags = tags
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

	// apptemplate_id, name_template, username, app_config and user_data are
	// create-only and tagged no_refresh, so the read above cannot populate them
	// and they stay null in state. Mark the adoption window so the plan
	// modifiers on those attributes know this null came from an import rather
	// than from a server created without them - the first apply adopts the
	// configured values into state and clears the marker. The marker also stops
	// that first apply from rebuilding the server, see Update.
	resp.Diagnostics.Append(markImportAdoptionPending(ctx, resp.Private)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudBaremetalServerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// The adoption marker is deliberately NOT cleared here. A plan whose only
	// effect would be the clear is a no-op plan, and Terraform never persists
	// PlannedPrivate for no-op plans, so the clear would silently evaporate.
	// Update owns retiring the marker.
	//
	// What this does do is warn while an adoption is actually planned. Adopted
	// values are written to state without being sent to the API, and no later
	// refresh can correct them, so this is the only moment a user can catch a
	// wrong value before it becomes permanent.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	adopting, diags := importAdoptionPending(ctx, req.Private)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var plan, state *CloudBaremetalServerModel

	// Both warnings are advisory, so a decode problem must never fail the plan.
	if req.Plan.Get(ctx, &plan).HasError() || req.State.Get(ctx, &state).HasError() {
		return
	}

	if plan == nil || state == nil {
		return
	}

	// A replacement runs Create, not Update: nothing is adopted and nothing is
	// rebuilt, so both warnings would be describing the opposite of what is
	// about to happen.
	//
	// This only catches the replacement it can see, which is the kind an
	// attribute change forces. A -replace on the command line is invisible here,
	// so both warnings are worded to stay true when they are emitted on a plan
	// the CLI goes on to print as a replacement.
	if planForcesReplacement(ctx, req.Plan, req.State, req.Config, adopting) {
		return
	}

	// Changing image_id or user_data is applied by reinstalling the operating
	// system, and Terraform plans that as an ordinary in-place update. This is
	// the only place the user is told the disks are about to be erased.
	//
	// A value that is still unknown warns too. Whether it ends up differing is
	// not decided yet, and a plan that might erase the disks is worth the same
	// sentence as one that will - staying quiet until the apply would be telling
	// the user only once it is too late to look.
	if rebuild := rebuildsOnApply(*plan, *state, adopting); rebuild == rebuildYes || rebuild == rebuildUnknown {
		resp.Diagnostics.Append(rebuildWarning())
	}

	if !adopting {
		return
	}

	if adopted := adoptedCreateOnlyAttributes(*plan, *state); len(adopted) > 0 {
		resp.Diagnostics.Append(importAdoptionWarning(adopted))
	}
}
