// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_placement_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*CloudPlacementGroupResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Placement Groups allow you to specific a policy that determines whether Virtual Machines will be hosted on the same physical server or on different ones.",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"region_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "The name of the server group.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"policy": schema.StringAttribute{
				Description: "The server group policy.\nAvailable values: \"affinity\", \"anti-affinity\", \"soft-anti-affinity\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"affinity",
						"anti-affinity",
						"soft-anti-affinity",
					),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"servergroup_id": schema.StringAttribute{
				Description: "The ID of the server group.",
				Computed:    true,
			},
			"instances": schema.ListNestedAttribute{
				Description: "The list of instances in this server group.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudPlacementGroupInstancesModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_id": schema.StringAttribute{
							Description: "The ID of the instance, corresponding to the attribute 'id'.",
							Computed:    true,
						},
						"instance_name": schema.StringAttribute{
							Description: "The name of the instance, corresponding to the attribute 'name'.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *CloudPlacementGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudPlacementGroupResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
