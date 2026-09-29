// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_listener

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudLoadBalancerListenerResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudLoadBalancerListenerResource)(nil)
var _ resource.ResourceWithImportState = (*CloudLoadBalancerListenerResource)(nil)

func NewResource() resource.Resource {
	return &CloudLoadBalancerListenerResource{}
}

// CloudLoadBalancerListenerResource defines the resource implementation.
type CloudLoadBalancerListenerResource struct {
	client *gcore.Client
}

func (r *CloudLoadBalancerListenerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_load_balancer_listener"
}

func (r *CloudLoadBalancerListenerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudLoadBalancerListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudLoadBalancerListenerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dataBytes, err := data.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}

	params := cloud.LoadBalancerListenerNewParams{}
	if err := params.UnmarshalJSON(dataBytes); err != nil {
		resp.Diagnostics.AddError("failed to deserialize into params", err.Error())
		return
	}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err = r.client.Cloud.LoadBalancers.Listeners.NewAndPoll(
		ctx,
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudLoadBalancerListenerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudLoadBalancerListenerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read the adoption marker before any of the branches below: every one of
	// them falls through to a state write, which is where the window closes.
	adopting, adoptionDiags := importAdoptionPending(ctx, req.Private)
	resp.Diagnostics.Append(adoptionDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// A failed adopting apply must not consume the adoption window: the
	// refresh that preceded this Update already advanced the marker, and
	// without a re-arm the next refresh would retire it and the retry would
	// plan a destroy+recreate. Terraform persists private state returned from
	// an errored apply, so this write reaches the state file.
	defer func() {
		if adopting && resp.Diagnostics.HasError() {
			resp.Diagnostics.Append(rearmImportAdoption(ctx, resp.Private)...)
		}
	}()

	// Keeps adopted create-only values out of the PATCH body; see patchBaselineFor.
	dataBytes, err := data.MarshalJSONForUpdate(patchBaselineFor(*state, *data, adopting))
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}

	dataStr := strings.TrimSpace(string(dataBytes))

	// If no fields have changed, skip the update and just refresh from API
	if dataStr == "{}" || dataStr == "null" || len(dataBytes) == 0 {
		// No changes to send - just read current state
		getParams := cloud.LoadBalancerListenerGetParams{}
		if !data.ProjectID.IsNull() {
			getParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			getParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		res := new(http.Response)
		_, err := r.client.Cloud.LoadBalancers.Listeners.Get(
			ctx,
			data.ID.ValueString(),
			getParams,
			option.WithResponseBodyInto(&res),
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to read listener", err.Error())
			return
		}
		bytes, _ := io.ReadAll(res.Body)
		err = apijson.UnmarshalComputed(bytes, &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize response", err.Error())
			return
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

		// Clear only after the write succeeded: a marker cleared before
		// success would turn a transient failure into a destroy+recreate on
		// the retry, because the next plan would no longer recognise the
		// import.
		if adopting && !resp.Diagnostics.HasError() {
			resp.Diagnostics.Append(clearImportAdoptionPending(ctx, resp.Private)...)
		}
		return
	}

	params := cloud.LoadBalancerListenerUpdateParams{}
	if err := params.UnmarshalJSON(dataBytes); err != nil {
		resp.Diagnostics.AddError("failed to deserialize into params", err.Error())
		return
	}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	resUpdate := new(http.Response)
	_, err = r.client.Cloud.LoadBalancers.Listeners.UpdateAndPoll(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithResponseBodyInto(&resUpdate),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	bytesUpdate, _ := io.ReadAll(resUpdate.Body)
	err = apijson.UnmarshalComputed(bytesUpdate, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Clear only after the write succeeded: a marker cleared before success
	// would turn a transient failure into a destroy+recreate on the retry,
	// because the next plan would no longer recognise the import.
	if adopting && !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(clearImportAdoptionPending(ctx, resp.Private)...)
	}
}

func (r *CloudLoadBalancerListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudLoadBalancerListenerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.LoadBalancerListenerGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.LoadBalancers.Listeners.Get(
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
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	// Advance the import-adoption marker (see import_adoption.go). Must run on
	// every successful Read: Terraform persists this refresh in `terraform
	// apply` even when the resulting plan is a no-op, which is the only hook
	// that retires a marker whose import adopted nothing.
	resp.Diagnostics.Append(advanceImportAdoption(ctx, resp.Private)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudLoadBalancerListenerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.LoadBalancerListenerDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.LoadBalancers.Listeners.DeleteAndPoll(
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

func (r *CloudLoadBalancerListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudLoadBalancerListenerModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_listener_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<listener_id>",
		&path_project_id,
		&path_region_id,
		&path_listener_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.ID = types.StringValue(path_listener_id)

	res := new(http.Response)
	_, err := r.client.Cloud.LoadBalancers.Listeners.Get(
		ctx,
		path_listener_id,
		cloud.LoadBalancerListenerGetParams{
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

	// default_pool_id and insert_x_forwarded are create-only and tagged
	// no_refresh, so the Get above cannot populate them and they stay null in
	// state. Mark the adoption window so the plan modifiers on those
	// attributes know this null came from an import rather than from a
	// listener created without them - the first apply adopts the configured
	// values into state and clears the marker.
	resp.Diagnostics.Append(markImportAdoptionPending(ctx, resp.Private)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerListenerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// The adoption marker is deliberately NOT touched here. A plan whose only
	// effect would be a marker write is a no-op plan, and Terraform never
	// persists PlannedPrivate for no-op plans, so the write would silently
	// evaporate. Retirement lives in Read's phase advance instead (see
	// advanceImportAdoption in import_adoption.go), whose private-state writes
	// every apply persists; Update handles the success (clear) and failure
	// (re-arm) cases.
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

	if resp.Diagnostics.HasError() || !adopting {
		return
	}

	var plan, state *CloudLoadBalancerListenerModel

	// The warning is advisory, so a decode problem must never fail the plan.
	if req.Plan.Get(ctx, &plan).HasError() || req.State.Get(ctx, &state).HasError() {
		return
	}

	if plan == nil || state == nil {
		return
	}

	// A replacement runs Create, not Update: nothing is adopted and every
	// configured value reaches the API, so the warning would describe the
	// opposite of what is about to happen.
	if planForcesReplacement(ctx, req.Plan, req.State, req.Config, adopting) {
		return
	}

	if adopted := adoptedCreateOnlyAttributes(*plan, *state); len(adopted) > 0 {
		resp.Diagnostics.Append(importAdoptionWarning(adopted))
	}
}
