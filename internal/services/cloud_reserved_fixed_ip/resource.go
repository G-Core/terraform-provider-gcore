// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_reserved_fixed_ip

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
	"github.com/stainless-sdks/gcore-terraform/internal/importpath"
	"github.com/stainless-sdks/gcore-terraform/internal/logging"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudReservedFixedIPResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudReservedFixedIPResource)(nil)
var _ resource.ResourceWithImportState = (*CloudReservedFixedIPResource)(nil)

func NewResource() resource.Resource {
	return &CloudReservedFixedIPResource{}
}

// CloudReservedFixedIPResource defines the resource implementation.
type CloudReservedFixedIPResource struct {
	client *gcore.Client
}

func (r *CloudReservedFixedIPResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_reserved_fixed_ip"
}

func (r *CloudReservedFixedIPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudReservedFixedIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudReservedFixedIPModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.ReservedFixedIPNewParams{}

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

	// Use NewAndPoll to automatically handle async operation and task polling
	reservedFixedIP, err := r.client.Cloud.ReservedFixedIPs.NewAndPoll(
		ctx,
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to create reserved fixed IP", err.Error())
		return
	}

	// Use RawJSON from the NewAndPoll response directly
	err = apijson.UnmarshalComputed([]byte(reservedFixedIP.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	// Set the ID to port_id for Terraform resource identification
	if !data.PortID.IsNull() && data.PortID.ValueString() != "" {
		data.ID = types.StringValue(data.PortID.ValueString())
	} else {
		resp.Diagnostics.AddError("Resource creation incomplete",
			"Could not extract port_id from API response")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudReservedFixedIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *CloudReservedFixedIPModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check for unsupported changes to allowed_address_pairs
	// Note: allowed_address_pairs is Computed-only in the schema, so users shouldn't be able to modify it directly.
	// But we add this check for completeness in case the schema changes in the future.
	if !plan.AllowedAddressPairs.Equal(state.AllowedAddressPairs) {
		resp.Diagnostics.AddError("Update Not Supported",
			"Updating 'allowed_address_pairs' is not supported yet. This feature requires Ports API integration which is not available in the current SDK version.")
		return
	}

	// Check for any other attempted changes to fields that should be immutable
	if !plan.Type.Equal(state.Type) {
		resp.Diagnostics.AddError("Update Not Supported",
			"The 'type' field cannot be changed after creation. Please recreate the resource with the new type.")
		return
	}

	if !plan.NetworkID.Equal(state.NetworkID) || !plan.SubnetID.Equal(state.SubnetID) || !plan.PortID.Equal(state.PortID) || !plan.IPAddress.Equal(state.IPAddress) {
		resp.Diagnostics.AddError("Update Not Supported",
			"Network configuration fields (network_id, subnet_id, port_id, ip_address) cannot be changed after creation. Please recreate the resource with the new configuration.")
		return
	}

	// Handle is_vip update via Vip.Toggle() API
	if !plan.IsVip.Equal(state.IsVip) {
		params := cloud.ReservedFixedIPVipToggleParams{
			IsVip: plan.IsVip.ValueBool(),
		}

		if !plan.ProjectID.IsNull() {
			params.ProjectID = param.NewOpt(plan.ProjectID.ValueInt64())
		}

		if !plan.RegionID.IsNull() {
			params.RegionID = param.NewOpt(plan.RegionID.ValueInt64())
		}

		res := new(http.Response)
		_, err := r.client.Cloud.ReservedFixedIPs.Vip.Toggle(
			ctx,
			state.PortID.ValueString(),
			params,
			option.WithResponseBodyInto(&res),
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to toggle VIP status", err.Error())
			return
		}

		// Update state with response
		bytes, _ := io.ReadAll(res.Body)
		err = apijson.UnmarshalComputed(bytes, &plan)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
			return
		}
	}

	// If we get here, update succeeded or no changes needed, save plan to state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CloudReservedFixedIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudReservedFixedIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If ID is null or empty, the resource doesn't exist yet - skip the read
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		return
	}

	params := cloud.ReservedFixedIPGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.ReservedFixedIPs.Get(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	// Check if it's a 404 error (resource not found)
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

func (r *CloudReservedFixedIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudReservedFixedIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.ReservedFixedIPDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	_, err := r.client.Cloud.ReservedFixedIPs.Delete(
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

func (r *CloudReservedFixedIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data *CloudReservedFixedIPModel = new(CloudReservedFixedIPModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_port_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<port_id>",
		&path_project_id,
		&path_region_id,
		&path_port_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.PortID = types.StringValue(path_port_id)

	res := new(http.Response)
	_, err := r.client.Cloud.ReservedFixedIPs.Get(
		ctx,
		path_port_id,
		cloud.ReservedFixedIPGetParams{
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
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudReservedFixedIPResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}
