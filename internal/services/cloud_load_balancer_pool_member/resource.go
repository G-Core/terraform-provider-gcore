// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_member

import (
	"context"
	"fmt"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/logging"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudLoadBalancerPoolMemberResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudLoadBalancerPoolMemberResource)(nil)

func NewResource() resource.Resource {
	return &CloudLoadBalancerPoolMemberResource{}
}

// CloudLoadBalancerPoolMemberResource defines the resource implementation.
type CloudLoadBalancerPoolMemberResource struct {
	client *gcore.Client
}

func (r *CloudLoadBalancerPoolMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_load_balancer_pool_member"
}

func (r *CloudLoadBalancerPoolMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudLoadBalancerPoolMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudLoadBalancerPoolMemberModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.LoadBalancerPoolMemberAddParams{}

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
	member, err := r.client.Cloud.LoadBalancers.Pools.Members.AddAndPoll(
		ctx,
		data.PoolID.ValueString(),
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	err = apijson.UnmarshalComputed([]byte(member.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerPoolMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudLoadBalancerPoolMemberModel
	var state *CloudLoadBalancerPoolMemberModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

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

	// 1. Get the pool to access all members
	pool, err := r.client.Cloud.LoadBalancers.Pools.Get(
		ctx,
		data.PoolID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to get pool for member update", err.Error())
		return
	}

	// 2. Rebuild members list with updated member
	var updatedMembers []cloud.LoadBalancerPoolUpdateParamsMember
	memberFound := false

	for _, member := range pool.Members {
		if member.ID == data.ID.ValueString() {
			// This is the member we're updating
			memberFound = true
			updatedMember := cloud.LoadBalancerPoolUpdateParamsMember{
				Address:      data.Address.ValueString(),
				ProtocolPort: data.ProtocolPort.ValueInt64(),
			}

			if !data.InstanceID.IsNull() {
				updatedMember.InstanceID = param.NewOpt(data.InstanceID.ValueString())
			}
			if !data.MonitorAddress.IsNull() {
				updatedMember.MonitorAddress = param.NewOpt(data.MonitorAddress.ValueString())
			}
			if !data.MonitorPort.IsNull() {
				updatedMember.MonitorPort = param.NewOpt(data.MonitorPort.ValueInt64())
			}
			if !data.SubnetID.IsNull() {
				updatedMember.SubnetID = param.NewOpt(data.SubnetID.ValueString())
			}
			if !data.AdminStateUp.IsNull() {
				updatedMember.AdminStateUp = param.NewOpt(data.AdminStateUp.ValueBool())
			}
			if !data.Backup.IsNull() {
				updatedMember.Backup = param.NewOpt(data.Backup.ValueBool())
			}
			if !data.Weight.IsNull() {
				updatedMember.Weight = param.NewOpt(data.Weight.ValueInt64())
			}

			updatedMembers = append(updatedMembers, updatedMember)
		} else {
			// Keep other members unchanged
			existingMember := cloud.LoadBalancerPoolUpdateParamsMember{
				Address:      member.Address,
				ProtocolPort: member.ProtocolPort,
			}

			if member.MonitorAddress != "" {
				existingMember.MonitorAddress = param.NewOpt(member.MonitorAddress)
			}
			if member.MonitorPort != 0 {
				existingMember.MonitorPort = param.NewOpt(member.MonitorPort)
			}
			if member.SubnetID != "" {
				existingMember.SubnetID = param.NewOpt(member.SubnetID)
			}
			existingMember.AdminStateUp = param.NewOpt(member.AdminStateUp)
			existingMember.Backup = param.NewOpt(member.Backup)
			if member.Weight != 0 {
				existingMember.Weight = param.NewOpt(member.Weight)
			}

			updatedMembers = append(updatedMembers, existingMember)
		}
	}

	if !memberFound {
		resp.Diagnostics.AddError("member not found in pool",
			fmt.Sprintf("member ID %s not found in pool %s", data.ID.ValueString(), data.PoolID.ValueString()))
		return
	}

	// 3. Update the pool with the new members list
	updateParams := cloud.LoadBalancerPoolUpdateParams{
		Name:    param.NewOpt(pool.Name),
		Members: updatedMembers,
	}
	if !data.ProjectID.IsNull() {
		updateParams.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}
	if !data.RegionID.IsNull() {
		updateParams.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	poolUpdated, err := r.client.Cloud.LoadBalancers.Pools.UpdateAndPoll(
		ctx,
		data.PoolID.ValueString(),
		updateParams,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to update load balancer pool member", err.Error())
		return
	}

	// 4. Read back the updated member data
	for _, member := range poolUpdated.Members {
		if member.ID == data.ID.ValueString() {
			err = apijson.UnmarshalComputed([]byte(member.RawJSON()), &data)
			if err != nil {
				resp.Diagnostics.AddError("failed to deserialize updated member", err.Error())
				return
			}
			break
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerPoolMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudLoadBalancerPoolMemberModel

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

	// Get the pool to access its members list
	pool, err := r.client.Cloud.LoadBalancers.Pools.Get(
		ctx,
		data.PoolID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to get load balancer pool", err.Error())
		return
	}

	// Find the member in the pool's members list
	memberID := data.ID.ValueString()
	found := false

	for _, member := range pool.Members {
		if member.ID == memberID {
			// Update the data model with the member's current state
			err = apijson.UnmarshalComputed([]byte(member.RawJSON()), &data)
			if err != nil {
				resp.Diagnostics.AddError("failed to deserialize member data", err.Error())
				return
			}
			found = true
			break
		}
	}

	if !found {
		// Member not found in pool - remove from state
		resp.Diagnostics.AddWarning("Member not found",
			fmt.Sprintf("Load balancer pool member with ID %s not found in pool %s. Removing from state.",
				memberID, data.PoolID.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerPoolMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudLoadBalancerPoolMemberModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.LoadBalancerPoolMemberRemoveParams{
		PoolID: data.PoolID.ValueString(),
	}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.LoadBalancers.Pools.Members.RemoveAndPoll(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to delete load balancer pool member", err.Error())
		return
	}
}

func (r *CloudLoadBalancerPoolMemberResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {
}
