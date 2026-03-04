// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_placement_group

import (
	"context"
	"errors"
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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudPlacementGroupResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudPlacementGroupResource)(nil)
var _ resource.ResourceWithImportState = (*CloudPlacementGroupResource)(nil)

func NewResource() resource.Resource {
	return &CloudPlacementGroupResource{}
}

// CloudPlacementGroupResource defines the resource implementation.
type CloudPlacementGroupResource struct {
	client *gcore.Client
}

func (r *CloudPlacementGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_placement_group"
}

func (r *CloudPlacementGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudPlacementGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudPlacementGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.PlacementGroupNewParams{}

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
	_, err = r.client.Cloud.PlacementGroups.New(
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
	data.ID = data.ServergroupID

	if !data.Instances.IsNullOrUnknown() {
		instances, _ := data.Instances.AsStructSliceT(ctx)
		for _, instance := range instances {
			addToPlacementGroupParams := cloud.InstanceAddToPlacementGroupParams{
				ServergroupID: data.ServergroupID.ValueString(),
			}
			if !data.ProjectID.IsNull() {
				addToPlacementGroupParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
			}
			if !data.RegionID.IsNull() {
				addToPlacementGroupParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
			}
			_, err := r.client.Cloud.Instances.AddToPlacementGroupAndPoll(
				ctx,
				instance.InstanceID.ValueString(),
				addToPlacementGroupParams,
				option.WithMiddleware(logging.Middleware(ctx)),
			)
			if err != nil {
				resp.Diagnostics.AddError("failed to make http request", err.Error())
				return
			}
		}
	}

	// After adding all instances, re-read the placement group to get the complete state
	getParams := cloud.PlacementGroupGetParams{}
	if !data.ProjectID.IsNull() {
		getParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}
	if !data.RegionID.IsNull() {
		getParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res = new(http.Response)
	_, err = r.client.Cloud.PlacementGroups.Get(
		ctx,
		data.ServergroupID.ValueString(),
		getParams,
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to read placement group after adding instances", err.Error())
		return
	}
	bytes, _ = io.ReadAll(res.Body)
	err = apijson.Unmarshal(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize placement group", err.Error())
		return
	}
	data.ID = data.ServergroupID

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudPlacementGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudPlacementGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudPlacementGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get instances from plan and state
	dataInstances, _ := data.Instances.AsStructSliceT(ctx)
	stateInstances, _ := state.Instances.AsStructSliceT(ctx)

	desiredInstanceIDs := make(map[string]struct{})
	currentInstanceIDs := make(map[string]struct{})

	for _, inst := range dataInstances {
		if !inst.InstanceID.IsNull() {
			desiredInstanceIDs[inst.InstanceID.ValueString()] = struct{}{}
		}
	}
	for _, inst := range stateInstances {
		if !inst.InstanceID.IsNull() {
			currentInstanceIDs[inst.InstanceID.ValueString()] = struct{}{}
		}
	}

	var instancesToAdd []string
	var instancesToRemove []string

	for id := range desiredInstanceIDs {
		if _, ok := currentInstanceIDs[id]; !ok {
			instancesToAdd = append(instancesToAdd, id)
		}
	}
	for id := range currentInstanceIDs {
		if _, ok := desiredInstanceIDs[id]; !ok {
			instancesToRemove = append(instancesToRemove, id)
		}
	}

	// Add new instances
	for _, instanceID := range instancesToAdd {
		addParams := cloud.InstanceAddToPlacementGroupParams{
			ServergroupID: state.ServergroupID.ValueString(),
		}
		if !data.ProjectID.IsNull() {
			addParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			addParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		_, err := r.client.Cloud.Instances.AddToPlacementGroupAndPoll(
			ctx,
			instanceID,
			addParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to add instance to placement group", err.Error())
			return
		}
	}

	// Remove instances
	for _, instanceID := range instancesToRemove {
		removeParams := cloud.InstanceRemoveFromPlacementGroupParams{}
		if !data.ProjectID.IsNull() {
			removeParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			removeParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		_, err := r.client.Cloud.Instances.RemoveFromPlacementGroupAndPoll(
			ctx,
			instanceID,
			removeParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			var apierr *gcore.Error
			if errors.As(err, &apierr) {
				if apierr.StatusCode == 404 {
					resp.Diagnostics.AddWarning("instance not found in a placement group.", err.Error())
				} else {
					resp.Diagnostics.AddError("failed to remove instance from placement group", err.Error())
					return
				}
			}
		}
	}

	// Re-read the placement group to get authoritative state
	getParams := cloud.PlacementGroupGetParams{}
	if !data.ProjectID.IsNull() {
		getParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}
	if !data.RegionID.IsNull() {
		getParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.PlacementGroups.Get(
		ctx,
		state.ServergroupID.ValueString(),
		getParams,
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to read placement group after update", err.Error())
		return
	}
	bytes, _ := io.ReadAll(res.Body)
	err = apijson.Unmarshal(bytes, &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize placement group", err.Error())
		return
	}
	data.ID = data.ServergroupID

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudPlacementGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudPlacementGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.PlacementGroupGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.PlacementGroups.Get(
		ctx,
		data.ServergroupID.ValueString(),
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
	data.ID = data.ServergroupID

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudPlacementGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudPlacementGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.PlacementGroupDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	_, err := r.client.Cloud.PlacementGroups.Delete(
		ctx,
		data.ServergroupID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	data.ID = data.ServergroupID

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudPlacementGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data *CloudPlacementGroupModel = new(CloudPlacementGroupModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_group_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<group_id>",
		&path_project_id,
		&path_region_id,
		&path_group_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.ServergroupID = types.StringValue(path_group_id)

	res := new(http.Response)
	_, err := r.client.Cloud.PlacementGroups.Get(
		ctx,
		path_group_id,
		cloud.PlacementGroupGetParams{
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
	data.ID = data.ServergroupID

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudPlacementGroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Don't modify plan during destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state *CloudPlacementGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If this is a create operation (no state), nothing to modify
	if req.State.Raw.IsNull() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If instances in config is null, preserve state value
	if plan.Instances.IsNull() && !state.Instances.IsNull() {
		plan.Instances = state.Instances
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
		return
	}

	// If both plan and state have instances, we need to handle the case where
	// new instances are being added with unknown instance_name values
	if !plan.Instances.IsNull() && !state.Instances.IsNull() {
		planInstances, diags := plan.Instances.AsStructSliceT(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		stateInstances, diags := state.Instances.AsStructSliceT(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Build a map of state instances by ID for quick lookup
		stateInstanceMap := make(map[string]CloudPlacementGroupInstancesModel)
		for _, inst := range stateInstances {
			if !inst.InstanceID.IsNull() {
				stateInstanceMap[inst.InstanceID.ValueString()] = inst
			}
		}

		// For each instance in the plan, if it exists in state, preserve the instance_name
		modified := false
		for i, planInst := range planInstances {
			if !planInst.InstanceID.IsNull() {
				if stateInst, exists := stateInstanceMap[planInst.InstanceID.ValueString()]; exists {
					// Instance exists in state, preserve the instance_name if plan has it as unknown
					if planInst.InstanceName.IsUnknown() && !stateInst.InstanceName.IsNull() {
						planInstances[i].InstanceName = stateInst.InstanceName
						modified = true
					}
				}
			}
		}

		// If we modified any instances, update the plan
		if modified {
			newInstances, diags := customfield.NewObjectSet(ctx, planInstances)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			plan.Instances = newInstances
			resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
		}
	}
}
