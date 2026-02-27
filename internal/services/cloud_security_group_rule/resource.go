// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group_rule

import (
	"context"
	"errors"
	"fmt"

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
var _ resource.ResourceWithConfigure = (*CloudSecurityGroupRuleResource)(nil)
var _ resource.ResourceWithModifyPlan = (*CloudSecurityGroupRuleResource)(nil)
var _ resource.ResourceWithImportState = (*CloudSecurityGroupRuleResource)(nil)

func NewResource() resource.Resource {
	return &CloudSecurityGroupRuleResource{}
}

// CloudSecurityGroupRuleResource defines the resource implementation.
type CloudSecurityGroupRuleResource struct {
	client *gcore.Client
}

func (r *CloudSecurityGroupRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_security_group_rule"
}

func (r *CloudSecurityGroupRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudSecurityGroupRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CloudSecurityGroupRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.SecurityGroupRuleNewParams{}

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
	rule, err := r.client.Cloud.SecurityGroups.Rules.NewAndPoll(
		ctx,
		data.GroupID.ValueString(),
		params,
		option.WithRequestBody("application/json", dataBytes),
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
	err = apijson.UnmarshalComputed([]byte(rule.RawJSON()), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudSecurityGroupRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update is not supported for this resource
}

func (r *CloudSecurityGroupRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CloudSecurityGroupRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.SecurityGroupGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	sg, err := r.client.Cloud.SecurityGroups.Get(
		ctx,
		data.GroupID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		var apierr *gcore.Error
		if errors.As(err, &apierr) && apierr.StatusCode == 404 {
			resp.Diagnostics.AddWarning("Resource not found", "The security group was not found on the server and the rule will be removed from state.")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	ruleID := data.ID.ValueString()
	found := false

	for _, rule := range sg.SecurityGroupRules {
		if rule.ID == ruleID {
			err = apijson.UnmarshalComputed([]byte(rule.RawJSON()), &data)
			if err != nil {
				resp.Diagnostics.AddError("failed to deserialize rule data", err.Error())
				return
			}
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddWarning("Rule not found",
			fmt.Sprintf("Security group rule with ID %s not found in security group %s. Removing from state.",
				ruleID, data.GroupID.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudSecurityGroupRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CloudSecurityGroupRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := cloud.SecurityGroupRuleDeleteParams{
		GroupID: data.GroupID.ValueString(),
	}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	err := r.client.Cloud.SecurityGroups.Rules.DeleteAndPoll(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		var apierr *gcore.Error
		if errors.As(err, &apierr) && apierr.StatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}
}

func (r *CloudSecurityGroupRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data = new(CloudSecurityGroupRuleModel)

	path_project_id := int64(0)
	path_region_id := int64(0)
	path_group_id := ""
	path_rule_id := ""
	diags := importpath.ParseImportID(
		req.ID,
		"<project_id>/<region_id>/<group_id>/<rule_id>",
		&path_project_id,
		&path_region_id,
		&path_group_id,
		&path_rule_id,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.GroupID = types.StringValue(path_group_id)
	data.ID = types.StringValue(path_rule_id)

	sg, err := r.client.Cloud.SecurityGroups.Get(
		ctx,
		path_group_id,
		cloud.SecurityGroupGetParams{
			ProjectID: param.NewOpt(path_project_id),
			RegionID:  param.NewOpt(path_region_id),
		},
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	found := false
	for _, rule := range sg.SecurityGroupRules {
		if rule.ID == path_rule_id {
			err = apijson.Unmarshal([]byte(rule.RawJSON()), &data)
			if err != nil {
				resp.Diagnostics.AddError("failed to deserialize rule data", err.Error())
				return
			}
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError("rule not found",
			fmt.Sprintf("Security group rule with ID %s not found in security group %s.",
				path_rule_id, path_group_id))
		return
	}

	// Re-set path fields since Unmarshal may not populate them
	data.ProjectID = types.Int64Value(path_project_id)
	data.RegionID = types.Int64Value(path_region_id)
	data.GroupID = types.StringValue(path_group_id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudSecurityGroupRuleResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {
}
