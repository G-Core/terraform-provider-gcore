// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_cluster

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
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/G-Core/terraform-provider-gcore/internal/importpath"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudK8SClusterResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudK8SClusterResource)(nil)
var _ resource.ResourceWithImportState = (*CloudK8SClusterResource)(nil)

func NewResource() resource.Resource {
	return &CloudK8SClusterResource{}
}

// CloudK8SClusterResource defines the resource implementation.
type CloudK8SClusterResource struct {
	client *gcore.Client
}

func (r *CloudK8SClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_k8s_cluster"
}

func (r *CloudK8SClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudK8SClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudK8SClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.K8SClusterNewParams{}

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
	// explicitly set params.Name which is needed to perform the get request after the polling is done
	params.Name = data.Name.ValueString()
	cluster, err := r.client.Cloud.K8S.Clusters.NewAndPoll(
		ctx,
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	err = apijson.UnmarshalComputed([]byte(cluster.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}
	data.ID = data.Name
	data.FilterServerManagedLabels(ctx)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudK8SClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudK8SClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *CloudK8SClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save the planned pool order - we'll need to restore this after refreshing from API
	// because the API may return pools in a different order
	plannedPoolOrder := getPoolOrder(data.Pools)

	var stateHasChanged bool

	clusterName := data.Name.ValueString()

	// 1. Check for cluster upgrade.
	// We check if the cluster version is being updated, if we so we need upgrade the cluster before executing the
	// normal update flow
	if !data.Version.IsNull() && data.Version.ValueString() != state.Version.ValueString() {
		upgradeParams := cloud.K8SClusterUpgradeParams{}
		if !data.ProjectID.IsNull() {
			upgradeParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			upgradeParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		upgradeParams.Version = data.Version.ValueString()
		cluster, err := r.client.Cloud.K8S.Clusters.UpgradeAndPoll(
			ctx,
			data.Name.ValueString(),
			upgradeParams,
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to upgrade k8s cluster", err.Error())
			return
		}
		err = apijson.UnmarshalComputed([]byte(cluster.RawJSON()), &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
			return
		}
		data.ID = data.Name
		data.FilterServerManagedLabels(ctx)
		stateHasChanged = true
	}

	// 2. Check for changes in cluster updatable fields.
	// Check for changes in attributes or blocks that can be updated in-place. Pools are handled in a separate block below.
	// TODO: deal with changes to ddos_profile and logging
	if !data.Authentication.Equal(state.Authentication) ||
		!data.AutoscalerConfig.Equal(state.AutoscalerConfig) ||
		!data.Cni.Equal(state.Cni) ||
		// Check if value changed from null to set, or between set values
		//(!data.DDOSProfile.IsNull() && state.DDOSProfile.IsNull()) ||
		//(!data.DDOSProfile.IsNull() && !state.DDOSProfile.IsNull() && !data.DDOSProfile.Equal(state.DDOSProfile)) ||
		//!data.Logging.Equal(state.Logging) ||
		!data.AddOns.Equal(state.AddOns) {
		params := cloud.K8SClusterUpdateParams{}

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
		cluster, err := r.client.Cloud.K8S.Clusters.UpdateAndPoll(
			ctx,
			data.Name.ValueString(),
			params,
			option.WithRequestBody("application/json", dataBytes),
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to make http request", err.Error())
			return
		}
		err = apijson.UnmarshalComputed([]byte(cluster.RawJSON()), &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
			return
		}
		data.ID = data.Name
		stateHasChanged = true
	}

	// 3. Handle pool changes
	// 1 pool   => Allow in-place updates and add/delete, but return error on replace.
	//             Users must create a new pool with different name in such case.
	// 2+ pools => Allow all operations, but make sure we don't end up with 0 pools at any moment.
	//             This means we process each pool change one-by-one, and perform adds before deletes.
	if !poolsEqual(data.Pools, state.Pools) {
		var projectID, regionID param.Opt[int64]
		if !data.ProjectID.IsNull() {
			projectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			regionID = param.NewOpt(data.RegionID.ValueInt64())
		}

		oldPools := poolsToMap(state.Pools)
		newPools := poolsToMap(data.Pools)

		// Check quota limits before making any changes.
		// This avoids deleting pools and then failing to create new ones due to quota limits.
		if err := checkPoolQuotaLimits(ctx, r.client, projectID, regionID, oldPools, newPools); err != nil {
			resp.Diagnostics.AddError("quota limits exceeded", err.Error())
			return
		}

		// 1. Create new pools first (before any deletes/replaces, to avoid having 0 pools)
		for poolName, newPool := range newPools {
			if _, exists := oldPools[poolName]; !exists {
				if err := createPool(ctx, r.client, clusterName, projectID, regionID, newPool); err != nil {
					resp.Diagnostics.AddError("failed to create pool", err.Error())
					return
				}
				stateHasChanged = true
			}
		}

		// 2. Handle pool replacements and updates
		for poolName, newPool := range newPools {
			oldPool, exists := oldPools[poolName]
			if !exists {
				continue // already handled above
			}
			if poolNeedsReplace(oldPool, newPool) {
				// Check if this is the only pool - can't replace the only pool
				if len(oldPools) == 1 && len(newPools) == 1 {
					resp.Diagnostics.AddError(
						"cannot replace the only pool",
						"Cannot replace the only pool in the cluster. Please create a new pool with a different name first, then remove this one.",
					)
					return
				}
				// Delete old pool
				if err := deletePool(ctx, r.client, clusterName, projectID, regionID, poolName); err != nil {
					resp.Diagnostics.AddError("failed to delete pool for replacement", err.Error())
					return
				}
				// Create new pool with same name
				if err := createPool(ctx, r.client, clusterName, projectID, regionID, newPool); err != nil {
					resp.Diagnostics.AddError("failed to create replacement pool", err.Error())
					return
				}
				stateHasChanged = true
			} else if poolNeedsUpdate(oldPool, newPool) {
				if err := updatePool(ctx, r.client, clusterName, projectID, regionID, poolName, newPool); err != nil {
					resp.Diagnostics.AddError("failed to update pool", err.Error())
					return
				}
				stateHasChanged = true
			}
		}

		// 3. Delete removed pools (after all additions to avoid 0 pools)
		// Check if we would delete all pools (sanity check - should not happen due to schema validation)
		if len(newPools) == 0 && len(oldPools) > 0 {
			resp.Diagnostics.AddError(
				"cannot delete all pools",
				"A cluster must have at least one pool. Please ensure at least one pool is configured.",
			)
			return
		}
		for poolName := range oldPools {
			if _, exists := newPools[poolName]; !exists {
				if err := deletePool(ctx, r.client, clusterName, projectID, regionID, poolName); err != nil {
					resp.Diagnostics.AddError("failed to delete pool", err.Error())
					return
				}
				stateHasChanged = true
			}
		}
	}

	// 4. Refresh state from API after pool changes
	// The overhead of this is necessary as pool operations can change cluster state
	if stateHasChanged {
		getParams := cloud.K8SClusterGetParams{}
		if !data.ProjectID.IsNull() {
			getParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
		}
		if !data.RegionID.IsNull() {
			getParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
		}
		res := new(http.Response)
		_, err := r.client.Cloud.K8S.Clusters.Get(
			ctx,
			clusterName,
			getParams,
			option.WithResponseBodyInto(&res),
			option.WithMiddleware(logging.Middleware(ctx)),
		)
		if err != nil {
			resp.Diagnostics.AddError("failed to refresh cluster state after pool changes", err.Error())
			return
		}
		bytes, _ := io.ReadAll(res.Body)
		err = apijson.Unmarshal(bytes, &data)
		if err != nil {
			resp.Diagnostics.AddError("failed to deserialize cluster state", err.Error())
			return
		}
		data.ID = data.Name
		data.FilterServerManagedLabels(ctx)

		// Reorder pools to match the planned order
		// The API may return pools in arbitrary order, but Terraform expects them
		// to match the plan to avoid "inconsistent result after apply" errors
		data.Pools = reorderPoolsToMatch(data.Pools, plannedPoolOrder)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudK8SClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudK8SClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.K8SClusterGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	res := new(http.Response)
	_, err := r.client.Cloud.K8S.Clusters.Get(
		ctx,
		data.Name.ValueString(),
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
	data.ID = data.Name
	data.FilterServerManagedLabels(ctx)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudK8SClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudK8SClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.K8SClusterDeleteParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.K8S.Clusters.DeleteAndPoll(
		ctx,
		data.Name.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	data.ID = data.Name

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudK8SClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudK8SClusterModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_cluster_name := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<cluster_name>",
		&path_project_id,
		&path_region_id,
		&path_cluster_name,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.Name = types.StringValue(path_cluster_name)

	res := new(http.Response)
	_, err := r.client.Cloud.K8S.Clusters.Get(
		ctx,
		path_cluster_name,
		cloud.K8SClusterGetParams{
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
	data.ID = data.Name
	data.FilterServerManagedLabels(ctx)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudK8SClusterResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}

// createPool creates a new pool in the cluster and waits for completion
func createPool(ctx context.Context, client *gcore.Client, clusterName string, projectID, regionID param.Opt[int64], pool *CloudK8SClusterPoolsModel) error {
	params := cloud.K8SClusterPoolNewParams{
		ProjectID:    projectID,
		RegionID:     regionID,
		Name:         pool.Name.ValueString(),
		FlavorID:     pool.FlavorID.ValueString(),
		MinNodeCount: pool.MinNodeCount.ValueInt64(),
	}

	if !pool.MaxNodeCount.IsNull() && !pool.MaxNodeCount.IsUnknown() {
		params.MaxNodeCount = param.NewOpt(pool.MaxNodeCount.ValueInt64())
	}
	if !pool.AutoHealingEnabled.IsNull() && !pool.AutoHealingEnabled.IsUnknown() {
		params.AutoHealingEnabled = param.NewOpt(pool.AutoHealingEnabled.ValueBool())
	}
	if !pool.BootVolumeSize.IsNull() && !pool.BootVolumeSize.IsUnknown() {
		params.BootVolumeSize = param.NewOpt(pool.BootVolumeSize.ValueInt64())
	}
	if !pool.BootVolumeType.IsNull() && !pool.BootVolumeType.IsUnknown() {
		params.BootVolumeType = cloud.K8SClusterPoolNewParamsBootVolumeType(pool.BootVolumeType.ValueString())
	}
	if !pool.IsPublicIpv4.IsNull() && !pool.IsPublicIpv4.IsUnknown() {
		params.IsPublicIpv4 = param.NewOpt(pool.IsPublicIpv4.ValueBool())
	}
	if !pool.ServergroupPolicy.IsNull() && !pool.ServergroupPolicy.IsUnknown() {
		params.ServergroupPolicy = cloud.K8SClusterPoolNewParamsServergroupPolicy(pool.ServergroupPolicy.ValueString())
	}

	// Handle map fields
	if !pool.Labels.IsNull() && !pool.Labels.IsUnknown() {
		labelsMap, _ := pool.Labels.Value(ctx)
		params.Labels = make(map[string]string)
		for k, v := range labelsMap {
			params.Labels[k] = v.ValueString()
		}
	}
	if !pool.Taints.IsNull() && !pool.Taints.IsUnknown() {
		taintsMap, _ := pool.Taints.Value(ctx)
		params.Taints = make(map[string]string)
		for k, v := range taintsMap {
			params.Taints[k] = v.ValueString()
		}
	}
	if !pool.CrioConfig.IsNull() && !pool.CrioConfig.IsUnknown() {
		crioMap, _ := pool.CrioConfig.Value(ctx)
		params.CrioConfig = make(map[string]string)
		for k, v := range crioMap {
			params.CrioConfig[k] = v.ValueString()
		}
	}
	if !pool.KubeletConfig.IsNull() && !pool.KubeletConfig.IsUnknown() {
		kubeletMap, _ := pool.KubeletConfig.Value(ctx)
		params.KubeletConfig = make(map[string]string)
		for k, v := range kubeletMap {
			params.KubeletConfig[k] = v.ValueString()
		}
	}

	_, err := client.Cloud.K8S.Clusters.Pools.NewAndPoll(
		ctx,
		clusterName,
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		return fmt.Errorf("create pool %s: %w", pool.Name.ValueString(), err)
	}

	return nil
}

// updatePool updates an existing pool in the cluster
func updatePool(ctx context.Context, client *gcore.Client, clusterName string, projectID, regionID param.Opt[int64], poolName string, pool *CloudK8SClusterPoolsModel) error {
	params := cloud.K8SClusterPoolUpdateParams{
		ProjectID:   projectID,
		RegionID:    regionID,
		ClusterName: clusterName,
	}

	if !pool.MinNodeCount.IsNull() && !pool.MinNodeCount.IsUnknown() {
		params.MinNodeCount = param.NewOpt(pool.MinNodeCount.ValueInt64())
	}
	if !pool.MaxNodeCount.IsNull() && !pool.MaxNodeCount.IsUnknown() {
		params.MaxNodeCount = param.NewOpt(pool.MaxNodeCount.ValueInt64())
	}
	if !pool.AutoHealingEnabled.IsNull() && !pool.AutoHealingEnabled.IsUnknown() {
		params.AutoHealingEnabled = param.NewOpt(pool.AutoHealingEnabled.ValueBool())
	}

	// Handle map fields
	if !pool.Labels.IsNull() && !pool.Labels.IsUnknown() {
		labelsMap, _ := pool.Labels.Value(ctx)
		params.Labels = make(map[string]string)
		for k, v := range labelsMap {
			params.Labels[k] = v.ValueString()
		}
	}
	if !pool.Taints.IsNull() && !pool.Taints.IsUnknown() {
		taintsMap, _ := pool.Taints.Value(ctx)
		params.Taints = make(map[string]string)
		for k, v := range taintsMap {
			params.Taints[k] = v.ValueString()
		}
	}

	_, err := client.Cloud.K8S.Clusters.Pools.Update(
		ctx,
		poolName,
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		return fmt.Errorf("update pool %s: %w", poolName, err)
	}

	return nil
}

// deletePool deletes a pool from the cluster and waits for completion
func deletePool(ctx context.Context, client *gcore.Client, clusterName string, projectID, regionID param.Opt[int64], poolName string) error {
	err := client.Cloud.K8S.Clusters.Pools.DeleteAndPoll(
		ctx,
		poolName,
		cloud.K8SClusterPoolDeleteParams{
			ProjectID:   projectID,
			RegionID:    regionID,
			ClusterName: clusterName,
		},
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		return fmt.Errorf("delete pool %s: %w", poolName, err)
	}

	return nil
}

// poolsToMap converts pools slice to a map keyed by pool name
func poolsToMap(pools *[]*CloudK8SClusterPoolsModel) map[string]*CloudK8SClusterPoolsModel {
	result := make(map[string]*CloudK8SClusterPoolsModel)
	if pools == nil {
		return result
	}
	for _, pool := range *pools {
		if pool != nil {
			result[pool.Name.ValueString()] = pool
		}
	}
	return result
}

// poolsEqual checks if two pool slices are equal
func poolsEqual(a, b *[]*CloudK8SClusterPoolsModel) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(*a) != len(*b) {
		return false
	}
	oldPools := poolsToMap(a)
	newPools := poolsToMap(b)

	for name, oldPool := range oldPools {
		newPool, exists := newPools[name]
		if !exists {
			return false
		}
		if poolNeedsReplace(oldPool, newPool) || poolNeedsUpdate(oldPool, newPool) {
			return false
		}
	}
	return true
}

// poolNeedsReplace checks if pool changes require delete + create
func poolNeedsReplace(old, new *CloudK8SClusterPoolsModel) bool {
	if old.FlavorID.ValueString() != new.FlavorID.ValueString() {
		return true
	}
	if old.BootVolumeType.ValueString() != new.BootVolumeType.ValueString() {
		return true
	}
	if old.BootVolumeSize.ValueInt64() != new.BootVolumeSize.ValueInt64() {
		return true
	}
	if old.IsPublicIpv4.ValueBool() != new.IsPublicIpv4.ValueBool() {
		return true
	}
	if old.ServergroupPolicy.ValueString() != new.ServergroupPolicy.ValueString() {
		return true
	}
	// For computed map fields, treat null/unknown as equivalent to empty map
	if !mapsEqualOrBothEmpty(old.CrioConfig, new.CrioConfig) {
		return true
	}
	if !mapsEqualOrBothEmpty(old.KubeletConfig, new.KubeletConfig) {
		return true
	}
	return false
}

// mapsEqualOrBothEmpty returns true if maps are equal, or if both are "empty"
// (null, unknown, or zero-length). This handles the case where state has {}
// but plan has null for computed optional map fields.
func mapsEqualOrBothEmpty(a, b customfield.Map[types.String]) bool {
	aEmpty := a.IsNull() || a.IsUnknown() || len(a.Elements()) == 0
	bEmpty := b.IsNull() || b.IsUnknown() || len(b.Elements()) == 0
	if aEmpty && bEmpty {
		return true
	}
	return a.Equal(b)
}

// poolNeedsUpdate checks if pool changes can be done in-place
func poolNeedsUpdate(old, new *CloudK8SClusterPoolsModel) bool {
	if old.MinNodeCount.ValueInt64() != new.MinNodeCount.ValueInt64() {
		return true
	}
	if old.MaxNodeCount.ValueInt64() != new.MaxNodeCount.ValueInt64() {
		return true
	}
	if old.AutoHealingEnabled.ValueBool() != new.AutoHealingEnabled.ValueBool() {
		return true
	}
	// For computed map fields, treat null/unknown as equivalent to empty map
	if !mapsEqualOrBothEmpty(old.Labels, new.Labels) {
		return true
	}
	if !mapsEqualOrBothEmpty(old.Taints, new.Taints) {
		return true
	}
	return false
}

// checkPoolQuotaLimits checks if the pool changes would exceed quota limits
func checkPoolQuotaLimits(ctx context.Context, client *gcore.Client, projectID, regionID param.Opt[int64], oldPools, newPools map[string]*CloudK8SClusterPoolsModel) error {
	for poolName, newPool := range newPools {
		oldPool, exists := oldPools[poolName]

		var params cloud.K8SClusterPoolCheckQuotaParams

		if !exists || poolNeedsReplace(oldPool, newPool) {
			// New pool or pool replacement - check full quota
			params = cloud.K8SClusterPoolCheckQuotaParams{
				ProjectID:      projectID,
				RegionID:       regionID,
				FlavorID:       newPool.FlavorID.ValueString(),
				Name:           param.NewOpt(newPool.Name.ValueString()),
				MinNodeCount:   param.NewOpt(newPool.MinNodeCount.ValueInt64()),
				MaxNodeCount:   param.NewOpt(newPool.MaxNodeCount.ValueInt64()),
				BootVolumeSize: param.NewOpt(newPool.BootVolumeSize.ValueInt64()),
			}
			if !newPool.ServergroupPolicy.IsNull() && !newPool.ServergroupPolicy.IsUnknown() {
				params.ServergroupPolicy = cloud.K8SClusterPoolCheckQuotaParamsServergroupPolicy(newPool.ServergroupPolicy.ValueString())
			}
		} else if poolNeedsUpdate(oldPool, newPool) {
			// Pool update - check only the delta for node count changes
			minDelta := newPool.MinNodeCount.ValueInt64() - oldPool.MinNodeCount.ValueInt64()
			if minDelta <= 0 {
				// Scaling down or no change - no quota check needed
				continue
			}
			params = cloud.K8SClusterPoolCheckQuotaParams{
				ProjectID:      projectID,
				RegionID:       regionID,
				FlavorID:       newPool.FlavorID.ValueString(),
				Name:           param.NewOpt(newPool.Name.ValueString()),
				MinNodeCount:   param.NewOpt(minDelta),
				MaxNodeCount:   param.NewOpt(newPool.MaxNodeCount.ValueInt64() - oldPool.MaxNodeCount.ValueInt64()),
				BootVolumeSize: param.NewOpt(newPool.BootVolumeSize.ValueInt64()),
			}
			if !newPool.ServergroupPolicy.IsNull() && !newPool.ServergroupPolicy.IsUnknown() {
				params.ServergroupPolicy = cloud.K8SClusterPoolCheckQuotaParamsServergroupPolicy(newPool.ServergroupPolicy.ValueString())
			}
		} else {
			// No changes to this pool
			continue
		}

		quota, err := client.Cloud.K8S.Clusters.Pools.CheckQuota(ctx, params)
		if err != nil {
			return fmt.Errorf("check quota for pool %s: %w", poolName, err)
		}

		// Check if any quota would be exceeded
		exceeded := getExceededQuotas(quota)
		if len(exceeded) > 0 {
			return fmt.Errorf("quota limits exceeded for pool %s: %v", poolName, exceeded)
		}
	}

	return nil
}

// getExceededQuotas returns a list of quota names where usage + requested > limit
func getExceededQuotas(quota *cloud.K8SClusterPoolQuota) []string {
	var exceeded []string

	// Check each quota type
	if quota.CPUCountUsage+quota.CPUCountRequested > quota.CPUCountLimit && quota.CPUCountLimit > 0 {
		exceeded = append(exceeded, fmt.Sprintf("cpu_count (limit: %d, usage: %d, requested: %d)",
			quota.CPUCountLimit, quota.CPUCountUsage, quota.CPUCountRequested))
	}
	if quota.RamUsage+quota.RamRequested > quota.RamLimit && quota.RamLimit > 0 {
		exceeded = append(exceeded, fmt.Sprintf("ram (limit: %d, usage: %d, requested: %d)",
			quota.RamLimit, quota.RamUsage, quota.RamRequested))
	}
	if quota.VmCountUsage+quota.VmCountRequested > quota.VmCountLimit && quota.VmCountLimit > 0 {
		exceeded = append(exceeded, fmt.Sprintf("vm_count (limit: %d, usage: %d, requested: %d)",
			quota.VmCountLimit, quota.VmCountUsage, quota.VmCountRequested))
	}
	if quota.VolumeCountUsage+quota.VolumeCountRequested > quota.VolumeCountLimit && quota.VolumeCountLimit > 0 {
		exceeded = append(exceeded, fmt.Sprintf("volume_count (limit: %d, usage: %d, requested: %d)",
			quota.VolumeCountLimit, quota.VolumeCountUsage, quota.VolumeCountRequested))
	}
	if quota.VolumeSizeUsage+quota.VolumeSizeRequested > quota.VolumeSizeLimit && quota.VolumeSizeLimit > 0 {
		exceeded = append(exceeded, fmt.Sprintf("volume_size (limit: %d, usage: %d, requested: %d)",
			quota.VolumeSizeLimit, quota.VolumeSizeUsage, quota.VolumeSizeRequested))
	}
	if quota.ServergroupCountUsage+quota.ServergroupCountRequested > quota.ServergroupCountLimit && quota.ServergroupCountLimit > 0 {
		exceeded = append(exceeded, fmt.Sprintf("servergroup_count (limit: %d, usage: %d, requested: %d)",
			quota.ServergroupCountLimit, quota.ServergroupCountUsage, quota.ServergroupCountRequested))
	}
	if quota.FloatingCountUsage+quota.FloatingCountRequested > quota.FloatingCountLimit && quota.FloatingCountLimit > 0 {
		exceeded = append(exceeded, fmt.Sprintf("floating_count (limit: %d, usage: %d, requested: %d)",
			quota.FloatingCountLimit, quota.FloatingCountUsage, quota.FloatingCountRequested))
	}
	if quota.GPUCountUsage+quota.GPUCountRequested > quota.GPUCountLimit && quota.GPUCountLimit > 0 {
		exceeded = append(exceeded, fmt.Sprintf("gpu_count (limit: %d, usage: %d, requested: %d)",
			quota.GPUCountLimit, quota.GPUCountUsage, quota.GPUCountRequested))
	}

	return exceeded
}

// getPoolOrder returns a slice of pool names in their current order
func getPoolOrder(pools *[]*CloudK8SClusterPoolsModel) []string {
	if pools == nil {
		return nil
	}
	var order []string
	for _, pool := range *pools {
		if pool != nil {
			order = append(order, pool.Name.ValueString())
		}
	}
	return order
}

// reorderPoolsToMatch reorders pools to match the desired order
// Pools not in desiredOrder are appended at the end
func reorderPoolsToMatch(pools *[]*CloudK8SClusterPoolsModel, desiredOrder []string) *[]*CloudK8SClusterPoolsModel {
	if pools == nil {
		return nil
	}

	// Build a map for quick lookup
	poolsByName := make(map[string]*CloudK8SClusterPoolsModel)
	for _, pool := range *pools {
		if pool != nil {
			poolsByName[pool.Name.ValueString()] = pool
		}
	}

	// Build reordered list based on desired order
	var reordered []*CloudK8SClusterPoolsModel
	seen := make(map[string]bool)
	for _, name := range desiredOrder {
		if pool, exists := poolsByName[name]; exists {
			reordered = append(reordered, pool)
			seen[name] = true
		}
	}

	// Append any pools not in desired order (shouldn't happen, but be safe)
	for _, pool := range *pools {
		if pool != nil {
			name := pool.Name.ValueString()
			if !seen[name] {
				reordered = append(reordered, pool)
			}
		}
	}

	return &reordered
}
