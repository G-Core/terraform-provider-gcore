// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_quota_request

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

var _ datasource.DataSourceWithConfigValidators = (*CloudQuotaRequestsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"status": schema.ListAttribute{
				Description: "List of limit requests statuses for filtering",
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOfCaseInsensitive(
							"done",
							"in progress",
							"rejected",
						),
					),
				},
				ElementType: types.StringType,
			},
			"limit": schema.Int64Attribute{
				Description: "Optional. Limit the number of returned items",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtMost(1000),
				},
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
				CustomType:  customfield.NewNestedObjectListType[CloudQuotaRequestsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Request ID",
							Computed:    true,
						},
						"client_id": schema.Int64Attribute{
							Description: "Client ID",
							Computed:    true,
						},
						"requested_limits": schema.SingleNestedAttribute{
							Description: "Requested limits.",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudQuotaRequestsRequestedLimitsDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"global_limits": schema.SingleNestedAttribute{
									Description: "Global entity quota limits",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectType[CloudQuotaRequestsRequestedLimitsGlobalLimitsDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"inference_cpu_millicore_count_limit": schema.Int64Attribute{
											Description: "Inference CPU millicore count limit",
											Computed:    true,
										},
										"inference_gpu_a100_count_limit": schema.Int64Attribute{
											Description: "Inference GPU A100 Count limit",
											Computed:    true,
										},
										"inference_gpu_h100_count_limit": schema.Int64Attribute{
											Description: "Inference GPU H100 Count limit",
											Computed:    true,
										},
										"inference_gpu_l40s_count_limit": schema.Int64Attribute{
											Description: "Inference GPU L40s Count limit",
											Computed:    true,
										},
										"inference_instance_count_limit": schema.Int64Attribute{
											Description: "Inference instance count limit",
											Computed:    true,
										},
										"keypair_count_limit": schema.Int64Attribute{
											Description: "SSH Keys Count limit",
											Computed:    true,
										},
										"project_count_limit": schema.Int64Attribute{
											Description: "Projects Count limit",
											Computed:    true,
										},
									},
								},
								"regional_limits": schema.ListNestedAttribute{
									Description: "Regions and their quota limits",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudQuotaRequestsRequestedLimitsRegionalLimitsDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"baremetal_basic_count_limit": schema.Int64Attribute{
												Description: "Basic bare metal servers count limit",
												Computed:    true,
											},
											"baremetal_gpu_a100_count_limit": schema.Int64Attribute{
												Description: "Baremetal A100 GPU card count limit",
												Computed:    true,
											},
											"baremetal_gpu_count_limit": schema.Int64Attribute{
												Description:        "Total number of AI GPU bare metal servers. This field is deprecated and is now always calculated automatically as the sum of `baremetal_gpu_a100_count_limit`, `baremetal_gpu_h100_count_limit`, `baremetal_gpu_h200_count_limit`, and `baremetal_gpu_l40s_count_limit`.",
												Computed:           true,
												DeprecationMessage: "This attribute is deprecated.",
											},
											"baremetal_gpu_h100_count_limit": schema.Int64Attribute{
												Description: "Baremetal H100 GPU card count limit",
												Computed:    true,
											},
											"baremetal_gpu_h200_count_limit": schema.Int64Attribute{
												Description: "Baremetal H200 GPU card count limit",
												Computed:    true,
											},
											"baremetal_gpu_l40s_count_limit": schema.Int64Attribute{
												Description: "Baremetal L40S GPU card count limit",
												Computed:    true,
											},
											"baremetal_hf_count_limit": schema.Int64Attribute{
												Description: "High-frequency bare metal servers count limit",
												Computed:    true,
											},
											"baremetal_infrastructure_count_limit": schema.Int64Attribute{
												Description: "Infrastructure bare metal servers count limit",
												Computed:    true,
											},
											"baremetal_network_count_limit": schema.Int64Attribute{
												Description: "Bare metal Network Count limit",
												Computed:    true,
											},
											"baremetal_storage_count_limit": schema.Int64Attribute{
												Description: "Storage bare metal servers count limit",
												Computed:    true,
											},
											"caas_container_count_limit": schema.Int64Attribute{
												Description: "Containers count limit",
												Computed:    true,
											},
											"caas_cpu_count_limit": schema.Int64Attribute{
												Description: "mCPU count for containers limit",
												Computed:    true,
											},
											"caas_gpu_count_limit": schema.Int64Attribute{
												Description: "Containers gpu count limit",
												Computed:    true,
											},
											"caas_ram_size_limit": schema.Int64Attribute{
												Description: "MB memory count for containers limit",
												Computed:    true,
											},
											"cluster_count_limit": schema.Int64Attribute{
												Description: "K8s clusters count limit",
												Computed:    true,
											},
											"cpu_count_limit": schema.Int64Attribute{
												Description: "vCPU Count limit",
												Computed:    true,
											},
											"dbaas_postgres_cluster_count_limit": schema.Int64Attribute{
												Description: "DBaaS cluster count limit",
												Computed:    true,
											},
											"external_ip_count_limit": schema.Int64Attribute{
												Description: "External IP Count limit",
												Computed:    true,
											},
											"faas_cpu_count_limit": schema.Int64Attribute{
												Description: "mCPU count for functions limit",
												Computed:    true,
											},
											"faas_function_count_limit": schema.Int64Attribute{
												Description: "Functions count limit",
												Computed:    true,
											},
											"faas_namespace_count_limit": schema.Int64Attribute{
												Description: "Functions namespace count limit",
												Computed:    true,
											},
											"faas_ram_size_limit": schema.Int64Attribute{
												Description: "MB memory count for functions limit",
												Computed:    true,
											},
											"firewall_count_limit": schema.Int64Attribute{
												Description: "Firewalls Count limit",
												Computed:    true,
											},
											"floating_count_limit": schema.Int64Attribute{
												Description: "Floating IP Count limit",
												Computed:    true,
											},
											"gpu_count_limit": schema.Int64Attribute{
												Description: "GPU Count limit",
												Computed:    true,
											},
											"gpu_virtual_a100_count_limit": schema.Int64Attribute{
												Description: "Virtual A100 GPU card count limit",
												Computed:    true,
											},
											"gpu_virtual_h100_count_limit": schema.Int64Attribute{
												Description: "Virtual H100 GPU card count limit",
												Computed:    true,
											},
											"gpu_virtual_h200_count_limit": schema.Int64Attribute{
												Description: "Virtual H200 GPU card count limit",
												Computed:    true,
											},
											"gpu_virtual_l40s_count_limit": schema.Int64Attribute{
												Description: "Virtual L40S GPU card count limit",
												Computed:    true,
											},
											"image_count_limit": schema.Int64Attribute{
												Description: "Images Count limit",
												Computed:    true,
											},
											"image_size_limit": schema.Int64Attribute{
												Description: "Images Size, GiB limit",
												Computed:    true,
											},
											"ipu_count_limit": schema.Int64Attribute{
												Description: "IPU Count limit",
												Computed:    true,
											},
											"laas_topic_count_limit": schema.Int64Attribute{
												Description: "LaaS Topics Count limit",
												Computed:    true,
											},
											"loadbalancer_count_limit": schema.Int64Attribute{
												Description: "Load Balancers Count limit",
												Computed:    true,
											},
											"network_count_limit": schema.Int64Attribute{
												Description: "Networks Count limit",
												Computed:    true,
											},
											"ram_limit": schema.Int64Attribute{
												Description: "RAM Size, GiB limit",
												Computed:    true,
											},
											"region_id": schema.Int64Attribute{
												Description: "Region ID",
												Computed:    true,
											},
											"registry_count_limit": schema.Int64Attribute{
												Description: "Registries count limit",
												Computed:    true,
											},
											"registry_storage_limit": schema.Int64Attribute{
												Description: "Registries volume usage, GiB limit",
												Computed:    true,
											},
											"router_count_limit": schema.Int64Attribute{
												Description: "Routers Count limit",
												Computed:    true,
											},
											"secret_count_limit": schema.Int64Attribute{
												Description: "Secret Count limit",
												Computed:    true,
											},
											"servergroup_count_limit": schema.Int64Attribute{
												Description: "Placement Group Count limit",
												Computed:    true,
											},
											"sfs_count_limit": schema.Int64Attribute{
												Description: "Shared file system Count limit",
												Computed:    true,
											},
											"sfs_size_limit": schema.Int64Attribute{
												Description: "Shared file system Size, GiB limit",
												Computed:    true,
											},
											"shared_vm_count_limit": schema.Int64Attribute{
												Description: "Basic VMs Count limit",
												Computed:    true,
											},
											"snapshot_schedule_count_limit": schema.Int64Attribute{
												Description: "Snapshot Schedules Count limit",
												Computed:    true,
											},
											"subnet_count_limit": schema.Int64Attribute{
												Description: "Subnets Count limit",
												Computed:    true,
											},
											"vm_count_limit": schema.Int64Attribute{
												Description: "Instances Dedicated Count limit",
												Computed:    true,
											},
											"volume_count_limit": schema.Int64Attribute{
												Description: "Volumes Count limit",
												Computed:    true,
											},
											"volume_size_limit": schema.Int64Attribute{
												Description: "Volumes Size, GiB limit",
												Computed:    true,
											},
											"volume_snapshots_count_limit": schema.Int64Attribute{
												Description: "Snapshots Count limit",
												Computed:    true,
											},
											"volume_snapshots_size_limit": schema.Int64Attribute{
												Description: "Snapshots Size, GiB limit",
												Computed:    true,
											},
										},
									},
								},
							},
						},
						"status": schema.StringAttribute{
							Description: "Request status",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Datetime when the request was created.",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"description": schema.StringAttribute{
							Description: "Describe the reason, in general terms.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Datetime when the request was updated.",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
					},
				},
			},
		},
	}
}

func (d *CloudQuotaRequestsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudQuotaRequestsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
