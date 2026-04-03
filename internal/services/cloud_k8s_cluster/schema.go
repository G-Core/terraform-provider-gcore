// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_cluster

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudK8SClusterResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Managed Kubernetes clusters with configurable worker node pools, networking, and cluster add-ons.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "The name of the cluster",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "The name of the cluster",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
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
			"keypair": schema.StringAttribute{
				Description:   "The keypair of the cluster",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"version": schema.StringAttribute{
				Description:   "The version of the k8s cluster",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"pools": schema.ListNestedAttribute{
				Description: "The pools of the cluster",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"flavor_id": schema.StringAttribute{
							Description: "Flavor ID",
							Required:    true,
						},
						"min_node_count": schema.Int64Attribute{
							Description: "Minimum node count",
							Computed:    true,
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.Between(1, 200),
							},
						},
						"name": schema.StringAttribute{
							Description: "Pool's name",
							Required:    true,
						},
						"auto_healing_enabled": schema.BoolAttribute{
							Description: "Enable auto healing",
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
						},
						"boot_volume_size": schema.Int64Attribute{
							Description: "Boot volume size",
							Computed:    true,
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.Between(10, 2000),
							},
						},
						"boot_volume_type": schema.StringAttribute{
							Description: "Boot volume type\nAvailable values: \"cold\", \"ssd_hiiops\", \"ssd_local\", \"ssd_lowlatency\", \"standard\", \"ultra\".",
							Computed:    true,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"cold",
									"ssd_hiiops",
									"ssd_local",
									"ssd_lowlatency",
									"standard",
									"ultra",
								),
							},
						},
						"crio_config": schema.MapAttribute{
							Description: "Cri-o configuration for pool nodes",
							Computed:    true,
							Optional:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"is_public_ipv4": schema.BoolAttribute{
							Description: "Enable public v4 address",
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
						},
						"kubelet_config": schema.MapAttribute{
							Description: "Kubelet configuration for pool nodes",
							Computed:    true,
							Optional:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"labels": schema.MapAttribute{
							Description: "Labels applied to the cluster pool",
							Computed:    true,
							Optional:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"max_node_count": schema.Int64Attribute{
							Description: "Maximum node count",
							Computed:    true,
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.Between(1, 200),
							},
						},
						"servergroup_policy": schema.StringAttribute{
							Description: "Server group policy: anti-affinity, soft-anti-affinity or affinity\nAvailable values: \"affinity\", \"anti-affinity\", \"soft-anti-affinity\".",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"affinity",
									"anti-affinity",
									"soft-anti-affinity",
								),
							},
						},
						"taints": schema.MapAttribute{
							Description: "Taints applied to the cluster pool",
							Computed:    true,
							Optional:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"fixed_network": schema.StringAttribute{
				Description:   "The network of the cluster",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
				Default:       stringdefault.StaticString(""),
			},
			"fixed_subnet": schema.StringAttribute{
				Description:   "The subnet of the cluster",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
				Default:       stringdefault.StaticString(""),
			},
			"is_ipv6": schema.BoolAttribute{
				Description:   "Enable public v6 address",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplaceIfConfigured()},
				Default:       booldefault.StaticBool(false),
			},
			"pods_ip_pool": schema.StringAttribute{
				Description:   "The IP pool for the pods",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"pods_ipv6_pool": schema.StringAttribute{
				Description:   "The IPv6 pool for the pods",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"services_ip_pool": schema.StringAttribute{
				Description:   "The IP pool for the services",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"services_ipv6_pool": schema.StringAttribute{
				Description:   "The IPv6 pool for the services",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"csi": schema.SingleNestedAttribute{
				Description: "Container Storage Interface (CSI) driver settings",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudK8SClusterCsiModel](ctx),
				Attributes: map[string]schema.Attribute{
					"nfs": schema.SingleNestedAttribute{
						Description: "NFS CSI driver settings",
						Computed:    true,
						Optional:    true,
						CustomType:  customfield.NewNestedObjectType[CloudK8SClusterCsiNfsModel](ctx),
						Attributes: map[string]schema.Attribute{
							"vast_enabled": schema.BoolAttribute{
								Description: "Enable or disable VAST NFS integration. The default value is `false`. When set to `true`, a dedicated StorageClass will be created in the cluster for each VAST NFS file share defined in the cloud. All file shares created prior to cluster creation will be available immediately, while those created afterward may take a few minutes for to appear.",
								Computed:    true,
								Optional:    true,
								Default:     booldefault.StaticBool(false),
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplaceIfConfigured()},
			},
			"autoscaler_config": schema.MapAttribute{
				Description: "Cluster autoscaler configuration.\n\nIt allows you to override the default cluster-autoscaler parameters provided by the platform with your preferred values.\n\nSupported parameters (in alphabetical order):\n- balance-similar-node-groups (boolean: true/false) - Detect similar node groups and balance the number of nodes between them.\n- expander (string: random, most-pods, least-waste, price, priority, grpc) - Type of node group expander to be used in scale up. Specifying multiple values separated by commas will call the expanders in succession until there is only one option remaining.\n- expendable-pods-priority-cutoff (float) - Pods with priority below cutoff will be expendable. They can be killed without any consideration during scale down and they don't cause scale up. Pods with null priority (PodPriority disabled) are non expendable.\n- ignore-daemonsets-utilization (boolean: true/false) - Should CA ignore DaemonSet pods when calculating resource utilization for scaling down.\n- max-empty-bulk-delete (integer) - Maximum number of empty nodes that can be deleted at the same time.\n- max-graceful-termination-sec (integer) - Maximum number of seconds CA waits for pod termination when trying to scale down a node.\n- max-node-provision-time (duration: e.g., '15m') - The default maximum time CA waits for node to be provisioned - the value can be overridden per node group.\n- max-total-unready-percentage (float) - Maximum percentage of unready nodes in the cluster. After this is exceeded, CA halts operations.\n- new-pod-scale-up-delay (duration: e.g., '10s') - Pods less than this old will not be considered for scale-up. Can be increased for individual pods through annotation.\n- ok-total-unready-count (integer) - Number of allowed unready nodes, irrespective of max-total-unready-percentage.\n- scale-down-delay-after-add (duration: e.g., '10m') - How long after scale up that scale down evaluation resumes.\n- scale-down-delay-after-delete (duration: e.g., '10s') - How long after node deletion that scale down evaluation resumes.\n- scale-down-delay-after-failure (duration: e.g., '3m') - How long after scale down failure that scale down evaluation resumes.\n- scale-down-enabled (boolean: true/false) - Should CA scale down the cluster.\n- scale-down-unneeded-time (duration: e.g., '10m') - How long a node should be unneeded before it is eligible for scale down.\n- scale-down-unready-time (duration: e.g., '20m') - How long an unready node should be unneeded before it is eligible for scale down.\n- scale-down-utilization-threshold (float) - The maximum value between the sum of cpu requests and sum of memory requests of all pods running on the node divided by node's corresponding allocatable resource, below which a node can be considered for scale down.\n- scan-interval (duration: e.g., '10s') - How often cluster is reevaluated for scale up or down.\n- skip-nodes-with-custom-controller-pods (boolean: true/false) - If true cluster autoscaler will never delete nodes with pods owned by custom controllers.\n- skip-nodes-with-local-storage (boolean: true/false) - If true cluster autoscaler will never delete nodes with pods with local storage, e.g. EmptyDir or HostPath.\n- skip-nodes-with-system-pods (boolean: true/false) - If true cluster autoscaler will never delete nodes with pods from kube-system (except for DaemonSet or mirror pods).",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"add_ons": schema.SingleNestedAttribute{
				Description: "Cluster add-ons configuration",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudK8SClusterAddOnsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"slurm": schema.SingleNestedAttribute{
						Description: "Slurm add-on configuration",
						Computed:    true,
						Optional:    true,
						CustomType:  customfield.NewNestedObjectType[CloudK8SClusterAddOnsSlurmModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "The Slurm add-on will be enabled in the cluster.\n\nThis add-on is only supported in clusters running Kubernetes v1.31 and v1.32 with at least 1 GPU cluster pool and VAST NFS support enabled.",
								Computed:    true,
								Optional:    true,
							},
							"file_share_id": schema.StringAttribute{
								Description: "ID of a VAST file share to be used as Slurm storage.\n\nThe Slurm add-on will create separate Persistent Volume Claims for different purposes (controller spool, worker spool, jail) on that file share.\n\nThe file share must have `root_squash` disabled, while `path_length` and `allowed_characters` settings must be set to `NPL`.",
								Computed:    true,
								Optional:    true,
							},
							"ssh_key_ids": schema.ListAttribute{
								Description: "IDs of SSH keys to authorize for SSH connection to Slurm login nodes.",
								Computed:    true,
								Optional:    true,
								CustomType:  customfield.NewListType[types.String](ctx),
								ElementType: types.StringType,
							},
							"worker_count": schema.Int64Attribute{
								Description: "Size of the worker pool, i.e. the number of Slurm worker nodes.\n\nEach Slurm worker node will be backed by a Pod scheduled on one of cluster's GPU nodes.\n\nNote: Downscaling (reducing worker count) is not supported.",
								Computed:    true,
								Optional:    true,
								Validators: []validator.Int64{
									int64validator.AtLeast(1),
								},
							},
						},
					},
				},
			},
			"authentication": schema.SingleNestedAttribute{
				Description: "Authentication settings",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudK8SClusterAuthenticationModel](ctx),
				Attributes: map[string]schema.Attribute{
					"oidc": schema.SingleNestedAttribute{
						Description: "OIDC authentication settings",
						Computed:    true,
						Optional:    true,
						CustomType:  customfield.NewNestedObjectType[CloudK8SClusterAuthenticationOidcModel](ctx),
						Attributes: map[string]schema.Attribute{
							"client_id": schema.StringAttribute{
								Description: "Client ID",
								Computed:    true,
								Optional:    true,
							},
							"groups_claim": schema.StringAttribute{
								Description: "JWT claim to use as the user's group",
								Computed:    true,
								Optional:    true,
							},
							"groups_prefix": schema.StringAttribute{
								Description: "Prefix prepended to group claims",
								Computed:    true,
								Optional:    true,
							},
							"issuer_url": schema.StringAttribute{
								Description: "Issuer URL",
								Computed:    true,
								Optional:    true,
							},
							"required_claims": schema.MapAttribute{
								Description: "Key-value pairs that describe required claims in the token",
								Computed:    true,
								Optional:    true,
								CustomType:  customfield.NewMapType[types.String](ctx),
								ElementType: types.StringType,
							},
							"signing_algs": schema.ListAttribute{
								Description: "Accepted signing algorithms",
								Computed:    true,
								Optional:    true,
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
								Optional:    true,
							},
							"username_prefix": schema.StringAttribute{
								Description: "Prefix prepended to username claims to prevent clashes",
								Computed:    true,
								Optional:    true,
							},
						},
					},
				},
			},
			"cni": schema.SingleNestedAttribute{
				Description: "Cluster CNI settings",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudK8SClusterCniModel](ctx),
				Attributes: map[string]schema.Attribute{
					"cilium": schema.SingleNestedAttribute{
						Description: "Cilium settings",
						Computed:    true,
						Optional:    true,
						CustomType:  customfield.NewNestedObjectType[CloudK8SClusterCniCiliumModel](ctx),
						Attributes: map[string]schema.Attribute{
							"encryption": schema.BoolAttribute{
								Description: "Wireguard encryption",
								Computed:    true,
								Optional:    true,
								Default:     booldefault.StaticBool(false),
							},
							"hubble_relay": schema.BoolAttribute{
								Description: "Hubble Relay",
								Computed:    true,
								Optional:    true,
								Default:     booldefault.StaticBool(false),
							},
							"hubble_ui": schema.BoolAttribute{
								Description: "Hubble UI",
								Computed:    true,
								Optional:    true,
								Default:     booldefault.StaticBool(false),
							},
							"lb_acceleration": schema.BoolAttribute{
								Description: "LoadBalancer acceleration",
								Computed:    true,
								Optional:    true,
								Default:     booldefault.StaticBool(false),
							},
							"lb_mode": schema.StringAttribute{
								Description: "LoadBalancer mode\nAvailable values: \"dsr\", \"hybrid\", \"snat\".",
								Computed:    true,
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive(
										"dsr",
										"hybrid",
										"snat",
									),
								},
								Default: stringdefault.StaticString("snat"),
							},
							"mask_size": schema.Int64Attribute{
								Description: "Mask size for IPv4",
								Computed:    true,
								Optional:    true,
								Default:     int64default.StaticInt64(24),
							},
							"mask_size_v6": schema.Int64Attribute{
								Description: "Mask size for IPv6",
								Computed:    true,
								Optional:    true,
								Default:     int64default.StaticInt64(120),
							},
							"routing_mode": schema.StringAttribute{
								Description: "Routing mode\nAvailable values: \"native\", \"tunnel\".",
								Computed:    true,
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("native", "tunnel"),
								},
								Default: stringdefault.StaticString("tunnel"),
							},
							"tunnel": schema.StringAttribute{
								Description: "CNI provider\nAvailable values: \"\", \"geneve\", \"vxlan\".",
								Computed:    true,
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive(
										"",
										"geneve",
										"vxlan",
									),
								},
								Default: stringdefault.StaticString("geneve"),
							},
						},
					},
					"cloud_k8s_cluster_provider": schema.StringAttribute{
						Description: "CNI provider\nAvailable values: \"calico\", \"cilium\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("calico", "cilium"),
						},
						Default: stringdefault.StaticString("calico"),
					},
				},
			},
			"ddos_profile": schema.SingleNestedAttribute{
				Description: "Advanced DDoS Protection profile",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudK8SClusterDDOSProfileModel](ctx),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Enable advanced DDoS protection",
						Required:    true,
					},
					"fields": schema.ListNestedAttribute{
						Description: "DDoS profile parameters",
						Computed:    true,
						Optional:    true,
						CustomType:  customfield.NewNestedObjectListType[CloudK8SClusterDDOSProfileFieldsModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"base_field": schema.Int64Attribute{
									Required: true,
								},
								"field_value": schema.StringAttribute{
									Description: "Complex value. Only one of 'value' or 'field_value' must be specified",
									Optional:    true,
									CustomType:  jsontypes.NormalizedType{},
								},
								"value": schema.StringAttribute{
									Description: "Basic value. Only one of 'value' or 'field_value' must be specified",
									Optional:    true,
								},
							},
						},
					},
					"profile_template": schema.Int64Attribute{
						Description: "DDoS profile template ID",
						Optional:    true,
					},
					"profile_template_name": schema.StringAttribute{
						Description: "DDoS profile template name",
						Optional:    true,
					},
				},
			},
			"logging": schema.SingleNestedAttribute{
				Description: "Logging configuration",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudK8SClusterLoggingModel](ctx),
				Attributes: map[string]schema.Attribute{
					"destination_region_id": schema.Int64Attribute{
						Description: "Destination region id to which the logs will be written",
						Optional:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "Enable/disable forwarding logs to LaaS",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"retention_policy": schema.SingleNestedAttribute{
						Description: "The logs retention policy",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"period": schema.Int64Attribute{
								Description: "Duration of days for which logs must be kept.",
								Required:    true,
								Validators: []validator.Int64{
									int64validator.AtMost(1825),
								},
							},
						},
					},
					"topic_name": schema.StringAttribute{
						Description: "The topic name to which the logs will be written",
						Optional:    true,
					},
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Function creation date",
				Computed:    true,
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"is_public": schema.BoolAttribute{
				Description: "Cluster is public",
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
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"tasks": schema.ListAttribute{
				Description: "List of task IDs representing asynchronous operations. Use these IDs to monitor operation progress:\n- `GET /v1/tasks/{task_id}` - Check individual task status and details\nPoll task status until completion (`FINISHED`/`ERROR`) before proceeding with dependent operations.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
		},
	}
}

func (r *CloudK8SClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudK8SClusterResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
