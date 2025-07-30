// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_quota_request

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*CloudQuotaRequestResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"request_id": schema.StringAttribute{
				Description:   "LimitRequest ID",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"description": schema.StringAttribute{
				Description:   "Describe the reason, in general terms.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"requested_limits": schema.SingleNestedAttribute{
				Description: "Limits you want to increase.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"global_limits": schema.SingleNestedAttribute{
						Description: "Global entity quota limits",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"inference_cpu_millicore_count_limit": schema.Int64Attribute{
								Description: "Inference CPU millicore count limit",
								Optional:    true,
							},
							"inference_gpu_a100_count_limit": schema.Int64Attribute{
								Description: "Inference GPU A100 Count limit",
								Optional:    true,
							},
							"inference_gpu_h100_count_limit": schema.Int64Attribute{
								Description: "Inference GPU H100 Count limit",
								Optional:    true,
							},
							"inference_gpu_l40s_count_limit": schema.Int64Attribute{
								Description: "Inference GPU L40s Count limit",
								Optional:    true,
							},
							"inference_instance_count_limit": schema.Int64Attribute{
								Description: "Inference instance count limit",
								Optional:    true,
							},
							"keypair_count_limit": schema.Int64Attribute{
								Description: "SSH Keys Count limit",
								Optional:    true,
							},
							"project_count_limit": schema.Int64Attribute{
								Description: "Projects Count limit",
								Optional:    true,
							},
						},
					},
					"regional_limits": schema.ListNestedAttribute{
						Description: "Regions and their quota limits",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"baremetal_basic_count_limit": schema.Int64Attribute{
									Description: "Basic bare metal servers count limit",
									Optional:    true,
								},
								"baremetal_gpu_a100_count_limit": schema.Int64Attribute{
									Description: "Baremetal A100 GPU card count limit",
									Optional:    true,
								},
								"baremetal_gpu_count_limit": schema.Int64Attribute{
									Description:        "Total number of AI GPU bare metal servers. This field is deprecated and is now always calculated automatically as the sum of `baremetal_gpu_a100_count_limit`, `baremetal_gpu_h100_count_limit`, `baremetal_gpu_h200_count_limit`, and `baremetal_gpu_l40s_count_limit`.",
									Optional:           true,
									DeprecationMessage: "This attribute is deprecated.",
								},
								"baremetal_gpu_h100_count_limit": schema.Int64Attribute{
									Description: "Baremetal H100 GPU card count limit",
									Optional:    true,
								},
								"baremetal_gpu_h200_count_limit": schema.Int64Attribute{
									Description: "Baremetal H200 GPU card count limit",
									Optional:    true,
								},
								"baremetal_gpu_l40s_count_limit": schema.Int64Attribute{
									Description: "Baremetal L40S GPU card count limit",
									Optional:    true,
								},
								"baremetal_hf_count_limit": schema.Int64Attribute{
									Description: "High-frequency bare metal servers count limit",
									Optional:    true,
								},
								"baremetal_infrastructure_count_limit": schema.Int64Attribute{
									Description: "Infrastructure bare metal servers count limit",
									Optional:    true,
								},
								"baremetal_network_count_limit": schema.Int64Attribute{
									Description: "Bare metal Network Count limit",
									Optional:    true,
								},
								"baremetal_storage_count_limit": schema.Int64Attribute{
									Description: "Storage bare metal servers count limit",
									Optional:    true,
								},
								"caas_container_count_limit": schema.Int64Attribute{
									Description: "Containers count limit",
									Optional:    true,
								},
								"caas_cpu_count_limit": schema.Int64Attribute{
									Description: "mCPU count for containers limit",
									Optional:    true,
								},
								"caas_gpu_count_limit": schema.Int64Attribute{
									Description: "Containers gpu count limit",
									Optional:    true,
								},
								"caas_ram_size_limit": schema.Int64Attribute{
									Description: "MB memory count for containers limit",
									Optional:    true,
								},
								"cluster_count_limit": schema.Int64Attribute{
									Description: "K8s clusters count limit",
									Optional:    true,
								},
								"cpu_count_limit": schema.Int64Attribute{
									Description: "vCPU Count limit",
									Optional:    true,
								},
								"dbaas_postgres_cluster_count_limit": schema.Int64Attribute{
									Description: "DBaaS cluster count limit",
									Optional:    true,
								},
								"external_ip_count_limit": schema.Int64Attribute{
									Description: "External IP Count limit",
									Optional:    true,
								},
								"faas_cpu_count_limit": schema.Int64Attribute{
									Description: "mCPU count for functions limit",
									Optional:    true,
								},
								"faas_function_count_limit": schema.Int64Attribute{
									Description: "Functions count limit",
									Optional:    true,
								},
								"faas_namespace_count_limit": schema.Int64Attribute{
									Description: "Functions namespace count limit",
									Optional:    true,
								},
								"faas_ram_size_limit": schema.Int64Attribute{
									Description: "MB memory count for functions limit",
									Optional:    true,
								},
								"firewall_count_limit": schema.Int64Attribute{
									Description: "Firewalls Count limit",
									Optional:    true,
								},
								"floating_count_limit": schema.Int64Attribute{
									Description: "Floating IP Count limit",
									Optional:    true,
								},
								"gpu_count_limit": schema.Int64Attribute{
									Description: "GPU Count limit",
									Optional:    true,
								},
								"gpu_virtual_a100_count_limit": schema.Int64Attribute{
									Description: "Virtual A100 GPU card count limit",
									Optional:    true,
								},
								"gpu_virtual_h100_count_limit": schema.Int64Attribute{
									Description: "Virtual H100 GPU card count limit",
									Optional:    true,
								},
								"gpu_virtual_h200_count_limit": schema.Int64Attribute{
									Description: "Virtual H200 GPU card count limit",
									Optional:    true,
								},
								"gpu_virtual_l40s_count_limit": schema.Int64Attribute{
									Description: "Virtual L40S GPU card count limit",
									Optional:    true,
								},
								"image_count_limit": schema.Int64Attribute{
									Description: "Images Count limit",
									Optional:    true,
								},
								"image_size_limit": schema.Int64Attribute{
									Description: "Images Size, GiB limit",
									Optional:    true,
								},
								"ipu_count_limit": schema.Int64Attribute{
									Description: "IPU Count limit",
									Optional:    true,
								},
								"laas_topic_count_limit": schema.Int64Attribute{
									Description: "LaaS Topics Count limit",
									Optional:    true,
								},
								"loadbalancer_count_limit": schema.Int64Attribute{
									Description: "Load Balancers Count limit",
									Optional:    true,
								},
								"network_count_limit": schema.Int64Attribute{
									Description: "Networks Count limit",
									Optional:    true,
								},
								"ram_limit": schema.Int64Attribute{
									Description: "RAM Size, GiB limit",
									Optional:    true,
								},
								"region_id": schema.Int64Attribute{
									Description: "Region ID",
									Optional:    true,
								},
								"registry_count_limit": schema.Int64Attribute{
									Description: "Registries count limit",
									Optional:    true,
								},
								"registry_storage_limit": schema.Int64Attribute{
									Description: "Registries volume usage, GiB limit",
									Optional:    true,
								},
								"router_count_limit": schema.Int64Attribute{
									Description: "Routers Count limit",
									Optional:    true,
								},
								"secret_count_limit": schema.Int64Attribute{
									Description: "Secret Count limit",
									Optional:    true,
								},
								"servergroup_count_limit": schema.Int64Attribute{
									Description: "Placement Group Count limit",
									Optional:    true,
								},
								"sfs_count_limit": schema.Int64Attribute{
									Description: "Shared file system Count limit",
									Optional:    true,
								},
								"sfs_size_limit": schema.Int64Attribute{
									Description: "Shared file system Size, GiB limit",
									Optional:    true,
								},
								"shared_vm_count_limit": schema.Int64Attribute{
									Description: "Basic VMs Count limit",
									Optional:    true,
								},
								"snapshot_schedule_count_limit": schema.Int64Attribute{
									Description: "Snapshot Schedules Count limit",
									Optional:    true,
								},
								"subnet_count_limit": schema.Int64Attribute{
									Description: "Subnets Count limit",
									Optional:    true,
								},
								"vm_count_limit": schema.Int64Attribute{
									Description: "Instances Dedicated Count limit",
									Optional:    true,
								},
								"volume_count_limit": schema.Int64Attribute{
									Description: "Volumes Count limit",
									Optional:    true,
								},
								"volume_size_limit": schema.Int64Attribute{
									Description: "Volumes Size, GiB limit",
									Optional:    true,
								},
								"volume_snapshots_count_limit": schema.Int64Attribute{
									Description: "Snapshots Count limit",
									Optional:    true,
								},
								"volume_snapshots_size_limit": schema.Int64Attribute{
									Description: "Snapshots Size, GiB limit",
									Optional:    true,
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
			},
			"client_id": schema.Int64Attribute{
				Description:   "Client ID that requests the limit increase.",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the request was created.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"id": schema.Int64Attribute{
				Description: "Request ID",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Request status",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the request was updated.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *CloudQuotaRequestResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudQuotaRequestResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
