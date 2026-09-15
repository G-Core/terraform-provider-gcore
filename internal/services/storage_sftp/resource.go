// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_sftp

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/storage"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*StorageSftpResource)(nil)
var _ resource.ResourceWithModifyPlan = (*StorageSftpResource)(nil)
var _ resource.ResourceWithImportState = (*StorageSftpResource)(nil)

func NewResource() resource.Resource {
	return &StorageSftpResource{}
}

// StorageSftpResource defines the resource implementation.
type StorageSftpResource struct {
	client *gcore.Client
}

func (r *StorageSftpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_sftp"
}

func (r *StorageSftpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StorageSftpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *StorageSftpModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Write-only values are null in the plan, so read the password from config.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo"), &data.PasswordWo)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataBytes, err := data.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}

	// POST requires a mode even when no password is set.
	action, _ := passwordActionFor(true, data.PasswordWo, types.Int64Null(), types.Int64Null(), types.BoolNull())
	passwordMode := passwordModeNone
	if action == passwordSet {
		passwordMode = passwordModeSet
	}

	res := new(http.Response)
	_, err = r.client.Storage.SftpStorages.NewAndPoll(
		ctx,
		storage.SftpStorageNewParams{},
		option.WithRequestBody("application/json", dataBytes),
		option.WithJSONSet("password_mode", passwordMode),
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

	// The API echoes the password. Do not store it in state.
	data.PasswordWo = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageSftpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *StorageSftpModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *StorageSftpModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dataBytes, err := data.MarshalJSONForUpdate(*state)
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}

	// Read the resolved write-only value from config. Keep it out of data so it
	// cannot enter state or the default PATCH body.
	var configPassword types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo"), &configPassword)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action, _ := passwordActionFor(false, configPassword, data.PasswordWoVersion, state.PasswordWoVersion, state.HasPassword)
	// this guard should be unreachable: passwordSet requires a non-null config
	// password, and config values are known at apply time. It stays as a cheap
	// invariant check on passwordActionFor
	if action == passwordSet && (configPassword.IsNull() || configPassword.IsUnknown()) {
		resp.Diagnostics.AddAttributeError(
			path.Root("password_wo_version"),
			missingPasswordSummary,
			`"password_wo_version" changed but "password_wo" is not available, so there is no password to send.`,
		)
		return
	}

	res := new(http.Response)
	requestOptions := []option.RequestOption{
		option.WithRequestBody("application/json", dataBytes),
	}
	// Add JSON fields after WithRequestBody, which replaces the full body.
	requestOptions = append(requestOptions, passwordRequestOptions(action, configPassword)...)
	requestOptions = append(requestOptions,
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)

	_, err = r.client.Storage.SftpStorages.UpdateAndPoll(
		ctx,
		data.ID.ValueInt64(),
		storage.SftpStorageUpdateParams{},
		requestOptions...,
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

	// The API may echo the password. Do not store it in state.
	data.PasswordWo = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageSftpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *StorageSftpModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res := new(http.Response)
	_, err := r.client.Storage.SftpStorages.Get(
		ctx,
		data.ID.ValueInt64(),
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

func (r *StorageSftpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *StorageSftpModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Storage.SftpStorages.DeleteAndPoll(
		ctx,
		data.ID.ValueInt64(),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageSftpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(StorageSftpModel)

	path := int64(0)
	diags := importpath.ParseImportID(
		req.ID,
		"<storage_id>",
		&path,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.Int64Value(path)

	res := new(http.Response)
	_, err := r.client.Storage.SftpStorages.Get(
		ctx,
		path,
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

func (r *StorageSftpResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	PasswordArgumentsRequiredTogether(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	PlanHasPassword(ctx, req, resp)
}
