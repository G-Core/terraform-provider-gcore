// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_router

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*CloudNetworkRouterResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"region_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"interfaces": schema.ListNestedAttribute{
				Description: "List of interfaces to attach to router immediately after creation.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudNetworkRouterInterfacesModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"subnet_id": schema.StringAttribute{
							Description: "id of the subnet to attach to.",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "must be 'subnet'.\nAvailable values: \"subnet\".",
							Computed:    true,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("subnet"),
							},
							Default: stringdefault.StaticString("subnet"),
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description: "name of router",
				Required:    true,
			},
			"external_gateway_info": schema.SingleNestedAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: customfield.NewNestedObjectType[CloudNetworkRouterExternalGatewayInfoModel](ctx),
				Attributes: map[string]schema.Attribute{
					"network_id": schema.StringAttribute{
						Description:   "id of the external network.",
						Computed:      true,
						Optional:      true,
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
					"enable_snat": schema.BoolAttribute{
						Description: "Is SNAT enabled. Defaults to true.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(true),
					},
					"type": schema.StringAttribute{
						Description: "must be 'manual'.\nAvailable values: \"manual\", \"default\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("manual", "default"),
						},
					},
				},
			},
			"routes": schema.ListNestedAttribute{
				Description: "List of custom routes.",
				Computed:    true,
				Optional:    true,
				CustomType:    customfield.NewNestedObjectListType[CloudNetworkRouterRoutesModel](ctx),
			PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"destination": schema.StringAttribute{
							Description: "CIDR of destination IPv4 subnet.",
							Required:    true,
						},
						"nexthop": schema.StringAttribute{
							Description: "IPv4 address to forward traffic to if it's destination IP matches 'destination' CIDR.",
							Required:    true,
						},
					},
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the router was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"distributed": schema.BoolAttribute{
				Description: "Whether the router is distributed or centralized.",
				Computed:    true,
			PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"status": schema.StringAttribute{
				Description: "Status of the router.",
				Computed:    true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the router was last updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"tasks": schema.ListAttribute{
				Description: "List of task IDs representing asynchronous operations. Use these IDs to monitor operation progress:\n* `GET /v1/tasks/{task_id}` - Check individual task status and details\nPoll task status until completion (`FINISHED`/`ERROR`) before proceeding with dependent operations.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *CloudNetworkRouterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudNetworkRouterResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
