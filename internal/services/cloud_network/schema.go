// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*CloudNetworkResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.Int64Attribute{
				Description:   "Project ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"region_id": schema.Int64Attribute{
				Description:   "Region ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"create_router": schema.BoolAttribute{
				Description:   "Defaults to True",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplaceIfConfigured()},
				Default:       booldefault.StaticBool(true),
			},
			"type": schema.StringAttribute{
				Description: "vlan or vxlan network type is allowed. Default value is vxlan\nAvailable values: \"vlan\", \"vxlan\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("vlan", "vxlan"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
				Default:       stringdefault.StaticString("vxlan"),
			},
			"name": schema.StringAttribute{
				Description: "Name.",
				Optional:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Update key-value tags using JSON Merge Patch semantics (RFC 7386). Provide key-value pairs to add or update tags. Set tag values to `null` to remove tags. Unspecified tags remain unchanged. Read-only tags are always preserved and cannot be modified.\n\n**Examples:**\n\n* **Add/update tags:** `{'tags': {'environment': 'production', 'team': 'backend'}}` adds new tags or updates existing ones.\n\n* **Delete tags:** `{'tags': {'old_tag': null}}` removes specific tags.\n\n* **Remove all tags:** `{'tags': null}` removes all user-managed tags (read-only tags are preserved).\n\n* **Partial update:** `{'tags': {'environment': 'staging'}}` only updates specified tags.\n\n* **Mixed operations:** `{'tags': {'environment': 'production', 'cost_center': 'engineering', 'deprecated_tag': null}}` adds/updates 'environment' and '`cost_center`' while removing '`deprecated_tag`', preserving other existing tags.\n\n* **Replace all:** first delete existing tags with null values, then add new ones in the same request.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the network was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"default": schema.BoolAttribute{
				Description: "True if network has `is_default` attribute",
				Computed:    true,
			},
			"external": schema.BoolAttribute{
				Description: "True if the network `router:external` attribute",
				Computed:    true,
			},
			"mtu": schema.Int64Attribute{
				Description: "MTU (maximum transmission unit). Default value is 1450",
				Computed:    true,
			},
			"port_security_enabled": schema.BoolAttribute{
				Description: "Indicates `port_security_enabled` status of all newly created in the network ports.",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"segmentation_id": schema.Int64Attribute{
				Description: "Id of network segment",
				Computed:    true,
			},
			"shared": schema.BoolAttribute{
				Description: "True when the network is shared with your project by external owner",
				Computed:    true,
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the network was last updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"subnets": schema.ListAttribute{
				Description: "List of subnetworks",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"tasks": schema.ListAttribute{
				Description: "List of task IDs representing asynchronous operations. Use these IDs to monitor operation progress:\n* `GET /v1/tasks/{task_id}` - Check individual task status and details\nPoll task status until completion (`FINISHED`/`ERROR`) before proceeding with dependent operations.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
		},
	}
}

func (r *CloudNetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudNetworkResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
