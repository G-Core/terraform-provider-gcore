// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_region

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudRegionsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				Description: "Filter regions by display name. Case-insensitive exact match.",
				Optional:    true,
			},
			"product": schema.StringAttribute{
				Description: "If defined then return only regions that support given product.\nAvailable values: \"containers\", \"inference\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("containers", "inference"),
				},
			},
			"order_by": schema.StringAttribute{
				Description: "Order by field and direction.\nAvailable values: \"created_at.asc\", \"created_at.desc\", \"display_name.asc\", \"display_name.desc\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"created_at.asc",
						"created_at.desc",
						"display_name.asc",
						"display_name.desc",
					),
				},
			},
			"show_volume_types": schema.BoolAttribute{
				Description: "If true, null `available_volume_type` is replaced with a list of available volume types.",
				Computed:    true,
				Optional:    true,
			},
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudRegionsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Region ID",
							Computed:    true,
						},
						"access_level": schema.StringAttribute{
							Description: "The access level of the region.\nAvailable values: \"core\", \"edge\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("core", "edge"),
							},
						},
						"available_volume_types": schema.ListAttribute{
							Description: "List of available volume types, 'standard', 'ssd_hiiops', 'cold'].",
							Computed:    true,
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
						},
						"coordinates": schema.SingleNestedAttribute{
							Description: "Coordinates of the region",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudRegionsCoordinatesDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"latitude": schema.StringAttribute{
									Computed: true,
								},
								"longitude": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"country": schema.StringAttribute{
							Description: "Two-letter country code, ISO 3166-1 alpha-2",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Region creation date and time",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"created_on": schema.StringAttribute{
							Description:        "This field is deprecated. Use `created_at` instead.",
							Computed:           true,
							DeprecationMessage: "This attribute is deprecated.",
							CustomType:         timetypes.RFC3339Type{},
						},
						"display_name": schema.StringAttribute{
							Description: "Human-readable region name",
							Computed:    true,
						},
						"endpoint_type": schema.StringAttribute{
							Description: "Endpoint type\nAvailable values: \"admin\", \"internal\", \"public\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"admin",
									"internal",
									"public",
								),
							},
						},
						"external_network_id": schema.StringAttribute{
							Description: "External network ID for Neutron",
							Computed:    true,
						},
						"file_share_types": schema.ListAttribute{
							Description: "List of available file share types",
							Computed:    true,
							Validators: []validator.List{
								listvalidator.ValueStringsAre(
									stringvalidator.OneOfCaseInsensitive("standard", "vast"),
								),
							},
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
						},
						"has_ai": schema.BoolAttribute{
							Description: "Region has AI capability",
							Computed:    true,
						},
						"has_ai_gpu": schema.BoolAttribute{
							Description: "Region has AI GPU capability",
							Computed:    true,
						},
						"has_baremetal": schema.BoolAttribute{
							Description: "Region has bare metal capability",
							Computed:    true,
						},
						"has_basic_vm": schema.BoolAttribute{
							Description: "Region has basic vm capability",
							Computed:    true,
						},
						"has_dbaas": schema.BoolAttribute{
							Description: "Region has DBAAS service",
							Computed:    true,
						},
						"has_ddos": schema.BoolAttribute{
							Description: "Region has Advanced DDoS Protection capability",
							Computed:    true,
						},
						"has_k8s": schema.BoolAttribute{
							Description: "Region has managed kubernetes capability",
							Computed:    true,
						},
						"has_kvm": schema.BoolAttribute{
							Description: "Region has KVM virtualization capability",
							Computed:    true,
						},
						"has_sfs": schema.BoolAttribute{
							Description: "Region has SFS capability",
							Computed:    true,
						},
						"keystone_id": schema.Int64Attribute{
							Description: "Foreign key to Keystone entity",
							Computed:    true,
						},
						"keystone_name": schema.StringAttribute{
							Description: "Technical region name",
							Computed:    true,
						},
						"metrics_database_id": schema.Int64Attribute{
							Description: "Foreign key to Metrics database entity",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "Region state\nAvailable values: \"ACTIVE\", \"DELETED\", \"DELETING\", \"DELETION_FAILED\", \"INACTIVE\", \"MAINTENANCE\", \"NEW\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"ACTIVE",
									"DELETED",
									"DELETING",
									"DELETION_FAILED",
									"INACTIVE",
									"MAINTENANCE",
									"NEW",
								),
							},
						},
						"task_id": schema.StringAttribute{
							Description:        "This field is deprecated and can be ignored",
							Computed:           true,
							DeprecationMessage: "This attribute is deprecated.",
						},
						"vlan_physical_network": schema.StringAttribute{
							Description: "Physical network name to create vlan networks",
							Computed:    true,
						},
						"zone": schema.StringAttribute{
							Description: "Geographical zone\nAvailable values: \"AMERICAS\", \"APAC\", \"EMEA\", \"RUSSIA_AND_CIS\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"AMERICAS",
									"APAC",
									"EMEA",
									"RUSSIA_AND_CIS",
								),
							},
						},
						"ddos_endpoint_id": schema.Int64Attribute{
							Description:        "DDoS endpoint ID",
							Computed:           true,
							DeprecationMessage: "This attribute is deprecated.",
						},
					},
				},
			},
		},
	}
}

func (d *CloudRegionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudRegionsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
