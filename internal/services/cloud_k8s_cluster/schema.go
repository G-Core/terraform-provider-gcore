// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_cluster

import (
	"context"
	"strings"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
			},
			"region_id": schema.Int64Attribute{
				Description:   "Region ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
			},
			"keypair": schema.StringAttribute{
				Description:   "The keypair of the cluster",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"version": schema.StringAttribute{
				Description: "The version of the k8s cluster",
				Required:    true,
			},
			"pools": schema.ListNestedAttribute{
				Description: "The pools of the cluster",
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					poolNamesUniqueValidator{},
				},
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
							Description:   "Cri-o configuration for pool nodes",
							Computed:      true,
							Optional:      true,
							CustomType:    customfield.NewMapType[types.String](ctx),
							ElementType:   types.StringType,
							PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
						},
						"is_public_ipv4": schema.BoolAttribute{
							Description: "Enable public v4 address",
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
						},
						"kubelet_config": schema.MapAttribute{
							Description:   "Kubelet configuration for pool nodes",
							Computed:      true,
							Optional:      true,
							CustomType:    customfield.NewMapType[types.String](ctx),
							ElementType:   types.StringType,
							PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
						},
						"labels": schema.MapAttribute{
							Description:   "Labels applied to the cluster pool",
							Computed:      true,
							Optional:      true,
							CustomType:    customfield.NewMapType[types.String](ctx),
							ElementType:   types.StringType,
							PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
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
							Description:   "Taints applied to the cluster pool",
							Computed:      true,
							Optional:      true,
							CustomType:    customfield.NewMapType[types.String](ctx),
							ElementType:   types.StringType,
							PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					poolsNormalizeOrderPlanModifier(),
				},
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
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"pods_ipv6_pool": schema.StringAttribute{
				Description:   "The IPv6 pool for the pods",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{planmodifiers.UseStateForUnknownIncludingNullString(), stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"services_ip_pool": schema.StringAttribute{
				Description:   "The IP pool for the services",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"services_ipv6_pool": schema.StringAttribute{
				Description:   "The IPv6 pool for the services",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{planmodifiers.UseStateForUnknownIncludingNullString(), stringplanmodifier.RequiresReplaceIfConfigured()},
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
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
					objectplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"autoscaler_config": schema.MapAttribute{
				Description:   "Cluster autoscaler configuration.\n\nIt allows you to override the default cluster-autoscaler parameters provided by the platform with your preferred values.\n\nSupported parameters (in alphabetical order):\n- balance-similar-node-groups (boolean: true/false) - Detect similar node groups and balance the number of nodes between them.\n- expander (string: random, most-pods, least-waste, price, priority, grpc) - Type of node group expander to be used in scale up. Specifying multiple values separated by commas will call the expanders in succession until there is only one option remaining.\n- expendable-pods-priority-cutoff (float) - Pods with priority below cutoff will be expendable. They can be killed without any consideration during scale down and they don't cause scale up. Pods with null priority (PodPriority disabled) are non expendable.\n- ignore-daemonsets-utilization (boolean: true/false) - Should CA ignore DaemonSet pods when calculating resource utilization for scaling down.\n- max-empty-bulk-delete (integer) - Maximum number of empty nodes that can be deleted at the same time.\n- max-graceful-termination-sec (integer) - Maximum number of seconds CA waits for pod termination when trying to scale down a node.\n- max-node-provision-time (duration: e.g., '15m') - The default maximum time CA waits for node to be provisioned - the value can be overridden per node group.\n- max-total-unready-percentage (float) - Maximum percentage of unready nodes in the cluster. After this is exceeded, CA halts operations.\n- new-pod-scale-up-delay (duration: e.g., '10s') - Pods less than this old will not be considered for scale-up. Can be increased for individual pods through annotation.\n- ok-total-unready-count (integer) - Number of allowed unready nodes, irrespective of max-total-unready-percentage.\n- scale-down-delay-after-add (duration: e.g., '10m') - How long after scale up that scale down evaluation resumes.\n- scale-down-delay-after-delete (duration: e.g., '10s') - How long after node deletion that scale down evaluation resumes.\n- scale-down-delay-after-failure (duration: e.g., '3m') - How long after scale down failure that scale down evaluation resumes.\n- scale-down-enabled (boolean: true/false) - Should CA scale down the cluster.\n- scale-down-unneeded-time (duration: e.g., '10m') - How long a node should be unneeded before it is eligible for scale down.\n- scale-down-unready-time (duration: e.g., '20m') - How long an unready node should be unneeded before it is eligible for scale down.\n- scale-down-utilization-threshold (float) - The maximum value between the sum of cpu requests and sum of memory requests of all pods running on the node divided by node's corresponding allocatable resource, below which a node can be considered for scale down.\n- scan-interval (duration: e.g., '10s') - How often cluster is reevaluated for scale up or down.\n- skip-nodes-with-custom-controller-pods (boolean: true/false) - If true cluster autoscaler will never delete nodes with pods owned by custom controllers.\n- skip-nodes-with-local-storage (boolean: true/false) - If true cluster autoscaler will never delete nodes with pods with local storage, e.g. EmptyDir or HostPath.\n- skip-nodes-with-system-pods (boolean: true/false) - If true cluster autoscaler will never delete nodes with pods from kube-system (except for DaemonSet or mirror pods).",
				Computed:      true,
				Optional:      true,
				CustomType:    customfield.NewMapType[types.String](ctx),
				ElementType:   types.StringType,
				PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
			},
			"add_ons": schema.SingleNestedAttribute{
				Description:   "Cluster add-ons configuration",
				Computed:      true,
				Optional:      true,
				CustomType:    customfield.NewNestedObjectType[CloudK8SClusterAddOnsModel](ctx),
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Attributes: map[string]schema.Attribute{
					"slurm": schema.SingleNestedAttribute{
						Description:   "Slurm add-on configuration",
						Computed:      true,
						Optional:      true,
						CustomType:    customfield.NewNestedObjectType[CloudK8SClusterAddOnsSlurmModel](ctx),
						PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description:   "The Slurm add-on will be enabled in the cluster.\n\nThis add-on is only supported in clusters running Kubernetes v1.31 and v1.32 with at least 1 GPU cluster pool and VAST NFS support enabled.",
								Computed:      true,
								Optional:      true,
								PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
							},
							"file_share_id": schema.StringAttribute{
								Description:   "ID of a VAST file share to be used as Slurm storage.\n\nThe Slurm add-on will create separate Persistent Volume Claims for different purposes (controller spool, worker spool, jail) on that file share.\n\nThe file share must have `root_squash` disabled, while `path_length` and `allowed_characters` settings must be set to `NPL`.",
								Computed:      true,
								Optional:      true,
								PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
							},
							"ssh_key_ids": schema.ListAttribute{
								Description:   "IDs of SSH keys to authorize for SSH connection to Slurm login nodes.",
								Computed:      true,
								Optional:      true,
								CustomType:    customfield.NewListType[types.String](ctx),
								ElementType:   types.StringType,
								PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
							},
							"worker_count": schema.Int64Attribute{
								Description: "Size of the worker pool, i.e. the number of Slurm worker nodes.\n\nEach Slurm worker node will be backed by a Pod scheduled on one of cluster's GPU nodes.\n\nNote: Downscaling (reducing worker count) is not supported.",
								Computed:    true,
								Optional:    true,
								Validators: []validator.Int64{
									int64validator.AtLeast(1),
								},
								PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
							},
						},
					},
				},
			},
			"authentication": schema.SingleNestedAttribute{
				Description:   "Authentication settings",
				Computed:      true,
				Optional:      true,
				CustomType:    customfield.NewNestedObjectType[CloudK8SClusterAuthenticationModel](ctx),
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown(), authenticationRemovalPlanModifier()},
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
				Description:   "Cluster CNI settings",
				Computed:      true,
				Optional:      true,
				CustomType:    customfield.NewNestedObjectType[CloudK8SClusterCniModel](ctx),
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
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
			"logging": schema.SingleNestedAttribute{
				Description:   "Logging configuration",
				Computed:      true,
				Optional:      true,
				CustomType:    customfield.NewNestedObjectType[CloudK8SClusterLoggingModel](ctx),
				PlanModifiers: []planmodifier.Object{planmodifiers.UseStateForUnknownIncludingNullObject()},
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
				Description:   "Function creation date",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"creator_task_id": schema.StringAttribute{
				Description:   "Task that created this entity",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"is_public": schema.BoolAttribute{
				Description:   "Cluster is public",
				Computed:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"status": schema.StringAttribute{
				Description:   "Status\nAvailable values: \"Deleting\", \"Provisioned\", \"Provisioning\".",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"Deleting",
						"Provisioned",
						"Provisioning",
					),
				},
			},
		},
	}
}

func (r *CloudK8SClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudK8SClusterResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		&addOnsValidator{},
		&authenticationValidator{},
		&poolFlavorValidator{},
	}
}

type addOnsValidator struct{}

func (v *addOnsValidator) Description(_ context.Context) string {
	return "slurm attributes must all be set when slurm is configured"
}

func (v *addOnsValidator) MarkdownDescription(_ context.Context) string {
	return "slurm attributes must all be set when slurm is configured"
}

func (v *addOnsValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config CloudK8SClusterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If add_ons is not provided or is null/unknown, allow API to fill in defaults
	if config.AddOns.IsNull() || config.AddOns.IsUnknown() {
		return
	}

	addOns, diags := config.AddOns.Value(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If addOns is nil or slurm is not provided or is null/unknown, allow API to fill in defaults
	if addOns == nil || addOns.Slurm.IsNull() || addOns.Slurm.IsUnknown() {
		return
	}

	slurm, diags := addOns.Slurm.Value(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If slurm is nil, allow API to fill in defaults
	if slurm == nil {
		return
	}

	// When slurm block is provided, all attributes must be non-null
	if slurm.Enabled.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("add_ons").AtName("slurm").AtName("enabled"),
			"Missing required attribute",
			"The attribute 'enabled' is required when 'slurm' is configured.",
		)
	}

	if slurm.FileShareID.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("add_ons").AtName("slurm").AtName("file_share_id"),
			"Missing required attribute",
			"The attribute 'file_share_id' is required when 'slurm' is configured.",
		)
	}

	if slurm.SSHKeyIDs.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("add_ons").AtName("slurm").AtName("ssh_key_ids"),
			"Missing required attribute",
			"The attribute 'ssh_key_ids' is required when 'slurm' is configured.",
		)
	}

	if slurm.WorkerCount.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("add_ons").AtName("slurm").AtName("worker_count"),
			"Missing required attribute",
			"The attribute 'worker_count' is required when 'slurm' is configured.",
		)
	}
}

type authenticationValidator struct{}

func (v *authenticationValidator) Description(_ context.Context) string {
	return "oidc issuer_url and client_id must be set when oidc is configured"
}

func (v *authenticationValidator) MarkdownDescription(_ context.Context) string {
	return "oidc issuer_url and client_id must be set when oidc is configured"
}

func (v *authenticationValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config CloudK8SClusterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If authentication is not provided or is null/unknown, allow API to fill in defaults
	if config.Authentication.IsNull() || config.Authentication.IsUnknown() {
		return
	}

	authentication, diags := config.Authentication.Value(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If authentication is nil or oidc is not provided or is null/unknown, allow API to fill in defaults
	if authentication == nil || authentication.Oidc.IsNull() || authentication.Oidc.IsUnknown() {
		return
	}

	oidc, diags := authentication.Oidc.Value(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If oidc is nil, allow API to fill in defaults
	if oidc == nil {
		return
	}

	// When oidc block is provided, issuer_url and client_id must be non-null
	if oidc.IssuerURL.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("authentication").AtName("oidc").AtName("issuer_url"),
			"Missing required attribute",
			"The attribute 'issuer_url' is required when 'oidc' is configured.",
		)
	}

	if oidc.ClientID.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("authentication").AtName("oidc").AtName("client_id"),
			"Missing required attribute",
			"The attribute 'client_id' is required when 'oidc' is configured.",
		)
	}
}

// authenticationRemovalPlanModifier handles the case where a user removes the entire
// authentication block from their config. When this happens, we need to transform
// the null config value into { oidc = null } so the API receives the proper payload
// to remove OIDC authentication.
func authenticationRemovalPlanModifier() planmodifier.Object {
	return authenticationRemovalModifier{}
}

type authenticationRemovalModifier struct{}

func (m authenticationRemovalModifier) Description(_ context.Context) string {
	return "Transforms null authentication config to { oidc = null } when removing OIDC"
}

func (m authenticationRemovalModifier) MarkdownDescription(_ context.Context) string {
	return "Transforms null authentication config to `{ oidc = null }` when removing OIDC"
}

func (m authenticationRemovalModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// If there's no state (new resource), nothing to do
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	// If config is not null (user explicitly set authentication), let normal flow handle it
	if !req.ConfigValue.IsNull() {
		return
	}

	// Config is null (user removed authentication block) and state exists
	// Check if state has oidc configured
	var stateAuth CloudK8SClusterAuthenticationModel
	diags := req.StateValue.As(ctx, &stateAuth, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// If state doesn't have oidc configured, nothing to remove
	if stateAuth.Oidc.IsNull() {
		return
	}

	// State has oidc configured, so we need to plan for removal
	// Create an authentication object with oidc = null
	removalAuth := CloudK8SClusterAuthenticationModel{
		Oidc: customfield.NullObject[CloudK8SClusterAuthenticationOidcModel](ctx),
	}

	// Convert to NestedObject and set as plan value
	nestedObj, diags := customfield.NewObject(ctx, &removalAuth)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Extract the underlying ObjectValue from NestedObject
	resp.PlanValue = nestedObj.ObjectValue
}

type poolFlavorValidator struct{}

func (v *poolFlavorValidator) Description(_ context.Context) string {
	return "validates pool attributes based on flavor_id type (VM-based vs baremetal)"
}

func (v *poolFlavorValidator) MarkdownDescription(_ context.Context) string {
	return "validates pool attributes based on flavor_id type (VM-based vs baremetal)"
}

// isVMBasedFlavor checks if the flavor_id represents a VM-based flavor.
// VM-based flavors have a prefix "g" or "a" (e.g., "g1-standard-2-4", "a1-gpu-v100-16g").
func isVMBasedFlavor(flavorID string) bool {
	return strings.HasPrefix(flavorID, "g") || strings.HasPrefix(flavorID, "a")
}

func (v *poolFlavorValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config CloudK8SClusterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If pools is not provided, nothing to validate
	if config.Pools == nil {
		return
	}

	for i, pool := range *config.Pools {
		if pool == nil {
			continue
		}

		// Skip validation if flavor_id is unknown (e.g., from a variable)
		if pool.FlavorID.IsUnknown() || pool.FlavorID.IsNull() {
			continue
		}

		flavorID := pool.FlavorID.ValueString()
		poolPath := path.Root("pools").AtSetValue(pool.FlavorID)
		_ = i // poolPath uses set value, not index

		if isVMBasedFlavor(flavorID) {
			// VM-based flavor: servergroup_policy, boot_volume_size, boot_volume_type must NOT be null
			if pool.ServergroupPolicy.IsNull() {
				resp.Diagnostics.AddAttributeError(
					poolPath.AtName("servergroup_policy"),
					"Missing required attribute for VM-based flavor",
					"The attribute 'servergroup_policy' is required when 'flavor_id' is a VM-based flavor (starts with 'g' or 'a').",
				)
			}

			if pool.BootVolumeSize.IsNull() {
				resp.Diagnostics.AddAttributeError(
					poolPath.AtName("boot_volume_size"),
					"Missing required attribute for VM-based flavor",
					"The attribute 'boot_volume_size' is required when 'flavor_id' is a VM-based flavor (starts with 'g' or 'a').",
				)
			}

			if pool.BootVolumeType.IsNull() {
				resp.Diagnostics.AddAttributeError(
					poolPath.AtName("boot_volume_type"),
					"Missing required attribute for VM-based flavor",
					"The attribute 'boot_volume_type' is required when 'flavor_id' is a VM-based flavor (starts with 'g' or 'a').",
				)
			}
		} else {
			// Baremetal flavor: servergroup_policy, boot_volume_size, boot_volume_type must be null
			if !pool.ServergroupPolicy.IsNull() {
				resp.Diagnostics.AddAttributeError(
					poolPath.AtName("servergroup_policy"),
					"Invalid attribute for baremetal flavor",
					"The attribute 'servergroup_policy' must not be set when 'flavor_id' is a baremetal flavor (does not start with 'g' or 'a').",
				)
			}

			if !pool.BootVolumeSize.IsNull() {
				resp.Diagnostics.AddAttributeError(
					poolPath.AtName("boot_volume_size"),
					"Invalid attribute for baremetal flavor",
					"The attribute 'boot_volume_size' must not be set when 'flavor_id' is a baremetal flavor (does not start with 'g' or 'a').",
				)
			}

			if !pool.BootVolumeType.IsNull() {
				resp.Diagnostics.AddAttributeError(
					poolPath.AtName("boot_volume_type"),
					"Invalid attribute for baremetal flavor",
					"The attribute 'boot_volume_type' must not be set when 'flavor_id' is a baremetal flavor (does not start with 'g' or 'a').",
				)
			}
		}
	}
}

// poolsNormalizeOrderPlanModifier reorders plan pools to match state order
// (correlating by name) so that list index correlation works correctly.
// New pools are appended at the end.
func poolsNormalizeOrderPlanModifier() planmodifier.List {
	return poolsNormalizeOrderModifier{}
}

type poolsNormalizeOrderModifier struct{}

func (m poolsNormalizeOrderModifier) Description(_ context.Context) string {
	return "Reorders plan pools to match state order (by name) for stable diffs"
}

func (m poolsNormalizeOrderModifier) MarkdownDescription(_ context.Context) string {
	return "Reorders plan pools to match state order (by name) for stable diffs"
}

func (m poolsNormalizeOrderModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// If there's no state (new resource), nothing to reorder
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	// If plan is null/unknown, let normal flow handle it
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	// Build ordered list of state pool names and a set for quick lookup
	var stateOrder []string
	stateNames := make(map[string]bool)
	for _, elem := range req.StateValue.Elements() {
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}
		name := getPoolName(obj)
		if name != "" {
			stateOrder = append(stateOrder, name)
			stateNames[name] = true
		}
	}

	// Index plan pools by name and build ordered list
	planPoolsByName := make(map[string]types.Object)
	planNamesSet := make(map[string]bool)
	for _, elem := range req.PlanValue.Elements() {
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}
		name := getPoolName(obj)
		if name != "" {
			planPoolsByName[name] = obj
			planNamesSet[name] = true
		}
	}

	// Check if there are any additions or deletions
	// If so, don't reorder - preserve config order to avoid plan validation errors
	hasAdditions := false
	hasDeletions := false
	for name := range planNamesSet {
		if !stateNames[name] {
			hasAdditions = true
			break
		}
	}
	for name := range stateNames {
		if !planNamesSet[name] {
			hasDeletions = true
			break
		}
	}

	// Only reorder if we're just reordering existing pools (no additions/deletions)
	// This prevents plan validation errors when adding pools "in the middle"
	if hasAdditions || hasDeletions {
		return
	}

	// Build new list: pools in state order (since all pools exist in both)
	var reorderedElements []attr.Value
	for _, name := range stateOrder {
		if pool, exists := planPoolsByName[name]; exists {
			reorderedElements = append(reorderedElements, pool)
		}
	}

	// Create new list with reordered elements
	newList, diags := types.ListValue(req.PlanValue.ElementType(ctx), reorderedElements)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = newList
}

// getPoolName extracts the "name" attribute from a pool object
func getPoolName(obj types.Object) string {
	attrs := obj.Attributes()
	if nameAttr, ok := attrs["name"]; ok {
		if nameStr, ok := nameAttr.(types.String); ok && !nameStr.IsNull() && !nameStr.IsUnknown() {
			return nameStr.ValueString()
		}
	}
	return ""
}

type poolNamesUniqueValidator struct{}

func (v poolNamesUniqueValidator) Description(_ context.Context) string {
	return "ensures all pool names are unique"
}

func (v poolNamesUniqueValidator) MarkdownDescription(_ context.Context) string {
	return "ensures all pool names are unique"
}

func (v poolNamesUniqueValidator) ValidateList(_ context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	namesSeen := make(map[string]bool)
	for _, elem := range req.ConfigValue.Elements() {
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}

		name := getPoolName(obj)
		if name == "" {
			continue
		}

		if namesSeen[name] {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Duplicate pool name",
				"Each pool must have a unique name. The pool name '"+name+"' is used more than once.",
			)
			return
		}
		namesSeen[name] = true
	}
}
