// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_router

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
	"github.com/stainless-sdks/gcore-terraform/internal/importpath"
	"github.com/stainless-sdks/gcore-terraform/internal/logging"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudNetworkRouterResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudNetworkRouterResource)(nil)
var _ resource.ResourceWithImportState = (*CloudNetworkRouterResource)(nil)

func NewResource() resource.Resource {
	return &CloudNetworkRouterResource{}
}

// CloudNetworkRouterResource defines the resource implementation.
type CloudNetworkRouterResource struct {
	client *gcore.Client
}

func (r *CloudNetworkRouterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_network_router"
}

func (r *CloudNetworkRouterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudNetworkRouterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudNetworkRouterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.NetworkRouterNewParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	// According to API docs, interfaces and routes can be created in a single POST call
	dataBytes, err := data.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}
	router, err := r.client.Cloud.Networks.Routers.NewAndPoll(
		ctx,
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	// Use apijson.MarshalRoot instead of router.RawJSON() for consistent handling
	routerBytes, err := apijson.MarshalRoot(router)
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize router response", err.Error())
		return
	}
	err = apijson.UnmarshalComputed(routerBytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudNetworkRouterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudNetworkRouterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudNetworkRouterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	routerID := data.ID.ValueString()

	var projectIDOpt param.Opt[int64]
	if !data.ProjectID.IsNull() {
		projectIDOpt = param.NewOpt(data.ProjectID.ValueInt64())
	}

	var regionIDOpt param.Opt[int64]
	if !data.RegionID.IsNull() {
		regionIDOpt = param.NewOpt(data.RegionID.ValueInt64())
	}

	toAttach, toDetach, diags := diffInterfaces(ctx, data.Interfaces, state.Interfaces)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Step 1: Attach new interfaces FIRST (before PATCH)
	// This ensures interfaces exist before routes that reference them are added
	for _, subnetID := range toAttach {
		attachParams := cloud.NetworkRouterAttachSubnetParams{
			SubnetID:  subnetID,
			ProjectID: projectIDOpt,
			RegionID:  regionIDOpt,
		}

		_, err := r.client.Cloud.Networks.Routers.AttachSubnet(
			ctx,
			routerID,
			attachParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to attach subnet to router",
				fmt.Sprintf("SubnetID: %s, Error: %s", subnetID, err.Error()),
			)
			return
		}
	}

	// Step 2: Send PATCH for all field updates (routes, name, external_gateway_info)
	// This happens after attach and before detach to satisfy both:
	// - Routes can reference newly attached interfaces (nexthop validation)
	// - Routes are deleted before interfaces are detached (API constraint)
	params := cloud.NetworkRouterUpdateParams{}
	params.ProjectID = projectIDOpt
	params.RegionID = regionIDOpt

	dataBytes, err := data.MarshalJSONForUpdate(*state)
	if err != nil {
		resp.Diagnostics.AddError("failed to serialize http request", err.Error())
		return
	}

	if shouldSendPatch(dataBytes) {
		res := new(http.Response)
		_, err = r.client.Cloud.Networks.Routers.Update(
			ctx,
			routerID,
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
	}

	// Step 3: Detach old interfaces LAST (after PATCH)
	// This ensures routes referencing these interfaces are deleted first
	for _, subnetID := range toDetach {
		detachParams := cloud.NetworkRouterDetachSubnetParams{
			SubnetID:  cloud.SubnetIDParam{SubnetID: subnetID},
			ProjectID: projectIDOpt,
			RegionID:  regionIDOpt,
		}

		_, err := r.client.Cloud.Networks.Routers.DetachSubnet(
			ctx,
			routerID,
			detachParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to detach subnet from router",
				fmt.Sprintf("SubnetID: %s, Error: %s", subnetID, err.Error()),
			)
			return
		}
	}

	// Step 4: GET to refresh final state after all operations
	getParams := cloud.NetworkRouterGetParams{}
	getParams.ProjectID = projectIDOpt
	getParams.RegionID = regionIDOpt

	readRes := new(http.Response)
	_, err = r.client.Cloud.Networks.Routers.Get(
		ctx,
		routerID,
		getParams,
		option.WithResponseBodyInto(&readRes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to read router after update", err.Error())
		return
	}
	readBytes, _ := io.ReadAll(readRes.Body)
	err = apijson.UnmarshalComputed(readBytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize router after update", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudNetworkRouterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudNetworkRouterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.NetworkRouterGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.Networks.Routers.Get(
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

func (r *CloudNetworkRouterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudNetworkRouterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.NetworkRouterDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.Networks.Routers.DeleteAndPoll(
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

func (r *CloudNetworkRouterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudNetworkRouterModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_router_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<router_id>",
		&path_project_id,
		&path_region_id,
		&path_router_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.ID = types.StringValue(path_router_id)

	res := new(http.Response)
	_, err := r.client.Cloud.Networks.Routers.Get(
		ctx,
		path_router_id,
		cloud.NetworkRouterGetParams{
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

func (r *CloudNetworkRouterResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan *CloudNetworkRouterModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudNetworkRouterModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var configRoutes customfield.NestedObjectList[CloudNetworkRouterRoutesModel]
	diagsRoutes := req.Config.GetAttribute(ctx, path.Root("routes"), &configRoutes)

	var configRoutesList []CloudNetworkRouterRoutesModel
	configRoutesEmpty := false
	if !diagsRoutes.HasError() && !configRoutes.IsNull() {
		configRoutes.ElementsAs(ctx, &configRoutesList, false)
		configRoutesEmpty = len(configRoutesList) == 0
	}

	routesRemovedFromConfig := !diagsRoutes.HasError() &&
		(configRoutes.IsNull() || configRoutesEmpty) &&
		!state.Routes.IsNull() &&
		!state.Routes.IsUnknown()

	if routesRemovedFromConfig {
		plan.Routes = customfield.NewObjectListMust(ctx, []CloudNetworkRouterRoutesModel{})
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}
}

func diffInterfaces(ctx context.Context, plan customfield.NestedObjectList[CloudNetworkRouterInterfacesModel], state customfield.NestedObjectList[CloudNetworkRouterInterfacesModel]) (toAttach []string, toDetach []string, diags diag.Diagnostics) {
	if plan.Equal(state) {
		return
	}

	oldInterfaces, diagnostics := state.AsStructSliceT(ctx)
	diags.Append(diagnostics...)
	if diags.HasError() {
		return
	}

	newInterfaces, diagnostics := plan.AsStructSliceT(ctx)
	diags.Append(diagnostics...)
	if diags.HasError() {
		return
	}

	oldInterfaceMap := make(map[string]struct{}, len(oldInterfaces))
	for _, oldIface := range oldInterfaces {
		if subnetID := oldIface.SubnetID.ValueString(); subnetID != "" {
			oldInterfaceMap[subnetID] = struct{}{}
		}
	}

	newInterfaceMap := make(map[string]struct{}, len(newInterfaces))
	for _, newIface := range newInterfaces {
		if subnetID := newIface.SubnetID.ValueString(); subnetID != "" {
			newInterfaceMap[subnetID] = struct{}{}
			if _, exists := oldInterfaceMap[subnetID]; !exists {
				toAttach = append(toAttach, subnetID)
			}
		}
	}

	for _, oldIface := range oldInterfaces {
		if subnetID := oldIface.SubnetID.ValueString(); subnetID != "" {
			if _, exists := newInterfaceMap[subnetID]; !exists {
				toDetach = append(toDetach, subnetID)
			}
		}
	}

	return
}

func shouldSendPatch(dataBytes []byte) bool {
	return len(dataBytes) > 0 && string(dataBytes) != "{}" && string(dataBytes) != "null"
}
