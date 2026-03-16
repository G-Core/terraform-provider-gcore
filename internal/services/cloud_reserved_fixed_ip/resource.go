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
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	var data *CloudReservedFixedIPModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudReservedFixedIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.ReservedFixedIPUpdateParams{}

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
	_, err = r.client.Cloud.ReservedFixedIPs.Update(
		ctx,
		state.PortID.ValueString(),
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to update reserved fixed IP", err.Error())
		return
	}

	// Update state with response
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.UnmarshalComputed(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	// Restore ID field (API doesn't return id, but we need it for Terraform)
	if !data.PortID.IsNull() && data.PortID.ValueString() != "" {
		data.ID = types.StringValue(data.PortID.ValueString())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
		data.PortID.ValueString(),
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

	// Restore ID field (API doesn't return id, but we need it for Terraform)
	if !data.PortID.IsNull() && data.PortID.ValueString() != "" {
		data.ID = types.StringValue(data.PortID.ValueString())
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
		data.PortID.ValueString(),
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
	var data = new(CloudReservedFixedIPModel)

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

	// Restore ID field (API doesn't return id, but we need it for Terraform)
	if !data.PortID.IsNull() && data.PortID.ValueString() != "" {
		data.ID = types.StringValue(data.PortID.ValueString())
	}

	// Infer the create-only Type field from API response fields, since the API does
	// not return it after creation. Best-effort coverage:
	//   external   → is_external == true
	//   subnet     → is_external == false && subnet_id is set
	//   any_subnet → is_external == false && network_id is set (no subnet_id)
	// ip_address and port types cannot be distinguished from the API response;
	// importing resources created with those types will produce a plan diff on the
	// type field and require manual state fixup.
	if data.Type.IsNull() || data.Type.ValueString() == "" {
		if !data.IsExternal.IsNull() && data.IsExternal.ValueBool() {
			data.Type = types.StringValue("external")
		} else if !data.SubnetID.IsNull() && data.SubnetID.ValueString() != "" {
			data.Type = types.StringValue("subnet")
		} else if !data.NetworkID.IsNull() && data.NetworkID.ValueString() != "" {
			data.Type = types.StringValue("any_subnet")
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudReservedFixedIPResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If port_id is not being explicitly configured and we have a state value, preserve it
	// This prevents unnecessary replacements when only is_vip changes
	var plan, state *CloudReservedFixedIPModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() || state == nil || plan == nil {
		return
	}

	// If port_id is unknown in the plan but known in state, preserve the state value
	// This happens when is_vip changes trigger a refresh of computed fields
	if plan.PortID.IsUnknown() && !state.PortID.IsNull() && !state.PortID.IsUnknown() {
		plan.PortID = state.PortID
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}
}
