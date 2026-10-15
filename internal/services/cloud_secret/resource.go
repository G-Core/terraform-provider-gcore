// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_secret

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
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudSecretResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudSecretResource)(nil)
var _ resource.ResourceWithImportState = (*CloudSecretResource)(nil)

func NewResource() resource.Resource {
	return &CloudSecretResource{}
}

// CloudSecretResource defines the resource implementation.
type CloudSecretResource struct {
	client *gcore.Client
}

func (r *CloudSecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_secret"
}

func (r *CloudSecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudSecretModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Write-only fields are null in plan; read from config
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("payload"), &data.Payload)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.SecretUploadTlsCertificateParams{}

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
	_, err = r.client.Cloud.Secrets.UploadTlsCertificateAndPoll(
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update sends nothing to the API: there is no PATCH endpoint, so every
	// plan that reaches it is written straight to state. Plans that reach it
	// include the import-adoption apply, the in-place removal of
	// payload/payload_wo_version, and removing project_id/region_id from
	// config (int64planmodifier.RequiresReplaceIfConfigured returns without
	// requiring replacement when the config value is null - see schema.go).
	// settle re-arms the marker on an errored apply and clears it on a
	// successful one (see import_adoption.go).
	_, settle := secretImportAdoption.BeginUpdate(ctx, req.Private, resp.Private, &resp.Diagnostics)
	defer settle()
	if resp.Diagnostics.HasError() {
		return
	}

	var data *CloudSecretModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudSecretModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// None of these are settable in config, and no API call happens here to
	// refresh them: without a plan modifier of their own they plan as unknown
	// on any in-place update, which this resource never had before the
	// adoption window existed. Carry the prior state forward.
	data.Algorithm = state.Algorithm
	data.BitLength = state.BitLength
	data.Created = state.Created
	data.Mode = state.Mode
	data.SecretType = state.SecretType
	data.Status = state.Status
	data.ContentTypes = state.ContentTypes

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudSecretModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.SecretGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	priorExpiration := data.Expiration

	res := new(http.Response)
	_, err := r.client.Cloud.Secrets.Get(
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
	// The API echoes expiration in UTC ("...Z"); keep the configured spelling
	// of the same instant (for example "+03:00"), or every plan shows a diff.
	if sameInstant(priorExpiration.StringValue, data.Expiration.StringValue) {
		data.Expiration = priorExpiration
	}

	// Advance the import-adoption marker (see import_adoption.go). Must run on
	// every successful Read: Terraform persists this refresh in `terraform
	// apply` even when the resulting plan is a no-op, which is the only hook
	// that retires a marker whose import adopted nothing.
	resp.Diagnostics.Append(secretImportAdoption.Advance(ctx, resp.Private)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudSecretModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.SecretDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.Secrets.DeleteAndPoll(
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

func (r *CloudSecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudSecretModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_secret_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<secret_id>",
		&path_project_id,
		&path_region_id,
		&path_secret_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.ID = types.StringValue(path_secret_id)

	res := new(http.Response)
	_, err := r.client.Cloud.Secrets.Get(
		ctx,
		path_secret_id,
		cloud.SecretGetParams{
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

	// Open the import-adoption window: the API never returns payload_wo_version,
	// so it is null here. See import_adoption.go.
	resp.Diagnostics.Append(secretImportAdoption.Mark(ctx, resp.Private)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudSecretResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	requirePayloadOnCreate(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	secretImportAdoptionPlan.ModifyPlan(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	// See plan_modifiers.go: without this, the seven read-only attributes
	// below show "(known after apply)" on any Update-shaped plan, not only
	// the payload-removal one, because none of them has a plan modifier of
	// its own.
	carryComputedAttributesForward(ctx, req, resp)
}
