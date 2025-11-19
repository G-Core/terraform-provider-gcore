// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_floating_ip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*CloudFloatingIPResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "A floating IP is a static IP address that points to one of your Instances. It allows you to redirect network traffic to any of your Instances in the same datacenter.",
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
			"fixed_ip_address": schema.StringAttribute{
				Description:   "If the port has multiple IP addresses, a specific one can be selected using this field. If not specified, the first IP in the port's list will be used by default.",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"port_id": schema.StringAttribute{
				Description:   "If provided, the floating IP will be immediately attached to the specified port.",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"tags": schema.MapAttribute{
				Description: "Update key-value tags using JSON Merge Patch semantics (RFC 7386). Provide key-value pairs to add or update tags. Set tag values to `null` to remove tags. Unspecified tags remain unchanged. Read-only tags are always preserved and cannot be modified.\n\n**Examples:**\n\n* **Add/update tags:** `{'tags': {'environment': 'production', 'team': 'backend'}}` adds new tags or updates existing ones.\n\n* **Delete tags:** `{'tags': {'old_tag': null}}` removes specific tags.\n\n* **Remove all tags:** `{'tags': null}` removes all user-managed tags (read-only tags are preserved).\n\n* **Partial update:** `{'tags': {'environment': 'staging'}}` only updates specified tags.\n\n* **Mixed operations:** `{'tags': {'environment': 'production', 'cost_center': 'engineering', 'deprecated_tag': null}}` adds/updates 'environment' and '`cost_center`' while removing '`deprecated_tag`', preserving other existing tags.\n\n* **Replace all:** first delete existing tags with null values, then add new ones in the same request.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the floating IP was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"floating_ip_address": schema.StringAttribute{
				Description: "IP Address of the floating IP",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"router_id": schema.StringAttribute{
				Description: "Router ID",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Floating IP status. DOWN - unassigned (available). ACTIVE - attached to a port (in use). ERROR - error state.\nAvailable values: \"ACTIVE\", \"DOWN\", \"ERROR\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"ACTIVE",
						"DOWN",
						"ERROR",
					),
				},
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the floating IP was last updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
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

func (r *CloudFloatingIPResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudFloatingIPResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
