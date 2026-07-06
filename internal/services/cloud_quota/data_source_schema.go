// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_quota

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudQuotaDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Quotas define the maximum amount of cloud resources (compute, storage, networking, GPU, and more) available to a client, both globally and per region.",
		Attributes: map[string]schema.Attribute{
			"global_quotas": schema.SingleNestedAttribute{
				Description: "Global entity quotas",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudQuotaGlobalQuotasDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"inference_cpu_millicore_count_limit": schema.Int64Attribute{
						Description: "Inference CPU millicore count limit",
						Computed:    true,
					},
					"inference_cpu_millicore_count_usage": schema.Int64Attribute{
						Description: "Inference CPU millicore count usage",
						Computed:    true,
					},
					"inference_gpu_a100_count_limit": schema.Int64Attribute{
						Description: "Inference GPU A100 Count limit",
						Computed:    true,
					},
					"inference_gpu_a100_count_usage": schema.Int64Attribute{
						Description: "Inference GPU A100 Count usage",
						Computed:    true,
					},
					"inference_gpu_h100_count_limit": schema.Int64Attribute{
						Description: "Inference GPU H100 Count limit",
						Computed:    true,
					},
					"inference_gpu_h100_count_usage": schema.Int64Attribute{
						Description: "Inference GPU H100 Count usage",
						Computed:    true,
					},
					"inference_gpu_l40s_count_limit": schema.Int64Attribute{
						Description: "Inference GPU L40s Count limit",
						Computed:    true,
					},
					"inference_gpu_l40s_count_usage": schema.Int64Attribute{
						Description: "Inference GPU L40s Count usage",
						Computed:    true,
					},
					"inference_instance_count_limit": schema.Int64Attribute{
						Description: "Inference instance count limit",
						Computed:    true,
					},
					"inference_instance_count_usage": schema.Int64Attribute{
						Description: "Inference instance count usage",
						Computed:    true,
					},
					"inference_public_model_api_key_count_limit": schema.Int64Attribute{
						Description: "Public model API keys count limit",
						Computed:    true,
					},
					"inference_public_model_api_key_count_usage": schema.Int64Attribute{
						Description: "Public model API keys count usage",
						Computed:    true,
					},
					"keypair_count_limit": schema.Int64Attribute{
						Description: "SSH Keys Count limit",
						Computed:    true,
					},
					"keypair_count_usage": schema.Int64Attribute{
						Description: "SSH Keys Count usage",
						Computed:    true,
					},
					"project_count_limit": schema.Int64Attribute{
						Description: "Projects Count limit",
						Computed:    true,
					},
					"project_count_usage": schema.Int64Attribute{
						Description: "Projects Count usage",
						Computed:    true,
					},
				},
			},
			"regional_quotas": schema.ListNestedAttribute{
				Description: "Regional entity quotas. Only contains initialized quotas.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudQuotaRegionalQuotasDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"baremetal_basic_count_limit": schema.Int64Attribute{
							Description: "Basic bare metal servers count limit",
							Computed:    true,
						},
						"baremetal_basic_count_usage": schema.Int64Attribute{
							Description: "Basic bare metal servers count usage",
							Computed:    true,
						},
						"baremetal_gpu_a100_count_limit": schema.Int64Attribute{
							Description: "Bare metal A100 GPU server count limit",
							Computed:    true,
						},
						"baremetal_gpu_a100_count_usage": schema.Int64Attribute{
							Description: "Bare metal A100 GPU server count usage",
							Computed:    true,
						},
						"baremetal_gpu_count_limit": schema.Int64Attribute{
							Description:        "Total number of AI GPU bare metal servers. This field is deprecated and is now always calculated automatically as the sum of `baremetal_gpu_a100_count_limit`, `baremetal_gpu_h100_count_limit`, `baremetal_gpu_h200_count_limit`, and `baremetal_gpu_l40s_count_limit`.",
							Computed:           true,
							DeprecationMessage: "This attribute is deprecated.",
						},
						"baremetal_gpu_count_usage": schema.Int64Attribute{
							Description:        "Baremetal Gpu Count Usage. This field is deprecated and is now always calculated automatically as the sum of `baremetal_gpu_a100_count_usage`, `baremetal_gpu_h100_count_usage`, `baremetal_gpu_h200_count_usage`, and `baremetal_gpu_l40s_count_usage`.",
							Computed:           true,
							DeprecationMessage: "This attribute is deprecated.",
						},
						"baremetal_gpu_h100_count_limit": schema.Int64Attribute{
							Description: "Bare metal H100 GPU server count limit",
							Computed:    true,
						},
						"baremetal_gpu_h100_count_usage": schema.Int64Attribute{
							Description: "Bare metal H100 GPU server count usage",
							Computed:    true,
						},
						"baremetal_gpu_h200_count_limit": schema.Int64Attribute{
							Description: "Bare metal H200 GPU server count limit",
							Computed:    true,
						},
						"baremetal_gpu_h200_count_usage": schema.Int64Attribute{
							Description: "Bare metal H200 GPU server count usage",
							Computed:    true,
						},
						"baremetal_gpu_l40s_count_limit": schema.Int64Attribute{
							Description: "Bare metal L40S GPU server count limit",
							Computed:    true,
						},
						"baremetal_gpu_l40s_count_usage": schema.Int64Attribute{
							Description: "Bare metal L40S GPU server count usage",
							Computed:    true,
						},
						"baremetal_hf_count_limit": schema.Int64Attribute{
							Description: "High-frequency bare metal servers count limit",
							Computed:    true,
						},
						"baremetal_hf_count_usage": schema.Int64Attribute{
							Description: "High-frequency bare metal servers count usage",
							Computed:    true,
						},
						"baremetal_infrastructure_count_limit": schema.Int64Attribute{
							Description: "Infrastructure bare metal servers count limit",
							Computed:    true,
						},
						"baremetal_infrastructure_count_usage": schema.Int64Attribute{
							Description: "Infrastructure bare metal servers count usage",
							Computed:    true,
						},
						"baremetal_network_count_limit": schema.Int64Attribute{
							Description: "Bare metal Network Count limit",
							Computed:    true,
						},
						"baremetal_network_count_usage": schema.Int64Attribute{
							Description: "Bare metal Network Count usage",
							Computed:    true,
						},
						"baremetal_storage_count_limit": schema.Int64Attribute{
							Description: "Storage bare metal servers count limit",
							Computed:    true,
						},
						"baremetal_storage_count_usage": schema.Int64Attribute{
							Description: "Storage bare metal servers count usage",
							Computed:    true,
						},
						"caas_container_count_limit": schema.Int64Attribute{
							Description: "Containers count limit",
							Computed:    true,
						},
						"caas_container_count_usage": schema.Int64Attribute{
							Description: "Containers count usage",
							Computed:    true,
						},
						"caas_cpu_count_limit": schema.Int64Attribute{
							Description: "mCPU count for containers limit",
							Computed:    true,
						},
						"caas_cpu_count_usage": schema.Int64Attribute{
							Description: "mCPU count for containers usage",
							Computed:    true,
						},
						"caas_gpu_count_limit": schema.Int64Attribute{
							Description: "Containers gpu count limit",
							Computed:    true,
						},
						"caas_gpu_count_usage": schema.Int64Attribute{
							Description: "Containers gpu count usage",
							Computed:    true,
						},
						"caas_ram_size_limit": schema.Int64Attribute{
							Description: "MiB memory count for containers limit",
							Computed:    true,
						},
						"caas_ram_size_usage": schema.Int64Attribute{
							Description: "MiB memory count for containers usage",
							Computed:    true,
						},
						"cluster_count_limit": schema.Int64Attribute{
							Description: "K8s clusters count limit",
							Computed:    true,
						},
						"cluster_count_usage": schema.Int64Attribute{
							Description: "K8s clusters count usage",
							Computed:    true,
						},
						"cpu_count_limit": schema.Int64Attribute{
							Description: "vCPU Count limit",
							Computed:    true,
						},
						"cpu_count_usage": schema.Int64Attribute{
							Description: "vCPU Count usage",
							Computed:    true,
						},
						"dbaas_postgres_cluster_count_limit": schema.Int64Attribute{
							Description: "DBaaS cluster count limit",
							Computed:    true,
						},
						"dbaas_postgres_cluster_count_usage": schema.Int64Attribute{
							Description: "DBaaS cluster count usage",
							Computed:    true,
						},
						"external_ip_count_limit": schema.Int64Attribute{
							Description: "External IP Count limit",
							Computed:    true,
						},
						"external_ip_count_usage": schema.Int64Attribute{
							Description: "External IP Count usage",
							Computed:    true,
						},
						"faas_cpu_count_limit": schema.Int64Attribute{
							Description: "mCPU count for functions limit",
							Computed:    true,
						},
						"faas_cpu_count_usage": schema.Int64Attribute{
							Description: "mCPU count for functions usage",
							Computed:    true,
						},
						"faas_function_count_limit": schema.Int64Attribute{
							Description: "Functions count limit",
							Computed:    true,
						},
						"faas_function_count_usage": schema.Int64Attribute{
							Description: "Functions count usage",
							Computed:    true,
						},
						"faas_namespace_count_limit": schema.Int64Attribute{
							Description: "Functions namespace count limit",
							Computed:    true,
						},
						"faas_namespace_count_usage": schema.Int64Attribute{
							Description: "Functions namespace count usage",
							Computed:    true,
						},
						"faas_ram_size_limit": schema.Int64Attribute{
							Description: "MiB memory count for functions limit",
							Computed:    true,
						},
						"faas_ram_size_usage": schema.Int64Attribute{
							Description: "MiB memory count for functions usage",
							Computed:    true,
						},
						"firewall_count_limit": schema.Int64Attribute{
							Description: "Firewalls Count limit",
							Computed:    true,
						},
						"firewall_count_usage": schema.Int64Attribute{
							Description: "Firewalls Count usage",
							Computed:    true,
						},
						"floating_count_limit": schema.Int64Attribute{
							Description: "Floating IP Count limit",
							Computed:    true,
						},
						"floating_count_usage": schema.Int64Attribute{
							Description: "Floating IP Count usage",
							Computed:    true,
						},
						"gpu_count_limit": schema.Int64Attribute{
							Description: "GPU Count limit",
							Computed:    true,
						},
						"gpu_count_usage": schema.Int64Attribute{
							Description: "GPU Count usage",
							Computed:    true,
						},
						"gpu_virtual_a100_count_limit": schema.Int64Attribute{
							Description: "Virtual A100 GPU card count limit",
							Computed:    true,
						},
						"gpu_virtual_a100_count_usage": schema.Int64Attribute{
							Description: "Virtual A100 GPU card count usage",
							Computed:    true,
						},
						"gpu_virtual_h100_count_limit": schema.Int64Attribute{
							Description: "Virtual H100 GPU card count limit",
							Computed:    true,
						},
						"gpu_virtual_h100_count_usage": schema.Int64Attribute{
							Description: "Virtual H100 GPU card count usage",
							Computed:    true,
						},
						"gpu_virtual_h200_count_limit": schema.Int64Attribute{
							Description: "Virtual H200 GPU card count limit",
							Computed:    true,
						},
						"gpu_virtual_h200_count_usage": schema.Int64Attribute{
							Description: "Virtual H200 GPU card count usage",
							Computed:    true,
						},
						"gpu_virtual_l40s_count_limit": schema.Int64Attribute{
							Description: "Virtual L40S GPU card count limit",
							Computed:    true,
						},
						"gpu_virtual_l40s_count_usage": schema.Int64Attribute{
							Description: "Virtual L40S GPU card count usage",
							Computed:    true,
						},
						"image_count_limit": schema.Int64Attribute{
							Description: "Images Count limit",
							Computed:    true,
						},
						"image_count_usage": schema.Int64Attribute{
							Description: "Images Count usage",
							Computed:    true,
						},
						"image_size_limit": schema.Int64Attribute{
							Description: "Images Size, bytes limit",
							Computed:    true,
						},
						"image_size_usage": schema.Int64Attribute{
							Description: "Images Size, bytes usage",
							Computed:    true,
						},
						"ipu_count_limit": schema.Int64Attribute{
							Description: "IPU Count limit",
							Computed:    true,
						},
						"ipu_count_usage": schema.Int64Attribute{
							Description: "IPU Count usage",
							Computed:    true,
						},
						"laas_topic_count_limit": schema.Int64Attribute{
							Description: "LaaS Topics Count limit",
							Computed:    true,
						},
						"laas_topic_count_usage": schema.Int64Attribute{
							Description: "LaaS Topics Count usage",
							Computed:    true,
						},
						"loadbalancer_count_limit": schema.Int64Attribute{
							Description: "Load Balancers Count limit",
							Computed:    true,
						},
						"loadbalancer_count_usage": schema.Int64Attribute{
							Description: "Load Balancers Count usage",
							Computed:    true,
						},
						"network_count_limit": schema.Int64Attribute{
							Description: "Networks Count limit",
							Computed:    true,
						},
						"network_count_usage": schema.Int64Attribute{
							Description: "Networks Count usage",
							Computed:    true,
						},
						"ram_limit": schema.Int64Attribute{
							Description: "RAM Size, MiB limit",
							Computed:    true,
						},
						"ram_usage": schema.Int64Attribute{
							Description: "RAM Size, MiB usage",
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
						"registry_count_usage": schema.Int64Attribute{
							Description: "Registries count usage",
							Computed:    true,
						},
						"registry_storage_limit": schema.Int64Attribute{
							Description: "Registries volume usage, GiB limit",
							Computed:    true,
						},
						"registry_storage_usage": schema.Int64Attribute{
							Description: "Registries volume usage, GiB usage",
							Computed:    true,
						},
						"router_count_limit": schema.Int64Attribute{
							Description: "Routers Count limit",
							Computed:    true,
						},
						"router_count_usage": schema.Int64Attribute{
							Description: "Routers Count usage",
							Computed:    true,
						},
						"secret_count_limit": schema.Int64Attribute{
							Description: "Secret Count limit",
							Computed:    true,
						},
						"secret_count_usage": schema.Int64Attribute{
							Description: "Secret Count usage",
							Computed:    true,
						},
						"servergroup_count_limit": schema.Int64Attribute{
							Description: "Placement Group Count limit",
							Computed:    true,
						},
						"servergroup_count_usage": schema.Int64Attribute{
							Description: "Placement Group Count usage",
							Computed:    true,
						},
						"sfs_count_limit": schema.Int64Attribute{
							Description: "Shared file system Count limit",
							Computed:    true,
						},
						"sfs_count_usage": schema.Int64Attribute{
							Description: "Shared file system Count usage",
							Computed:    true,
						},
						"sfs_size_limit": schema.Int64Attribute{
							Description: "Shared file system Size, GiB limit",
							Computed:    true,
						},
						"sfs_size_usage": schema.Int64Attribute{
							Description: "Shared file system Size, GiB usage",
							Computed:    true,
						},
						"shared_vm_count_limit": schema.Int64Attribute{
							Description: "Basic VMs Count limit",
							Computed:    true,
						},
						"shared_vm_count_usage": schema.Int64Attribute{
							Description: "Basic VMs Count usage",
							Computed:    true,
						},
						"snapshot_schedule_count_limit": schema.Int64Attribute{
							Description: "Snapshot Schedules Count limit",
							Computed:    true,
						},
						"snapshot_schedule_count_usage": schema.Int64Attribute{
							Description: "Snapshot Schedules Count usage",
							Computed:    true,
						},
						"subnet_count_limit": schema.Int64Attribute{
							Description: "Subnets Count limit",
							Computed:    true,
						},
						"subnet_count_usage": schema.Int64Attribute{
							Description: "Subnets Count usage",
							Computed:    true,
						},
						"vm_count_limit": schema.Int64Attribute{
							Description: "Instances Dedicated Count limit",
							Computed:    true,
						},
						"vm_count_usage": schema.Int64Attribute{
							Description: "Instances Dedicated Count usage",
							Computed:    true,
						},
						"volume_count_limit": schema.Int64Attribute{
							Description: "Volumes Count limit",
							Computed:    true,
						},
						"volume_count_usage": schema.Int64Attribute{
							Description: "Volumes Count usage",
							Computed:    true,
						},
						"volume_size_limit": schema.Int64Attribute{
							Description: "Volumes Size, GiB limit",
							Computed:    true,
						},
						"volume_size_usage": schema.Int64Attribute{
							Description: "Volumes Size, GiB usage",
							Computed:    true,
						},
						"volume_snapshots_count_limit": schema.Int64Attribute{
							Description: "Snapshots Count limit",
							Computed:    true,
						},
						"volume_snapshots_count_usage": schema.Int64Attribute{
							Description: "Snapshots Count usage",
							Computed:    true,
						},
						"volume_snapshots_size_limit": schema.Int64Attribute{
							Description: "Snapshots Size, GiB limit",
							Computed:    true,
						},
						"volume_snapshots_size_usage": schema.Int64Attribute{
							Description: "Snapshots Size, GiB usage",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudQuotaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudQuotaDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
