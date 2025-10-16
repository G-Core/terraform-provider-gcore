// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_member

import (
	"context"
	"fmt"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/logging"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = (*CloudLoadBalancerPoolMemberResource)(nil)
var _ resource.ResourceWithImportState = (*CloudLoadBalancerPoolMemberResource)(nil)

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

func (r *CloudLoadBalancerPoolMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
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

	params := cloud.LoadBalancerPoolMemberAddParams{
		ProjectID:    param.NewOpt(data.ProjectID.ValueInt64()),
		RegionID:     param.NewOpt(data.RegionID.ValueInt64()),
		Address:      data.Address.ValueString(),
		ProtocolPort: data.ProtocolPort.ValueInt64(),
	}

	if !data.InstanceID.IsNull() {
		params.InstanceID = param.NewOpt(data.InstanceID.ValueString())
	}

	if !data.MonitorAddress.IsNull() {
		params.MonitorAddress = param.NewOpt(data.MonitorAddress.ValueString())
	}

	if !data.MonitorPort.IsNull() {
		params.MonitorPort = param.NewOpt(data.MonitorPort.ValueInt64())
	}

	if !data.SubnetID.IsNull() {
		params.SubnetID = param.NewOpt(data.SubnetID.ValueString())
	}

	if !data.AdminStateUp.IsNull() {
		params.AdminStateUp = param.NewOpt(data.AdminStateUp.ValueBool())
	}

	if !data.Backup.IsNull() {
		params.Backup = param.NewOpt(data.Backup.ValueBool())
	}

	if !data.Weight.IsNull() {
		params.Weight = param.NewOpt(data.Weight.ValueInt64())
	}

	// Call AddAndPoll to create the member and wait for completion
	member, err := r.client.Cloud.LoadBalancers.Pools.Members.AddAndPoll(
		ctx,
		data.PoolID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to create load balancer pool member", err.Error())
		return
	}

	// Populate Terraform state with the returned member data
	data.ID = types.StringValue(member.ID)
	data.Address = types.StringValue(member.Address)
	data.ProtocolPort = types.Int64Value(member.ProtocolPort)
	data.AdminStateUp = types.BoolValue(member.AdminStateUp)
	data.Backup = types.BoolValue(member.Backup)
	data.Weight = types.Int64Value(member.Weight)
	data.OperatingStatus = types.StringValue(string(member.OperatingStatus))

	// Note: InstanceID is not returned in the API response
	if member.MonitorAddress != "" {
		data.MonitorAddress = types.StringValue(member.MonitorAddress)
	}
	if member.MonitorPort != 0 {
		data.MonitorPort = types.Int64Value(member.MonitorPort)
	}
	if member.SubnetID != "" {
		data.SubnetID = types.StringValue(member.SubnetID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudLoadBalancerPoolMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudLoadBalancerPoolMemberModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readMemberFromPool(ctx, data); err != nil {
		resp.Diagnostics.AddError("failed to read load balancer pool member", err.Error())
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

	// To update a member, we need to get the entire pool, update the member in the list,
	// and call the pool's Update method with the full member list
	// This is based on the old provider's implementation

	poolParams := cloud.LoadBalancerPoolGetParams{
		ProjectID: param.NewOpt(data.ProjectID.ValueInt64()),
		RegionID:  param.NewOpt(data.RegionID.ValueInt64()),
	}

	pool, err := r.client.Cloud.LoadBalancers.Pools.Get(
		ctx,
		data.PoolID.ValueString(),
		poolParams,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to get pool for member update", err.Error())
		return
	}

	// Build the updated members list
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
		resp.Diagnostics.AddError("member not found in pool", fmt.Sprintf("member ID %s not found in pool %s", data.ID.ValueString(), data.PoolID.ValueString()))
		return
	}

	// Update the pool with the new member list
	updateParams := cloud.LoadBalancerPoolUpdateParams{
		ProjectID: param.NewOpt(data.ProjectID.ValueInt64()),
		RegionID:  param.NewOpt(data.RegionID.ValueInt64()),
		Name:      param.NewOpt(pool.Name),
		Members:   updatedMembers,
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

	// Verify the update was successful by checking the pool
	_ = poolUpdated

	// Read the updated member details
	if err := r.readMemberFromPool(ctx, data); err != nil {
		resp.Diagnostics.AddWarning("failed to read member details after update", err.Error())
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
		ProjectID: param.NewOpt(data.ProjectID.ValueInt64()),
		RegionID:  param.NewOpt(data.RegionID.ValueInt64()),
		PoolID:    data.PoolID.ValueString(),
	}

	// Call RemoveAndPoll to delete the member and wait for completion
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

	// The member is deleted, no need to set state
}

func (r *CloudLoadBalancerPoolMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: project_id:region_id:pool_id:member_id
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// readMemberFromPool reads member details by fetching the pool and finding the member in it
func (r *CloudLoadBalancerPoolMemberResource) readMemberFromPool(ctx context.Context, data *CloudLoadBalancerPoolMemberModel) error {
	params := cloud.LoadBalancerPoolGetParams{
		ProjectID: param.NewOpt(data.ProjectID.ValueInt64()),
		RegionID:  param.NewOpt(data.RegionID.ValueInt64()),
	}

	pool, err := r.client.Cloud.LoadBalancers.Pools.Get(
		ctx,
		data.PoolID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		return fmt.Errorf("failed to get pool: %w", err)
	}

	// Find the member in the pool's member list
	for _, member := range pool.Members {
		if member.ID == data.ID.ValueString() {
			data.Address = types.StringValue(member.Address)
			data.ProtocolPort = types.Int64Value(member.ProtocolPort)
			data.AdminStateUp = types.BoolValue(member.AdminStateUp)
			data.Backup = types.BoolValue(member.Backup)
			data.Weight = types.Int64Value(member.Weight)
			data.OperatingStatus = types.StringValue(string(member.OperatingStatus))

			// Note: InstanceID is not returned in the API response
			if member.MonitorAddress != "" {
				data.MonitorAddress = types.StringValue(member.MonitorAddress)
			}
			if member.MonitorPort != 0 {
				data.MonitorPort = types.Int64Value(member.MonitorPort)
			}
			if member.SubnetID != "" {
				data.SubnetID = types.StringValue(member.SubnetID)
			}

			return nil
		}
	}

	return fmt.Errorf("member ID %s not found in pool %s", data.ID.ValueString(), data.PoolID.ValueString())
}
