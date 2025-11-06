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

	// Handle interface attach/detach operations
	if !data.Interfaces.Equal(state.Interfaces) {
		// Get old and new interface lists
		oldInterfaces, diags := state.Interfaces.AsStructSliceT(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		newInterfaces, diags := data.Interfaces.AsStructSliceT(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Create a map of old interfaces for quick lookup
		oldInterfaceMap := make(map[string]bool)
		for _, oldIface := range oldInterfaces {
			if oldIface.SubnetID.ValueString() != "" {
				oldInterfaceMap[oldIface.SubnetID.ValueString()] = true
			}
		}

		// Process new interfaces: attach those that don't exist in old
		newInterfaceMap := make(map[string]bool)
		for _, newIface := range newInterfaces {
			subnetID := newIface.SubnetID.ValueString()
			if subnetID == "" {
				continue
			}
			newInterfaceMap[subnetID] = true

			// If this interface is not in old set, attach it
			if !oldInterfaceMap[subnetID] {
				attachParams := cloud.NetworkRouterAttachSubnetParams{
					SubnetID: subnetID,
				}
				if !data.ProjectID.IsNull() {
					attachParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
				}
				if !data.RegionID.IsNull() {
					attachParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
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
		}

		// Process old interfaces: detach those that are not in new set
		for _, oldIface := range oldInterfaces {
			subnetID := oldIface.SubnetID.ValueString()
			if subnetID == "" {
				continue
			}

			// If this interface is not in new set, detach it
			if !newInterfaceMap[subnetID] {
				detachParams := cloud.NetworkRouterDetachSubnetParams{
					SubnetID: cloud.SubnetIDParam{SubnetID: subnetID},
				}
				if !data.ProjectID.IsNull() {
					detachParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
				}
				if !data.RegionID.IsNull() {
					detachParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
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
		}
	}

	// Update other router attributes (name, routes, external_gateway_info)
	// Note: Route deletion is handled in ModifyPlan to update the plan before Update runs
	// Only send PATCH request if fields OTHER than interfaces have changed
	needsUpdate := !data.Name.Equal(state.Name) ||
		!data.Routes.Equal(state.Routes) ||
		!data.ExternalGatewayInfo.Equal(state.ExternalGatewayInfo)

	// Log route details for debugging
	var dataRoutes []CloudNetworkRouterRoutesModel
	var stateRoutes []CloudNetworkRouterRoutesModel
	data.Routes.ElementsAs(ctx, &dataRoutes, false)
	state.Routes.ElementsAs(ctx, &stateRoutes, false)

	// Add visible warning if routes are being removed
	if len(dataRoutes) == 0 && len(stateRoutes) > 0 {
		resp.Diagnostics.AddWarning(
			"Router route deletion detected",
			fmt.Sprintf("Removing %d routes from router. Plan has %d routes, State has %d routes. needsUpdate=%v",
				len(stateRoutes), len(dataRoutes), len(stateRoutes), needsUpdate),
		)
	}

	var err error
	if needsUpdate {
		params := cloud.NetworkRouterUpdateParams{}

		if !data.ProjectID.IsNull() {
			params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}

		if !data.RegionID.IsNull() {
			params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		// WORKAROUND: The SDK's NetworkRouterUpdateParams has `omitzero` on Routes field,
		// which causes empty routes array to be omitted from JSON. We need to explicitly
		// include routes=[] when deleting routes using option.WithJSONSet().
		var planRoutes, stateRoutes []CloudNetworkRouterRoutesModel
		data.Routes.ElementsAs(ctx, &planRoutes, false)
		state.Routes.ElementsAs(ctx, &stateRoutes, false)

		routesDeletionNeeded := len(planRoutes) == 0 && len(stateRoutes) > 0

		// Build update options
		updateOpts := []option.RequestOption{
			option.WithMiddleware(logging.Middleware(ctx)),
		}

		var dataBytes []byte
		dataBytes, err = data.MarshalJSONForUpdate(*state)
		if err != nil {
			resp.Diagnostics.AddError("failed to serialize http request", err.Error())
			return
		}

		// Skip PATCH if no actual changes (empty JSON body) and not deleting routes
		if (len(dataBytes) > 0 && string(dataBytes) != "{}" && string(dataBytes) != "null") || routesDeletionNeeded {
			res := new(http.Response)
			updateOpts = append(updateOpts,
				option.WithRequestBody("application/json", dataBytes),
				option.WithResponseBodyInto(&res),
			)

			// IMPORTANT: WithJSONSet must be added AFTER WithRequestBody so it can modify the body
			if routesDeletionNeeded {
				resp.Diagnostics.AddWarning(
					"Route deletion workaround",
					fmt.Sprintf("Forcing routes=[] in request body using WithJSONSet. Deleting %d routes.", len(stateRoutes)),
				)
				// Use the SDK's built-in mechanism to force routes field to be included
				updateOpts = append(updateOpts, option.WithJSONSet("routes", []interface{}{}))
			}

			_, err = r.client.Cloud.Networks.Routers.Update(
				ctx,
				routerID,
				params,
				updateOpts...,
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
	}

	// Do a final Read to get consistent state after attach/detach and update operations
	getParams := cloud.NetworkRouterGetParams{}
	if !data.ProjectID.IsNull() {
		getParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}
	if !data.RegionID.IsNull() {
		getParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

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
	var data *CloudNetworkRouterModel = new(CloudNetworkRouterModel)

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
	// Handle routes with computed_optional: when routes block is removed from config,
	// Terraform keeps the state value in plan (doesn't detect as changed). We need to
	// detect this and update the plan to have empty routes so deletion works properly.

	// Only apply during updates (not create or delete)
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

	// Check if routes were removed from config or set to empty array
	var configRoutes customfield.NestedObjectList[CloudNetworkRouterRoutesModel]
	diagsRoutes := req.Config.GetAttribute(ctx, path.Root("routes"), &configRoutes)

	// Get config routes as a list to check if empty
	var configRoutesList []CloudNetworkRouterRoutesModel
	configRoutesEmpty := false
	if !diagsRoutes.HasError() && !configRoutes.IsNull() {
		configRoutes.ElementsAs(ctx, &configRoutesList, false)
		configRoutesEmpty = len(configRoutesList) == 0
	}

	// Routes are considered "removed" if:
	// 1. Routes block is absent from config (configRoutes.IsNull()), OR
	// 2. Routes is explicitly set to empty array in config (configRoutesEmpty)
	routesRemovedFromConfig := !diagsRoutes.HasError() &&
		(configRoutes.IsNull() || configRoutesEmpty) &&
		!state.Routes.IsNull() &&
		!state.Routes.IsUnknown()

	if routesRemovedFromConfig {
		// Routes not in config (or set to empty) but exist in state - user removed routes.
		// Update plan to have empty routes so Terraform knows they'll be deleted.
		plan.Routes = customfield.NewObjectListMust(ctx, []CloudNetworkRouterRoutesModel{})
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}
}
