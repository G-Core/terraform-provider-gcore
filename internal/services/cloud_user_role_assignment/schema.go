// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_user_role_assignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*CloudUserRoleAssignmentResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "Assignment ID",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"role": schema.StringAttribute{
				Description: "User role\nAvailable values: \"ClientAdministrator\", \"InternalNetworkOnlyUser\", \"Observer\", \"ProjectAdministrator\", \"User\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"ClientAdministrator",
						"InternalNetworkOnlyUser",
						"Observer",
						"ProjectAdministrator",
						"User",
					),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID",
				Required:    true,
			},
			"client_id": schema.Int64Attribute{
				Description: "Client ID. Required if `project_id` is specified",
				Optional:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"assigned_by": schema.Int64Attribute{
				Computed: true,
			},
			"assignment_id": schema.Int64Attribute{
				Description: "Assignment ID",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Created timestamp",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"updated_at": schema.StringAttribute{
				Description: "Updated timestamp",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *CloudUserRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudUserRoleAssignmentResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
