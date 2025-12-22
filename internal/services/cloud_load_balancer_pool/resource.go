// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool

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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/importpath"
	"github.com/stainless-sdks/gcore-terraform/internal/logging"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudLoadBalancerPoolResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudLoadBalancerPoolResource)(nil)
var _ resource.ResourceWithImportState = (*CloudLoadBalancerPoolResource)(nil)

func NewResource() resource.Resource {
	return &CloudLoadBalancerPoolResource{}
}

// CloudLoadBalancerPoolResource defines the resource implementation.
type CloudLoadBalancerPoolResource struct {
	client *gcore.Client
}

func (r *CloudLoadBalancerPoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_load_balancer_pool"
}

func (r *CloudLoadBalancerPoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudLoadBalancerPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudLoadBalancerPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dataBytes, err := data.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}

	params := cloud.LoadBalancerPoolNewParams{}
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

	pool, err := r.client.Cloud.LoadBalancers.Pools.NewAndPoll(
		ctx,
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	// Use raw JSON from the response to unmarshal the "computed" fields into the data model
	err = apijson.UnmarshalComputed([]byte(pool.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudLoadBalancerPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudLoadBalancerPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Handle health monitor deletion via dedicated DELETE endpoint
	if state.Healthmonitor != nil && data.Healthmonitor == nil {
		deleteParams := cloud.LoadBalancerPoolHealthMonitorDeleteParams{}

		if !data.ProjectID.IsNull() {
			deleteParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}

		if !data.RegionID.IsNull() {
			deleteParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		err := r.client.Cloud.LoadBalancers.Pools.HealthMonitors.Delete(
			ctx,
			data.ID.ValueString(),
			deleteParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to delete health monitor", err.Error())
			return
		}
		// After successful DELETE, set state.Healthmonitor to nil so MarshalJSONForUpdate
		// doesn't include healthmonitor:null in the patch (deletion already handled above)
		state.Healthmonitor = nil
	}

	dataBytes, err := data.MarshalJSONForUpdate(*state)
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}

	dataStr := strings.TrimSpace(string(dataBytes))

	// If no fields have changed, skip the update and just refresh from API
	if dataStr == "{}" || dataStr == "null" || len(dataBytes) == 0 {
		// No changes to send - just read current state
		getParams := cloud.LoadBalancerPoolGetParams{}
		if !data.ProjectID.IsNull() {
			getParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			getParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		res := new(http.Response)
		_, err := r.client.Cloud.LoadBalancers.Pools.Get(
			ctx,
			data.ID.ValueString(),
			getParams,
			option.WithResponseBodyInto(&res),
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to read pool", err.Error())
			return
		}
		bytes, _ := io.ReadAll(res.Body)
		err = apijson.UnmarshalComputed(bytes, &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize response", err.Error())
			return
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	params := cloud.LoadBalancerPoolUpdateParams{}
	if err := params.UnmarshalJSON(dataBytes); err != nil {
		resp.Diagnostics.AddError("failed to deserialize into params", err.Error())
		return
	}

	// SDK's LoadBalancerPoolUpdateParamsHealthmonitor has required fields (Delay, MaxRetries, Timeout)
	// that are plain int64 (not param.Opt), so partial patch JSON leaves them at zero values.
	// If healthmonitor is being updated, we must ensure all required fields are populated.
	if data.Healthmonitor != nil && state.Healthmonitor != nil {
		// Check if any healthmonitor field was included in the patch (non-zero in params)
		hmChanged := params.Healthmonitor.Delay != 0 ||
			params.Healthmonitor.MaxRetries != 0 ||
			params.Healthmonitor.Timeout != 0 ||
			params.Healthmonitor.Type != "" ||
			params.Healthmonitor.HTTPMethod != "" ||
			params.Healthmonitor.URLPath.Valid() ||
			params.Healthmonitor.ExpectedCodes.Valid() ||
			params.Healthmonitor.MaxRetriesDown.Valid()

		if hmChanged {
			// Fill in required fields from plan if they're still zero
			if params.Healthmonitor.Delay == 0 && !data.Healthmonitor.Delay.IsNull() {
				params.Healthmonitor.Delay = data.Healthmonitor.Delay.ValueInt64()
			}
			if params.Healthmonitor.MaxRetries == 0 && !data.Healthmonitor.MaxRetries.IsNull() {
				params.Healthmonitor.MaxRetries = data.Healthmonitor.MaxRetries.ValueInt64()
			}
			if params.Healthmonitor.Timeout == 0 && !data.Healthmonitor.Timeout.IsNull() {
				params.Healthmonitor.Timeout = data.Healthmonitor.Timeout.ValueInt64()
			}
			if params.Healthmonitor.Type == "" && !data.Healthmonitor.Type.IsNull() {
				params.Healthmonitor.Type = cloud.LbHealthMonitorType(data.Healthmonitor.Type.ValueString())
			}
		}
	}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	pool, err := r.client.Cloud.LoadBalancers.Pools.UpdateAndPoll(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	// Use raw JSON from the response to unmarshal the "computed" fields into the data model
	err = apijson.UnmarshalComputed([]byte(pool.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudLoadBalancerPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.LoadBalancerPoolGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.LoadBalancers.Pools.Get(
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudLoadBalancerPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.LoadBalancerPoolDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.LoadBalancers.Pools.DeleteAndPoll(
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

func (r *CloudLoadBalancerPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudLoadBalancerPoolModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_pool_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<pool_id>",
		&path_project_id,
		&path_region_id,
		&path_pool_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.ID = types.StringValue(path_pool_id)

	res := new(http.Response)
	_, err := r.client.Cloud.LoadBalancers.Pools.Get(
		ctx,
		path_pool_id,
		cloud.LoadBalancerPoolGetParams{
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

	// Populate listener_id from listeners relationship for proper import.
	// The API doesn't return listener_id directly, but returns it in the listeners array.
	// Without this, terraform plan after import would show listener_id changing from null
	// to the config value, which would force replacement due to RequiresReplace() modifier.
	// See: GCLOUD2-20778
	if data.ListenerID.IsNull() || data.ListenerID.IsUnknown() {
		if !data.Listeners.IsNull() && !data.Listeners.IsUnknown() {
			listeners, diags := data.Listeners.AsStructSliceT(ctx)
			if !diags.HasError() && len(listeners) > 0 && !listeners[0].ID.IsNull() {
				data.ListenerID = listeners[0].ID
			}
		}
	}
	// Note: We intentionally do NOT populate load_balancer_id from loadbalancers relationship.
	// Users typically create pools with either listener_id OR load_balancer_id, not both.
	// If we populate load_balancer_id when user's config only has listener_id, it would
	// cause drift (state has load_balancer_id, config doesn't, triggers replacement).

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerPoolResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}
