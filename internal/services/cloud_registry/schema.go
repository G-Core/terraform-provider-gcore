// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*CloudRegistryResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "Registry ID",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown(), int64planmodifier.RequiresReplace()},
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
				Description:   "A name for the container registry.\nShould be in lowercase, consisting only of numbers, letters and -,\nwith maximum length of 24 characters",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"storage_limit": schema.Int64Attribute{
				Description: "Registry storage limit, GiB",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 1000),
				},
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
				Default:       int64default.StaticInt64(5),
			},
			"created_at": schema.StringAttribute{
				Description: "Registry creation date-time",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"repo_count": schema.Int64Attribute{
				Description: "Number of repositories in the registry",
				Computed:    true,
			},
			"storage_used": schema.Int64Attribute{
				Description: "Registry storage used, bytes",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Registry modification date-time",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"url": schema.StringAttribute{
				Description: "Registry url",
				Computed:    true,
			},
		},
	}
}

func (r *CloudRegistryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudRegistryResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
