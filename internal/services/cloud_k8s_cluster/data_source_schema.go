// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_cluster

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

var _ datasource.DataSourceWithConfigValidators = (*CloudK8SClusterDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Managed Kubernetes clusters with configurable worker node pools, networking, and cluster add-ons.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Cluster name",
				Computed:    true,
			},
			"cluster_name": schema.StringAttribute{
				Description: "Cluster name",
				Required:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Function creation date",
				Computed:    true,
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"fixed_network": schema.StringAttribute{
				Description: "Fixed network id",
				Computed:    true,
			},
			"fixed_subnet": schema.StringAttribute{
				Description: "Fixed subnet id",
				Computed:    true,
			},
			"is_ipv6": schema.BoolAttribute{
				Description: "Enable public v6 address",
				Computed:    true,
			},
			"is_public": schema.BoolAttribute{
				Description: "Cluster is public",
				Computed:    true,
			},
			"keypair": schema.StringAttribute{
				Description: "Keypair",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name",
				Computed:    true,
			},
			"pods_ip_pool": schema.StringAttribute{
				Description: "The IP pool for the pods",
				Computed:    true,
			},
			"pods_ipv6_pool": schema.StringAttribute{
				Description: "The IPv6 pool for the pods",
				Computed:    true,
			},
			"services_ip_pool": schema.StringAttribute{
				Description: "The IP pool for the services",
				Computed:    true,
			},
			"services_ipv6_pool": schema.StringAttribute{
				Description: "The IPv6 pool for the services",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status\nAvailable values: \"Deleting\", \"Provisioned\", \"Provisioning\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"Deleting",
						"Provisioned",
						"Provisioning",
					),
				},
			},
			"version": schema.StringAttribute{
				Description: "K8s version",
				Computed:    true,
			},
			"autoscaler_config": schema.MapAttribute{
				Description: "Cluster autoscaler configuration.\n\nIt contains overrides to the default cluster-autoscaler parameters provided by the platform.",
				Computed:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"add_ons": schema.SingleNestedAttribute{
				Description: "Cluster add-ons configuration",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudK8SClusterAddOnsDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"slurm": schema.SingleNestedAttribute{
						Description: "Slurm add-on configuration",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudK8SClusterAddOnsSlurmDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Indicates whether Slurm add-on is deployed in the cluster.\n\nThis add-on is only supported in clusters running Kubernetes v1.31 and v1.32 with at least 1 GPU cluster pool.",
								Computed:    true,
							},
							"file_share_id": schema.StringAttribute{
								Description: "ID of a VAST file share used as Slurm storage.\n\nThe Slurm add-on creates separate Persistent Volume Claims for different purposes (controller spool, worker spool, jail) on that file share.",
								Computed:    true,
							},
							"ssh_key_ids": schema.ListAttribute{
								Description: "IDs of SSH keys authorized for SSH connection to Slurm login nodes.",
								Computed:    true,
								CustomType:  customfield.NewListType[types.String](ctx),
								ElementType: types.StringType,
							},
							"worker_count": schema.Int64Attribute{
								Description: "Size of the worker pool, i.e. number of worker nodes.\n\nEach Slurm worker node is backed by a Pod scheduled on one of cluster's GPU nodes.\n\nNote: Downscaling (reducing worker count) is not supported.",
								Computed:    true,
							},
						},
					},
				},
			},
			"authentication": schema.SingleNestedAttribute{
				Description: "Cluster authentication settings",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudK8SClusterAuthenticationDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"kubeconfig_created_at": schema.StringAttribute{
						Description: "Kubeconfig creation date",
						Computed:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"kubeconfig_expires_at": schema.StringAttribute{
						Description: "Kubeconfig expiration date",
						Computed:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"oidc": schema.SingleNestedAttribute{
						Description: "OIDC authentication settings",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudK8SClusterAuthenticationOidcDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"client_id": schema.StringAttribute{
								Description: "Client ID",
								Computed:    true,
							},
							"groups_claim": schema.StringAttribute{
								Description: "JWT claim to use as the user's group",
								Computed:    true,
							},
							"groups_prefix": schema.StringAttribute{
								Description: "Prefix prepended to group claims",
								Computed:    true,
							},
							"issuer_url": schema.StringAttribute{
								Description: "Issuer URL",
								Computed:    true,
							},
							"required_claims": schema.MapAttribute{
								Description: "Key-value pairs that describe required claims in the token",
								Computed:    true,
								CustomType:  customfield.NewMapType[types.String](ctx),
								ElementType: types.StringType,
							},
							"signing_algs": schema.ListAttribute{
								Description: "Accepted signing algorithms",
								Computed:    true,
								Validators: []validator.List{
									listvalidator.ValueStringsAre(
										stringvalidator.OneOfCaseInsensitive(
											"ES256",
											"ES384",
											"ES512",
											"PS256",
											"PS384",
											"PS512",
											"RS256",
											"RS384",
											"RS512",
										),
									),
								},
								CustomType:  customfield.NewListType[types.String](ctx),
								ElementType: types.StringType,
							},
							"username_claim": schema.StringAttribute{
								Description: "JWT claim to use as the user name",
								Computed:    true,
							},
							"username_prefix": schema.StringAttribute{
								Description: "Prefix prepended to username claims to prevent clashes",
								Computed:    true,
							},
						},
					},
				},
			},
			"cni": schema.SingleNestedAttribute{
				Description: "Cluster CNI settings",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudK8SClusterCniDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"cilium": schema.SingleNestedAttribute{
						Description: "Cilium settings",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudK8SClusterCniCiliumDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"encryption": schema.BoolAttribute{
								Description: "Wireguard encryption",
								Computed:    true,
							},
							"hubble_relay": schema.BoolAttribute{
								Description: "Hubble Relay",
								Computed:    true,
							},
							"hubble_ui": schema.BoolAttribute{
								Description: "Hubble UI",
								Computed:    true,
							},
							"lb_acceleration": schema.BoolAttribute{
								Description: "LoadBalancer acceleration",
								Computed:    true,
							},
							"lb_mode": schema.StringAttribute{
								Description: "LoadBalancer mode\nAvailable values: \"dsr\", \"hybrid\", \"snat\".",
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive(
										"dsr",
										"hybrid",
										"snat",
									),
								},
							},
							"mask_size": schema.Int64Attribute{
								Description: "Mask size for IPv4",
								Computed:    true,
							},
							"mask_size_v6": schema.Int64Attribute{
								Description: "Mask size for IPv6",
								Computed:    true,
							},
							"routing_mode": schema.StringAttribute{
								Description: "Routing mode\nAvailable values: \"native\", \"tunnel\".",
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("native", "tunnel"),
								},
							},
							"tunnel": schema.StringAttribute{
								Description: "CNI provider\nAvailable values: \"\", \"geneve\", \"vxlan\".",
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive(
										"",
										"geneve",
										"vxlan",
									),
								},
							},
						},
					},
					"cloud_k8s_cluster_provider": schema.StringAttribute{
						Description: "CNI provider\nAvailable values: \"calico\", \"cilium\".",
						Computed:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("calico", "cilium"),
						},
					},
				},
			},
			"csi": schema.SingleNestedAttribute{
				Description: "Cluster CSI settings",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudK8SClusterCsiDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"nfs": schema.SingleNestedAttribute{
						Description: "NFS settings",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudK8SClusterCsiNfsDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"vast_enabled": schema.BoolAttribute{
								Description: "Indicates the status of VAST NFS integration",
								Computed:    true,
							},
						},
					},
				},
			},
			"logging": schema.SingleNestedAttribute{
				Description: "Logging configuration",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudK8SClusterLoggingDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"destination_region_id": schema.Int64Attribute{
						Description: "ID of the region in which the logs will be stored",
						Computed:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "Indicates if log streaming is enabled or disabled",
						Computed:    true,
					},
					"topic_name": schema.StringAttribute{
						Description: "The topic name to stream logs to",
						Computed:    true,
					},
					"retention_policy": schema.SingleNestedAttribute{
						Description: "Logs retention policy",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudK8SClusterLoggingRetentionPolicyDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"period": schema.Int64Attribute{
								Description: "Duration of days for which logs must be kept.",
								Computed:    true,
								Validators: []validator.Int64{
									int64validator.AtMost(1825),
								},
							},
						},
					},
				},
			},
			"pools": schema.ListNestedAttribute{
				Description: "pools",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudK8SClusterPoolsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "UUID of the cluster pool",
							Computed:    true,
						},
						"auto_healing_enabled": schema.BoolAttribute{
							Description: "Indicates the status of auto healing",
							Computed:    true,
						},
						"boot_volume_size": schema.Int64Attribute{
							Description: "Size of the boot volume",
							Computed:    true,
						},
						"boot_volume_type": schema.StringAttribute{
							Description: "Type of the boot volume",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Date of function creation",
							Computed:    true,
						},
						"crio_config": schema.MapAttribute{
							Description: "Crio configuration for pool nodes",
							Computed:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"flavor_id": schema.StringAttribute{
							Description: "ID of the cluster pool flavor",
							Computed:    true,
						},
						"is_public_ipv4": schema.BoolAttribute{
							Description: "Indicates if the pool is public",
							Computed:    true,
						},
						"kubelet_config": schema.MapAttribute{
							Description: "Kubelet configuration for pool nodes",
							Computed:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"labels": schema.MapAttribute{
							Description: "Labels applied to the cluster pool",
							Computed:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"max_node_count": schema.Int64Attribute{
							Description: "Maximum node count in the cluster pool",
							Computed:    true,
						},
						"min_node_count": schema.Int64Attribute{
							Description: "Minimum node count in the cluster pool",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the cluster pool",
							Computed:    true,
						},
						"node_count": schema.Int64Attribute{
							Description: "Node count in the cluster pool",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Status of the cluster pool",
							Computed:    true,
						},
						"taints": schema.MapAttribute{
							Description: "Taints applied to the cluster pool",
							Computed:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"servergroup_id": schema.StringAttribute{
							Description: "Server group ID",
							Computed:    true,
						},
						"servergroup_name": schema.StringAttribute{
							Description: "Server group name",
							Computed:    true,
						},
						"servergroup_policy": schema.StringAttribute{
							Description: "Anti-affinity, affinity or soft-anti-affinity server group policy",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudK8SClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudK8SClusterDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
