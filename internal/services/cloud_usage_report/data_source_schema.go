// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_usage_report

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudUsageReportDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"time_from": schema.StringAttribute{
				Description: "The start date of the report period (ISO 8601). The report starts from the beginning of this day.",
				Required:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"time_to": schema.StringAttribute{
				Description: "The end date of the report period (ISO 8601). The report ends just before the beginning of this day.",
				Required:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"projects": schema.ListAttribute{
				Description: "List of project IDs",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"regions": schema.ListAttribute{
				Description: "List of region IDs.",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"types": schema.ListAttribute{
				Description: "List of resource types to be filtered in the report.",
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOfCaseInsensitive(
							"ai_cluster",
							"ai_virtual_cluster",
							"backup",
							"baremetal",
							"basic_vm",
							"containers",
							"dbaas_postgresql_connection_pooler",
							"dbaas_postgresql_cpu",
							"dbaas_postgresql_memory",
							"dbaas_postgresql_public_network",
							"dbaas_postgresql_volume",
							"egress_traffic",
							"external_ip",
							"file_share",
							"floatingip",
							"functions",
							"functions_calls",
							"functions_traffic",
							"image",
							"inference",
							"instance",
							"load_balancer",
							"log_index",
							"snapshot",
							"volume",
						),
					),
				},
				ElementType: types.StringType,
			},
			"schema_filter": schema.SingleNestedAttribute{
				Description: "Extended filter for field filtering.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"field": schema.StringAttribute{
						Description: "Field name to filter by\nAvailable values: \"last_name\", \"last_size\", \"source_volume_uuid\", \"type\", \"uuid\", \"volume_type\", \"flavor\", \"attached_to_vm\", \"file_share_type\", \"ip_address\", \"instance_name\", \"instance_type\", \"port_id\", \"vm_id\", \"network_id\", \"subnet_id\", \"schedule_id\".",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"last_name",
								"last_size",
								"source_volume_uuid",
								"type",
								"uuid",
								"volume_type",
								"flavor",
								"attached_to_vm",
								"file_share_type",
								"ip_address",
								"instance_name",
								"instance_type",
								"port_id",
								"vm_id",
								"network_id",
								"subnet_id",
								"schedule_id",
							),
						},
					},
					"type": schema.StringAttribute{
						Description: `Available values: "snapshot", "instance", "ai_cluster", "ai_virtual_cluster", "basic_vm", "baremetal", "volume", "file_share", "image", "floatingip", "egress_traffic", "load_balancer", "external_ip", "backup", "log_index", "functions", "functions_calls", "functions_traffic", "containers", "inference", "dbaas_postgresql_volume", "dbaas_postgresql_public_network", "dbaas_postgresql_cpu", "dbaas_postgresql_memory", "dbaas_postgresql_connection_pooler".`,
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"snapshot",
								"instance",
								"ai_cluster",
								"ai_virtual_cluster",
								"basic_vm",
								"baremetal",
								"volume",
								"file_share",
								"image",
								"floatingip",
								"egress_traffic",
								"load_balancer",
								"external_ip",
								"backup",
								"log_index",
								"functions",
								"functions_calls",
								"functions_traffic",
								"containers",
								"inference",
								"dbaas_postgresql_volume",
								"dbaas_postgresql_public_network",
								"dbaas_postgresql_cpu",
								"dbaas_postgresql_memory",
								"dbaas_postgresql_connection_pooler",
							),
						},
					},
					"values": schema.ListAttribute{
						Description: "List of field values to filter",
						Required:    true,
						ElementType: types.StringType,
					},
				},
			},
			"tags": schema.SingleNestedAttribute{
				Description: "Filter by tags",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"conditions": schema.ListNestedAttribute{
						Description: "A list of tag filtering conditions defining how tags should match.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Description: "The name of the tag to filter (e.g., '`os_version`').",
									Optional:    true,
								},
								"strict": schema.BoolAttribute{
									Description: "Determines how strictly the tag value must match the specified value. If true, the tag value must exactly match the given value. If false, a less strict match (e.g., partial or case-insensitive match) may be applied.",
									Optional:    true,
								},
								"value": schema.StringAttribute{
									Description: "The value of the tag to filter (e.g., '22.04').",
									Optional:    true,
								},
							},
						},
					},
					"condition_type": schema.StringAttribute{
						Description: "Specifies whether conditions are combined using OR (default) or AND logic.\nAvailable values: \"AND\", \"OR\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("AND", "OR"),
						},
					},
				},
			},
			"enable_last_day": schema.BoolAttribute{
				Description: "Expenses for the last specified day are taken into account. As the default, False.",
				Computed:    true,
				Optional:    true,
			},
			"limit": schema.Int64Attribute{
				Description: "The response resources limit. Defaults to 10.",
				Computed:    true,
				Optional:    true,
			},
			"offset": schema.Int64Attribute{
				Description: "The response resources offset.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"sorting": schema.ListNestedAttribute{
				Description: "List of sorting filters (JSON objects) fields: project. directions: asc, desc.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudUsageReportSortingDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"billing_value": schema.StringAttribute{
							Description: `Available values: "asc", "desc".`,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("asc", "desc"),
							},
						},
						"first_seen": schema.StringAttribute{
							Description: `Available values: "asc", "desc".`,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("asc", "desc"),
							},
						},
						"last_name": schema.StringAttribute{
							Description: `Available values: "asc", "desc".`,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("asc", "desc"),
							},
						},
						"last_seen": schema.StringAttribute{
							Description: `Available values: "asc", "desc".`,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("asc", "desc"),
							},
						},
						"project": schema.StringAttribute{
							Description: `Available values: "asc", "desc".`,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("asc", "desc"),
							},
						},
						"region": schema.StringAttribute{
							Description: `Available values: "asc", "desc".`,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("asc", "desc"),
							},
						},
						"type": schema.StringAttribute{
							Description: `Available values: "asc", "desc".`,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("asc", "desc"),
							},
						},
					},
				},
			},
			"count": schema.Int64Attribute{
				Description: "Total count of the resources",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"resources": schema.ListNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectListType[CloudUsageReportResourcesDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"billing_metric_name": schema.StringAttribute{
							Description: "Name of the billing metric",
							Computed:    true,
						},
						"billing_value": schema.Float64Attribute{
							Description: "Value of the billing metric",
							Computed:    true,
						},
						"billing_value_unit": schema.StringAttribute{
							Description: "Unit of billing value",
							Computed:    true,
						},
						"first_seen": schema.StringAttribute{
							Description: "First time the resource was seen in the given period",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"flavor": schema.StringAttribute{
							Description: "Flavor of the Baremetal GPU cluster",
							Computed:    true,
						},
						"last_name": schema.StringAttribute{
							Description: "Name of the AI cluster",
							Computed:    true,
						},
						"last_seen": schema.StringAttribute{
							Description: "Last time the resource was seen in the given period",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"project_id": schema.Int64Attribute{
							Description: "ID of the project the resource belongs to",
							Computed:    true,
						},
						"region": schema.Int64Attribute{
							Description: "Region ID",
							Computed:    true,
						},
						"region_id": schema.Int64Attribute{
							Description: "Region ID",
							Computed:    true,
						},
						"tags": schema.ListAttribute{
							Description: "List of tags",
							Computed:    true,
							CustomType:  customfield.NewListType[customfield.Map[types.String]](ctx),
							ElementType: types.MapType{
								ElemType: types.StringType,
							},
						},
						"type": schema.StringAttribute{
							Description: `Available values: "ai_cluster", "ai_virtual_cluster", "baremetal", "basic_vm", "backup", "containers", "egress_traffic", "external_ip", "file_share", "floatingip", "functions", "functions_calls", "functions_traffic", "image", "inference", "instance", "load_balancer", "log_index", "snapshot", "volume", "dbaas_postgresql_connection_pooler", "dbaas_postgresql_memory", "dbaas_postgresql_public_network", "dbaas_postgresql_cpu", "dbaas_postgresql_volume".`,
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"ai_cluster",
									"ai_virtual_cluster",
									"baremetal",
									"basic_vm",
									"backup",
									"containers",
									"egress_traffic",
									"external_ip",
									"file_share",
									"floatingip",
									"functions",
									"functions_calls",
									"functions_traffic",
									"image",
									"inference",
									"instance",
									"load_balancer",
									"log_index",
									"snapshot",
									"volume",
									"dbaas_postgresql_connection_pooler",
									"dbaas_postgresql_memory",
									"dbaas_postgresql_public_network",
									"dbaas_postgresql_cpu",
									"dbaas_postgresql_volume",
								),
							},
						},
						"uuid": schema.StringAttribute{
							Description: "UUID of the Baremetal GPU cluster",
							Computed:    true,
						},
						"last_size": schema.Int64Attribute{
							Description: "Size of the backup in bytes",
							Computed:    true,
						},
						"schedule_id": schema.StringAttribute{
							Description: "ID of the backup schedule",
							Computed:    true,
						},
						"source_volume_uuid": schema.StringAttribute{
							Description: "UUID of the source volume",
							Computed:    true,
						},
						"instance_type": schema.StringAttribute{
							Description: "Type of the instance\nAvailable values: \"baremetal\", \"vm\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("baremetal", "vm"),
							},
						},
						"port_id": schema.StringAttribute{
							Description: "ID of the port the traffic is associated with",
							Computed:    true,
						},
						"size_unit": schema.StringAttribute{
							Description: "Unit of size",
							Computed:    true,
						},
						"vm_id": schema.StringAttribute{
							Description: "ID of the bare metal server the traffic is associated with",
							Computed:    true,
						},
						"instance_name": schema.StringAttribute{
							Description: "Name of the instance",
							Computed:    true,
						},
						"attached_to_vm": schema.StringAttribute{
							Description: "ID of the VM the IP is attached to",
							Computed:    true,
						},
						"ip_address": schema.StringAttribute{
							Description: "IP address",
							Computed:    true,
						},
						"network_id": schema.StringAttribute{
							Description: "ID of the network the IP is attached to",
							Computed:    true,
						},
						"subnet_id": schema.StringAttribute{
							Description: "ID of the subnet the IP is attached to",
							Computed:    true,
						},
						"file_share_type": schema.StringAttribute{
							Description: "Type of the file share",
							Computed:    true,
						},
						"volume_type": schema.StringAttribute{
							Description: "Type of the volume",
							Computed:    true,
						},
					},
				},
			},
			"totals": schema.ListNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectListType[CloudUsageReportTotalsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"billing_metric_name": schema.StringAttribute{
							Description: "Name of the billing metric",
							Computed:    true,
						},
						"billing_value": schema.Float64Attribute{
							Description: "Value of the billing metric",
							Computed:    true,
						},
						"billing_value_unit": schema.StringAttribute{
							Description: "Unit of billing value",
							Computed:    true,
						},
						"flavor": schema.StringAttribute{
							Description: "Flavor of the Baremetal GPU cluster",
							Computed:    true,
						},
						"region": schema.Int64Attribute{
							Description: "Region ID",
							Computed:    true,
						},
						"region_id": schema.Int64Attribute{
							Description: "Region ID",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: `Available values: "ai_cluster", "ai_virtual_cluster", "baremetal", "basic_vm", "containers", "egress_traffic", "external_ip", "file_share", "floatingip", "functions", "functions_calls", "functions_traffic", "image", "inference", "instance", "load_balancer", "log_index", "snapshot", "volume", "dbaas_postgresql_connection_pooler", "dbaas_postgresql_memory", "dbaas_postgresql_public_network", "dbaas_postgresql_cpu", "dbaas_postgresql_volume".`,
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"ai_cluster",
									"ai_virtual_cluster",
									"baremetal",
									"basic_vm",
									"containers",
									"egress_traffic",
									"external_ip",
									"file_share",
									"floatingip",
									"functions",
									"functions_calls",
									"functions_traffic",
									"image",
									"inference",
									"instance",
									"load_balancer",
									"log_index",
									"snapshot",
									"volume",
									"dbaas_postgresql_connection_pooler",
									"dbaas_postgresql_memory",
									"dbaas_postgresql_public_network",
									"dbaas_postgresql_cpu",
									"dbaas_postgresql_volume",
								),
							},
						},
						"instance_type": schema.StringAttribute{
							Description: "Type of the instance\nAvailable values: \"baremetal\", \"vm\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("baremetal", "vm"),
							},
						},
						"file_share_type": schema.StringAttribute{
							Description: "Type of the file share",
							Computed:    true,
						},
						"volume_type": schema.StringAttribute{
							Description: "Type of the volume",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudUsageReportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudUsageReportDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
