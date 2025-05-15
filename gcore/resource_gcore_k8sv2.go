package gcore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/servergroup/v1/servergroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	K8sPoint         = "k8s/clusters"
	tasksPoint       = "tasks"
	K8sCreateTimeout = 3600

	k8sSgMetadataKey = "gcloud_cluster_name"
)

var k8sCreateTimeout = time.Second * time.Duration(K8sCreateTimeout)

func resourceK8sV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceK8sV2Create,
		ReadContext:   resourceK8sV2Read,
		UpdateContext: resourceK8sV2Update,
		DeleteContext: resourceK8sV2Delete,
		Description:   "Represent k8s cluster with one default pool.",
		Timeouts: &schema.ResourceTimeout{
			Create: &k8sCreateTimeout,
			Update: &k8sCreateTimeout,
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, clusterName, err := ImportStringParser(d.Id())
				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.Set("name", clusterName)
				d.SetId(clusterName)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Cluster name.",
				Required:    true,
				ForceNew:    true,
			},
			"authentication": {
				Type:        schema.TypeList,
				Description: "Cluster authentication configuration.",
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oidc": {
							Type:        schema.TypeList,
							Description: "OpenID Connect configuration settings.",
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client_id": {
										Type:        schema.TypeString,
										Description: "A client id that all tokens must be issued for.",
										Optional:    true,
										Computed:    true,
									},
									"groups_claim": {
										Type:        schema.TypeString,
										Description: "JWT claim to use as the user's group.",
										Optional:    true,
										Computed:    true,
									},
									"groups_prefix": {
										Type:        schema.TypeString,
										Description: "Prefix prepended to group claims to prevent clashes with existing names.",
										Optional:    true,
										Computed:    true,
									},
									"issuer_url": {
										Type:        schema.TypeString,
										Description: "URL of the provider that allows the API server to discover public signing keys. Only URLs that use the https:// scheme are accepted.",
										Optional:    true,
										Computed:    true,
									},
									"required_claims": {
										Type:        schema.TypeMap,
										Description: "A map describing required claims in the ID Token. Each claim is verified to be present in the ID Token with a matching value.",
										Optional:    true,
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"signing_algs": {
										Type:        schema.TypeSet,
										Description: "Accepted signing algorithms. Supported values are: RS256, RS384, RS512, ES256, ES384, ES512, PS256, PS384, PS512.",
										Optional:    true,
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"username_claim": {
										Type:        schema.TypeString,
										Description: "JWT claim to use as the user name. When not specified, the `sub` claim will be used.",
										Optional:    true,
										Computed:    true,
									},
									"username_prefix": {
										Type:        schema.TypeString,
										Description: "Prefix prepended to username claims to prevent clashes with existing names.",
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"autoscaler_config": {
				Type:        schema.TypeMap,
				Description: "Cluster autoscaler configuration params. Keys and values are expected to follow the cluster-autoscaler option format.",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cni": {
				Type:        schema.TypeList,
				Description: "Cluster CNI configuration.",
				MaxItems:    1,
				MinItems:    1,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:        schema.TypeString,
							Description: "CNI provider to use when creating the cluster. Supported values are: calico, cilium. The default value is calico.",
							Optional:    true,
							ForceNew:    true,
							Default:     clusters.CalicoProvider.String(),
						},
						"cilium": {
							Type:        schema.TypeList,
							Description: "Cilium CNI configuration.",
							MaxItems:    1,
							MinItems:    1,
							Optional:    true,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mask_size": {
										Type:        schema.TypeInt,
										Description: "Specifies the size allocated from pods_ip_pool CIDR to node.ipam.podCIDRs. The default value is 24.",
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
									},
									"mask_size_v6": {
										Type:        schema.TypeInt,
										Description: "Specifies the size allocated from pods_ipv6_pool CIDR to node.ipam.podCIDRs. The default value is 120.",
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
									},
									"tunnel": {
										Type:        schema.TypeString,
										Description: "Tunneling protocol to use in tunneling mode and for ad-hoc tunnels. The default value is geneve.",
										Optional:    true,
										Default:     "geneve",
									},
									"encryption": {
										Type:        schema.TypeBool,
										Description: "Enables transparent network encryption. The default value is false.",
										Optional:    true,
										Default:     false,
									},
									"lb_mode": {
										Type:        schema.TypeString,
										Description: "The operation mode of load balancing for remote backends. Supported values are snat, dsr, hybrid. The default value is snat.",
										Optional:    true,
										Default:     "snat",
									},
									"lb_acceleration": {
										Type:        schema.TypeBool,
										Description: "Enables load balancer acceleration via XDP. The default value is false.",
										Optional:    true,
										Default:     false,
									},
									"routing_mode": {
										Type:        schema.TypeString,
										Description: "Enables native-routing mode or tunneling mode. The default value is tunnel.",
										Optional:    true,
										Default:     "tunnel",
									},
									"hubble_relay": {
										Type:        schema.TypeBool,
										Description: "Enables Hubble Relay. The default value is false.",
										Optional:    true,
										Default:     false,
									},
									"hubble_ui": {
										Type:        schema.TypeBool,
										Description: "Enables Hubble UI. Requires `hubble_relay=true`. The default value is false.",
										Optional:    true,
										Default:     false,
									},
								},
							},
						}},
				},
			},
			"ddos_profile": {
				Type:        schema.TypeList,
				Description: "DDoS profile configuration.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Indicates if the DDoS profile is enabled.",
							Required:    true,
						},
						"fields": {
							Type:        schema.TypeList,
							Description: "List of fields for the DDoS profile.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"base_field": {
										Type:        schema.TypeInt,
										Description: "Base field ID.",
										Required:    true,
									},
									"field_value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Complex value. Only one of 'value' or 'field_value' must be specified.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Basic type value. Only one of 'value' or 'field_value' must be specified.",
									},
								},
							},
						},
						"profile_template": {
							Type:        schema.TypeInt,
							Description: "Profile template ID.",
							Optional:    true,
						},
						"profile_template_name": {
							Type:        schema.TypeString,
							Description: "Profile template name.",
							Optional:    true,
						},
					},
				},
			},
			"fixed_network": {
				Type:        schema.TypeString,
				Description: "Fixed network used to allocate network addresses for cluster nodes.",
				Optional:    true,
				ForceNew:    true,
			},
			"fixed_subnet": {
				Type:        schema.TypeString,
				Description: "Fixed subnet used to allocate network addresses for cluster nodes. Subnet should have a router.",
				Optional:    true,
				ForceNew:    true,
			},
			"pods_ip_pool": {
				Type:        schema.TypeString,
				Description: "Pods IPv4 IP pool in CIDR notation.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"services_ip_pool": {
				Type:        schema.TypeString,
				Description: "Services IPv4 IP pool in CIDR notation.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"pods_ipv6_pool": {
				Type:        schema.TypeString,
				Description: "Pods IPv6 IP pool in CIDR notation.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"services_ipv6_pool": {
				Type:        schema.TypeString,
				Description: "Services IPv6 IP pool in CIDR notation.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"keypair": {
				Type:        schema.TypeString,
				Description: "Name of the keypair used for SSH access to nodes.",
				Required:    true,
				ForceNew:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Kubernetes version.",
				Required:    true,
			},
			"is_ipv6": {
				Type:        schema.TypeBool,
				Description: "Enable public IPv6 address.",
				Optional:    true,
				ForceNew:    true,
			},
			"pool": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Cluster pool name. Changing the value of this attribute will trigger recreation of the cluster pool.",
							Required:    true,
						},
						"flavor_id": {
							Type:        schema.TypeString,
							Description: "Cluster pool node flavor ID. Changing the value of this attribute will trigger recreation of the cluster pool.",
							Required:    true,
						},
						"min_node_count": {
							Type:        schema.TypeInt,
							Description: "Minimum number of nodes in the cluster pool.",
							Required:    true,
						},
						"servergroup_policy": {
							Type:        schema.TypeString,
							Description: "Server group policy: anti-affinity, soft-anti-affinity or affinity",
							Optional:    true,
						},
						"max_node_count": {
							Type:        schema.TypeInt,
							Description: "Maximum number of nodes in the cluster pool.",
							Optional:    true,
							Computed:    true,
						},
						"node_count": {
							Type:        schema.TypeInt,
							Description: "Current node count in the cluster pool.",
							Computed:    true,
						},
						"boot_volume_type": {
							Type:        schema.TypeString,
							Description: "Cluster pool boot volume type. Must be set only for VM pools. Available values are 'standard', 'ssd_hiiops', 'cold', 'ultra'. Changing the value of this attribute will trigger recreation of the cluster pool.",
							Optional:    true,
							Computed:    true,
						},
						"boot_volume_size": {
							Type:        schema.TypeInt,
							Description: "Cluster pool boot volume size. Must be set only for VM pools. Changing the value of this attribute will trigger recreation of the cluster pool.",
							Optional:    true,
							Computed:    true,
						},
						"auto_healing_enabled": {
							Type:        schema.TypeBool,
							Description: "Enable/disable auto healing of cluster pool nodes.",
							Optional:    true,
							Computed:    true,
						},
						"is_public_ipv4": {
							Type:        schema.TypeBool,
							Description: "Assign public IPv4 address to nodes in this pool. Changing the value of this attribute will trigger recreation of the cluster pool.",
							Optional:    true,
							Computed:    true,
						},
						"labels": {
							Type:        schema.TypeMap,
							Description: "Labels applied to the cluster pool nodes.",
							Optional:    true,
							Computed:    true,
						},
						"taints": {
							Type:        schema.TypeMap,
							Description: "Taints applied to the cluster pool nodes.",
							Optional:    true,
							Computed:    true,
						},
						"crio_config": {
							Type:        schema.TypeMap,
							Description: "Crio configuration for pool nodes. Keys and values are expected to follow the crio option format.",
							Optional:    true,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"kubelet_config": {
							Type:        schema.TypeMap,
							Description: "Kubelet configuration for pool nodes. Keys and values are expected to follow the kubelet configuration file format.",
							Optional:    true,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Cluster pool status.",
							Computed:    true,
						},
						"servergroup_name": {
							Type:        schema.TypeString,
							Description: "Server group name",
							Computed:    true,
						},
						"servergroup_id": {
							Type:        schema.TypeString,
							Description: "Server group id",
							Computed:    true,
						},
						"created_at": {
							Type:        schema.TypeString,
							Description: "Cluster pool creation date.",
							Computed:    true,
						},
					},
				},
			},
			"security_group_rules": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Firewall rules control what inbound(ingress) and outbound(egress) traffic is allowed to enter or leave a Instance. At least one 'egress' rule should be set",
				Set:         secGroupUniqueID,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:        schema.TypeString,
							Required:    true,
							Description: fmt.Sprintf("Available value is '%s', '%s'", types.RuleDirectionIngress, types.RuleDirectionEgress),
							ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
								val := v.(string)
								switch types.RuleDirection(val) {
								case types.RuleDirectionIngress, types.RuleDirectionEgress:
									return nil
								}
								return diag.Errorf("wrong direction '%s', available value is '%s', '%s'", val, types.RuleDirectionIngress, types.RuleDirectionEgress)
							},
						},
						"ethertype": {
							Type:        schema.TypeString,
							Required:    true,
							Description: fmt.Sprintf("Available value is '%s', '%s'", types.EtherTypeIPv4, types.EtherTypeIPv6),
							ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
								val := v.(string)
								switch types.EtherType(val) {
								case types.EtherTypeIPv4, types.EtherTypeIPv6:
									return nil
								}
								return diag.Errorf("wrong ethertype '%s', available value is '%s', '%s'", val, types.EtherTypeIPv4, types.EtherTypeIPv6)
							},
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: fmt.Sprintf("Available value is %s", strings.Join(types.Protocol("").StringList(), ",")),
						},
						"port_range_min": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          0,
							ValidateDiagFunc: validatePortRange,
						},
						"port_range_max": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          0,
							ValidateDiagFunc: validatePortRange,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"remote_ip_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group ID.",
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Cluster status.",
				Computed:    true,
			},
			"is_public": {
				Type:        schema.TypeBool,
				Description: "True if the cluster is public.",
				Computed:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Cluster creation date.",
				Computed:    true,
			},
			"creator_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ValidateChange("pool", func(ctx context.Context, old, new, meta interface{}) error {
				for _, p := range new.([]interface{}) {
					pool := p.(map[string]interface{})
					if resourceK8sV2IsVMFlavor(pool["flavor_id"].(string)) {
						if pool["servergroup_policy"].(string) == "" {
							return fmt.Errorf("servergroup_policy is required for flavor %v", pool["flavor_id"])
						}
					} else {
						if pool["servergroup_policy"].(string) != "" {
							return fmt.Errorf("servergroup_policy cannot be set for flavor %v", pool["flavor_id"])
						}
					}
				}
				return nil
			})),
	}
}

func resourceK8sV2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start k8s cluster creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, K8sPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := clusters.CreateOpts{
		Name:         d.Get("name").(string),
		FixedNetwork: d.Get("fixed_network").(string),
		FixedSubnet:  d.Get("fixed_subnet").(string),
		KeyPair:      d.Get("keypair").(string),
		Version:      d.Get("version").(string),
		IsIPV6:       d.Get("is_ipv6").(bool),
	}

	if authI, ok := d.GetOk("authentication"); ok {
		authA := authI.([]interface{})
		auth := authA[0].(map[string]interface{})
		opts.Authentication = &clusters.AuthenticationCreateOpts{}
		if oidcI, ok := auth["oidc"]; ok {
			oidcA := oidcI.([]interface{})
			if len(oidcA) != 0 {
				oidc := oidcA[0].(map[string]interface{})
				opts.Authentication.OIDC = &clusters.OIDCCreateOpts{
					ClientID:       oidc["client_id"].(string),
					GroupsClaim:    oidc["groups_claim"].(string),
					GroupsPrefix:   oidc["groups_prefix"].(string),
					IssuerURL:      oidc["issuer_url"].(string),
					UsernameClaim:  oidc["username_claim"].(string),
					UsernamePrefix: oidc["username_prefix"].(string),
				}
				if len(oidc["required_claims"].(map[string]interface{})) > 0 {
					opts.Authentication.OIDC.RequiredClaims = map[string]string{}
					for k, v := range oidc["required_claims"].(map[string]interface{}) {
						opts.Authentication.OIDC.RequiredClaims[k] = v.(string)
					}
				}
				if algs, ok := oidc["signing_algs"].(*schema.Set); ok {
					for _, alg := range algs.List() {
						opts.Authentication.OIDC.SigningAlgs = append(opts.Authentication.OIDC.SigningAlgs, alg.(string))
					}
				}
			}
		}
	}

	if autoscalerCfgI, ok := d.GetOk("autoscaler_config"); ok {
		autoscalerCfg := autoscalerCfgI.(map[string]interface{})
		opts.AutoscalerConfig = map[string]string{}
		for k, v := range autoscalerCfg {
			opts.AutoscalerConfig[k] = v.(string)
		}
	}

	if cniI, ok := d.GetOk("cni"); ok {
		cniA := cniI.([]interface{})
		cni := cniA[0].(map[string]interface{})
		opts.CNI = &clusters.CNICreateOpts{Provider: clusters.CNIProvider(cni["provider"].(string))}
		if cni["provider"].(string) == "cilium" {
			if ciliumI, ok := cni["cilium"]; ok {
				ciliumA := ciliumI.([]interface{})
				if len(ciliumA) != 0 {
					cilium := ciliumA[0].(map[string]interface{})
					opts.CNI.Cilium = &clusters.CiliumCreateOpts{
						MaskSize:                 cilium["mask_size"].(int),
						MaskSizeV6:               cilium["mask_size_v6"].(int),
						Tunnel:                   clusters.TunnelType(cilium["tunnel"].(string)),
						Encryption:               cilium["encryption"].(bool),
						LoadBalancerMode:         clusters.LBModeType(cilium["lb_mode"].(string)),
						LoadBalancerAcceleration: cilium["lb_acceleration"].(bool),
						RoutingMode:              clusters.RoutingModeType(cilium["routing_mode"].(string)),
						HubbleRelay:              cilium["hubble_relay"].(bool),
						HubbleUI:                 cilium["hubble_ui"].(bool),
					}
				}
			}
		}
	}

	if ddosProfileI, ok := d.GetOk("ddos_profile"); ok {
		ddosProfile := ddosProfileI.([]interface{})
		if len(ddosProfile) != 0 {
			profile := ddosProfile[0].(map[string]interface{})
			opts.DDoSProfile = &clusters.DDoSProfileCreateOpts{
				Enabled: profile["enabled"].(bool),
			}
			if fields, ok := profile["fields"].([]interface{}); ok {
				if len(fields) != 0 {
					opts.DDoSProfile.Fields = make([]clusters.DDoSProfileField, len(fields))
					for i, field := range fields {
						fieldMap := field.(map[string]interface{})
						opts.DDoSProfile.Fields[i] = clusters.DDoSProfileField{
							BaseField: fieldMap["base_field"].(int),
						}
						if value, ok := fieldMap["value"].(string); ok {
							opts.DDoSProfile.Fields[i].Value = &value
						}
						if fieldValue, ok := fieldMap["field_value"].(string); ok {
							opts.DDoSProfile.Fields[i].FieldValue = &fieldValue
						}
					}
				}
			}
			if profileTemplate, ok := profile["profile_template"].(int); ok {
				opts.DDoSProfile.ProfileTemplate = &profileTemplate
			}
			if profileTemplateName, ok := profile["profile_template_name"].(string); ok {
				opts.DDoSProfile.ProfileTemplateName = &profileTemplateName
			}
		}
	}

	if podsIP, ok := d.GetOk("pods_ip_pool"); ok {
		gccidr, err := parseCIDRFromString(podsIP.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.PodsIPPool = &gccidr
	}

	if svcIP, ok := d.GetOk("services_ip_pool"); ok {
		gccidr, err := parseCIDRFromString(svcIP.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.ServicesIPPool = &gccidr
	}

	if podsIPV6, ok := d.GetOk("pods_ipv6_pool"); ok {
		gccidr, err := parseCIDRFromString(podsIPV6.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.PodsIPV6Pool = &gccidr
	}

	if svcIPV6, ok := d.GetOk("services_ipv6_pool"); ok {
		gccidr, err := parseCIDRFromString(svcIPV6.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.ServicesIPV6Pool = &gccidr
	}

	for _, poolRaw := range d.Get("pool").([]interface{}) {
		pool := poolRaw.(map[string]interface{})
		poolOpts := pools.CreateOpts{
			Name:               pool["name"].(string),
			FlavorID:           pool["flavor_id"].(string),
			MinNodeCount:       pool["min_node_count"].(int),
			MaxNodeCount:       pool["max_node_count"].(int),
			BootVolumeSize:     pool["boot_volume_size"].(int),
			BootVolumeType:     volumes.VolumeType(pool["boot_volume_type"].(string)),
			AutoHealingEnabled: pool["auto_healing_enabled"].(bool),
			IsPublicIPv4:       pool["is_public_ipv4"].(bool),
			ServerGroupPolicy:  servergroups.ServerGroupPolicy(pool["servergroup_policy"].(string)),
		}
		if labels, ok := pool["labels"].(map[string]interface{}); ok {
			poolOpts.Labels = map[string]string{}
			for k, v := range labels {
				poolOpts.Labels[k] = v.(string)
			}
		}
		if taints, ok := pool["taints"].(map[string]interface{}); ok {
			poolOpts.Taints = map[string]string{}
			for k, v := range taints {
				poolOpts.Taints[k] = v.(string)
			}
		}
		if crioCfg, ok := pool["crio_config"].(map[string]interface{}); ok {
			poolOpts.CrioConfig = map[string]string{}
			for k, v := range crioCfg {
				poolOpts.CrioConfig[k] = v.(string)
			}
		}
		if kubeletCfg, ok := pool["kubelet_config"].(map[string]interface{}); ok {
			poolOpts.KubeletConfig = map[string]string{}
			for k, v := range kubeletCfg {
				poolOpts.KubeletConfig[k] = v.(string)
			}
		}
		opts.Pools = append(opts.Pools, poolOpts)
	}

	results, err := clusters.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)

	tasksClient, err := CreateClient(provider, d, tasksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	clusterName, err := tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		cluster, err := clusters.Get(client, opts.Name).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot create k8s cluster with name: %s. Error: %w", opts.Name, err)
		}
		return cluster.Name, nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterName.(string))

	sgClient, err := CreateClient(provider, d, securityGroupPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	sgs, err := securitygroups.ListAll(
		sgClient, securitygroups.ListOpts{MetadataKV: map[string]string{k8sSgMetadataKey: clusterName.(string)}},
	)
	sg := getSuitableSecurityGroup(sgs, clusterName.(string), d.Get("project_id").(int), d.Get("region_id").(int))
	if sg != nil {
		rawRules := d.Get("security_group_rules").(*schema.Set).List()
		if len(rawRules) != 0 {
			usersRules := convertToSecurityGroupRules(rawRules)

			for _, rule := range usersRules {
				_, err = securitygroups.AddRule(sgClient, sg.ID, rule).Extract()
				if err != nil {
					log.Println("[ERROR] Cannot add rule to security group", err)
				}
			}
		}
	}

	resourceK8sV2Read(ctx, d, m)
	log.Printf("[DEBUG] Finish k8s cluster creating (%s)", clusterName)
	return diags
}

func getSuitableSecurityGroup(sgs []securitygroups.SecurityGroup, name string, projectID, regionID int) *securitygroups.SecurityGroup {
	sgName := fmt.Sprintf("%s-%d-%d-worker", name, regionID, projectID)
	for _, sg := range sgs {
		if sg.Name == sgName {
			return &sg
		}
	}
	return nil
}

func filterSecurityGroupRules(rules []securitygroups.SecurityGroupRule) []securitygroups.SecurityGroupRule {
	newRules := []securitygroups.SecurityGroupRule{}
	for _, rule := range rules {
		if *rule.Description == "system" {
			continue
		}
		newRules = append(newRules, rule)
	}
	return newRules
}

func resourceK8sV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start k8s cluster reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, K8sPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterName := d.Get("name").(string)
	cluster, err := clusters.Get(client, clusterName).Extract()
	if err != nil {
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			d.SetId("")
			log.Printf("[WARNING] k8s cluster not found, removing from state")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("name", cluster.Name)
	d.Set("fixed_network", cluster.FixedNetwork)
	d.Set("fixed_subnet", cluster.FixedSubnet)
	d.Set("keypair", cluster.KeyPair)
	d.Set("version", cluster.Version)
	d.Set("status", cluster.Status)
	d.Set("is_public", cluster.IsPublic)
	d.Set("created_at", cluster.CreatedAt.Format(time.RFC850))
	d.Set("creator_task_id", cluster.CreatorTaskID)
	d.Set("task_id", cluster.TaskID)
	d.Set("is_ipv6", cluster.IsIPV6)
	d.Set("autoscaler_config", cluster.AutoscalerConfig)

	if cluster.PodsIPPool != nil {
		d.Set("pods_ip_pool", cluster.PodsIPPool.String())
	}
	if cluster.ServicesIPPool != nil {
		d.Set("services_ip_pool", cluster.ServicesIPPool.String())
	}
	if cluster.PodsIPV6Pool != nil {
		d.Set("pods_ipv6_pool", cluster.PodsIPV6Pool.String())
	}
	if cluster.ServicesIPV6Pool != nil {
		d.Set("services_ipv6_pool", cluster.ServicesIPV6Pool.String())
	}

	if cluster.Authentication != nil {
		v := map[string]interface{}{}
		if cluster.Authentication.OIDC != nil {
			v["oidc"] = []map[string]interface{}{{
				"client_id":       cluster.Authentication.OIDC.ClientID,
				"groups_claim":    cluster.Authentication.OIDC.GroupsClaim,
				"groups_prefix":   cluster.Authentication.OIDC.GroupsPrefix,
				"issuer_url":      cluster.Authentication.OIDC.IssuerURL,
				"required_claims": cluster.Authentication.OIDC.RequiredClaims,
				"signing_algs":    cluster.Authentication.OIDC.SigningAlgs,
				"username_claim":  cluster.Authentication.OIDC.UsernameClaim,
				"username_prefix": cluster.Authentication.OIDC.UsernamePrefix,
			}}
		}
		if err := d.Set("authentication", []interface{}{v}); err != nil {
			return diag.FromErr(err)
		}
	}

	if cluster.CNI != nil {
		v := map[string]interface{}{
			"provider": cluster.CNI.Provider.String(),
		}
		if cluster.CNI.Cilium != nil {
			v["cilium"] = []map[string]interface{}{{
				"mask_size":       cluster.CNI.Cilium.MaskSize,
				"mask_size_v6":    cluster.CNI.Cilium.MaskSizeV6,
				"tunnel":          cluster.CNI.Cilium.Tunnel.String(),
				"encryption":      cluster.CNI.Cilium.Encryption,
				"lb_mode":         cluster.CNI.Cilium.LoadBalancerMode.String(),
				"lb_acceleration": cluster.CNI.Cilium.LoadBalancerAcceleration,
				"routing_mode":    cluster.CNI.Cilium.RoutingMode.String(),
				"hubble_relay":    cluster.CNI.Cilium.HubbleRelay,
				"hubble_ui":       cluster.CNI.Cilium.HubbleUI,
			}}
		}
		if err := d.Set("cni", []interface{}{v}); err != nil {
			return diag.FromErr(err)
		}
	}

	if cluster.DDoSProfile != nil {
		v := map[string]interface{}{
			"enabled": cluster.DDoSProfile.Enabled,
		}
		if cluster.DDoSProfile.Fields != nil {
			v["fields"] = []interface{}{}
			for _, field := range cluster.DDoSProfile.Fields {
				v["fields"] = append(v["fields"].([]interface{}), map[string]interface{}{
					"base_field": field.BaseField,
					"value":      field.Value,
				})
			}
		}
		if cluster.DDoSProfile.ProfileTemplate != nil {
			v["profile_template"] = *cluster.DDoSProfile.ProfileTemplate
		}
		if cluster.DDoSProfile.ProfileTemplateName != nil {
			v["profile_template_name"] = *cluster.DDoSProfile.ProfileTemplateName
		}
	}

	poolMap := map[string]pools.ClusterPool{}
	for _, pool := range cluster.Pools {
		poolMap[pool.Name] = pool
	}

	// Returned pool order needs to match TF state or users will see broken diff,
	// so we first process all pools stored in the state file, and then append any remaining pools.
	var poolData []interface{}
	for _, rawPool := range d.Get("pool").([]interface{}) {
		pool := rawPool.(map[string]interface{})
		poolName := pool["name"].(string)
		if p, ok := poolMap[poolName]; ok {
			poolData = append(poolData, resourceK8sV2PoolDataFromPool(p))
			delete(poolMap, poolName)
		} else {
			// prevent breaking diff when a pool from state file is missing
			log.Printf("[DEBUG] Returning cluster pool placeholder for %q\n", poolName)
			poolData = append(poolData, map[string]interface{}{})
		}
	}
	for _, pool := range poolMap {
		poolData = append(poolData, resourceK8sV2PoolDataFromPool(pool))
	}
	if err := d.Set("pool", poolData); err != nil {
		return diag.FromErr(err)
	}

	// get cluster's security group
	sgClient, err := CreateClient(provider, d, securityGroupPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	sgs, err := securitygroups.ListAll(
		sgClient, securitygroups.ListOpts{MetadataKV: map[string]string{k8sSgMetadataKey: clusterName}},
	)
	sg := getSuitableSecurityGroup(sgs, clusterName, d.Get("project_id").(int), d.Get("region_id").(int))
	if err == nil && sg != nil {
		d.Set("security_group_id", sg.ID)
		// todo read security group for the cluster
		rulesRaw := convertSecurityGroupRules(filterSecurityGroupRules(sg.SecurityGroupRules))
		resultRules := schema.NewSet(secGroupUniqueID, rulesRaw)

		if err := d.Set("security_group_rules", resultRules); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Println("[DEBUG] Finish k8s cluster reading")
	return diags
}

func resourceK8sV2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start k8s cluster updating")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, K8sPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	tasksClient, err := CreateClient(provider, d, tasksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterName := d.Get("name").(string)

	if d.HasChange("version") {
		if err := resourceK8sV2UgradeCluster(client, tasksClient, clusterName, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("authentication", "autoscaler_config", "cni", "ddos_profile") {
		if err := resourceK8sV2UpdateCluster(client, tasksClient, clusterName, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("pool") {
		// 1 pool   => Allow in-place updates and add/delete, but return error on replace.
		//             Users must create a new pool with different name in such case.
		// 2+ pools => Allow all operations, but make sure we don't end up with 0 pools at any moment.
		//             This means we process each pool change one-by-one, and perform adds before deletes.
		o, n := d.GetChange("pool")
		old, new := o.([]interface{}), n.([]interface{})

		if err := resourceK8sV2CheckLimits(client, old, new); err != nil {
			return diag.FromErr(err)
		}

		// Any new pools must be created first, so that "replace" can safely delete pools that it will recreate.
		// This also covers pools that were renamed, because pool name must be unique.
		for _, pool := range new {
			if resourceK8sV2FindClusterPool(old, pool) == nil {
				if err := resourceK8sV2CreateClusterPool(client, tasksClient, clusterName, pool); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		// Check replaces before updates, because replace due to its nature results in all fields being updated.
		for _, pool := range new {
			if resourceK8sV2ClusterPoolNeedsReplace(old, pool) {
				if len(old) == 1 && len(new) == 1 {
					msg := "cannot replace the only pool in the cluster, please create a new pool with different name and remove this one"
					return diag.FromErr(fmt.Errorf("%s", msg))
				}
				if err := resourceK8sV2DeleteClusterPool(client, tasksClient, clusterName, pool); err != nil {
					return diag.FromErr(err)
				}
				if err := resourceK8sV2CreateClusterPool(client, tasksClient, clusterName, pool); err != nil {
					return diag.FromErr(err)
				}
			} else if resourceK8sV2ClusterPoolNeedsUpdate(old, pool) {
				if err := resourceK8sV2UpdateClusterPool(client, clusterName, pool); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		// Finish up by removing all pools that need to be deleted (explicit deletes and leftovers from renames).
		// This allows us to have replace working in case we are going down to 1 pool.
		for _, pool := range old {
			if resourceK8sV2FindClusterPool(new, pool) == nil {
				if err := resourceK8sV2DeleteClusterPool(client, tasksClient, clusterName, pool); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("security_group_rules") {
		o, n := d.GetChange("security_group_rules")
		newUsersRules := n.(*schema.Set)
		oldUsersRules := o.(*schema.Set)

		sgClient, err := CreateClient(provider, d, securityGroupPoint, versionPointV1)
		if err != nil {
			return diag.FromErr(err)
		}

		updateRequest := securitygroups.UpdateOpts{}
		for _, rule := range oldUsersRules.List() {
			r := rule.(map[string]interface{})
			updateRequest.ChangedRules = append(updateRequest.ChangedRules, securitygroups.UpdateSecurityGroupRuleOpts{
				Action:              types.ActionDelete,
				SecurityGroupRuleID: r["id"].(string),
			})
		}

		for _, rule := range newUsersRules.List() {
			r := rule.(map[string]interface{})
			changedRules := securitygroups.UpdateSecurityGroupRuleOpts{
				Action:    types.ActionCreate,
				Direction: types.RuleDirection(r["direction"].(string)),
				EtherType: types.EtherType(r["ethertype"].(string)),
				Protocol:  types.Protocol(r["protocol"].(string)),
			}

			if port := r["port_range_max"].(int); port != 0 {
				changedRules.PortRangeMax = &port
			}
			if port := r["port_range_min"].(int); port != 0 {
				changedRules.PortRangeMin = &port
			}
			if descr := r["description"].(string); descr != "" {
				changedRules.Description = &descr
			}
			if remoteIPPrefix := r["remote_ip_prefix"].(string); remoteIPPrefix != "" {
				changedRules.RemoteIPPrefix = &remoteIPPrefix
			}

			updateRequest.ChangedRules = append(updateRequest.ChangedRules, changedRules)
		}

		_, err = securitygroups.Update(sgClient, d.Get("security_group_id").(string), updateRequest).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

	}

	diags := resourceK8sV2Read(ctx, d, m)
	log.Printf("[DEBUG] Finish k8s cluster updating (%s)", clusterName)
	return diags
}

func resourceK8sV2UgradeCluster(client, tasksClient *gcorecloud.ServiceClient, clusterName string, d *schema.ResourceData) error {
	upgradeOpts := clusters.UpgradeOpts{
		Version: d.Get("version").(string),
	}
	results, err := clusters.Upgrade(client, clusterName, upgradeOpts).Extract()
	if err != nil {
		return fmt.Errorf("upgrade cluster: %w", err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		return nil, nil
	})
	if err != nil {
		return fmt.Errorf("wait for task %s: %w", taskID, err)
	}
	return nil
}

func resourceK8sV2UpdateCluster(client, tasksClient *gcorecloud.ServiceClient, clusterName string, d *schema.ResourceData) error {
	opts := clusters.UpdateOpts{}

	if d.HasChange("authentication") {
		if authI, ok := d.GetOk("authentication"); ok {
			authA := authI.([]interface{})
			auth := authA[0].(map[string]interface{})
			opts.Authentication = &clusters.AuthenticationCreateOpts{}
			if oidcI, ok := auth["oidc"]; ok {
				oidcA := oidcI.([]interface{})
				oidc := oidcA[0].(map[string]interface{})
				if len(oidc) > 0 {
					opts.Authentication.OIDC = &clusters.OIDCCreateOpts{
						ClientID:       oidc["client_id"].(string),
						GroupsClaim:    oidc["groups_claim"].(string),
						GroupsPrefix:   oidc["groups_prefix"].(string),
						IssuerURL:      oidc["issuer_url"].(string),
						UsernameClaim:  oidc["username_claim"].(string),
						UsernamePrefix: oidc["username_prefix"].(string),
					}
					if len(oidc["required_claims"].(map[string]interface{})) > 0 {
						opts.Authentication.OIDC.RequiredClaims = map[string]string{}
						for k, v := range oidc["required_claims"].(map[string]interface{}) {
							opts.Authentication.OIDC.RequiredClaims[k] = v.(string)
						}
					}
					if algs, ok := oidc["signing_algs"].(*schema.Set); ok {
						for _, alg := range algs.List() {
							opts.Authentication.OIDC.SigningAlgs = append(opts.Authentication.OIDC.SigningAlgs, alg.(string))
						}
					}
				}
			}
		}
	}

	if d.HasChange("autoscaler_config") {
		if autoscalerCfgI, ok := d.GetOk("autoscaler_config"); ok {
			autoscalerCfg := autoscalerCfgI.(map[string]interface{})
			opts.AutoscalerConfig = map[string]string{}
			for k, v := range autoscalerCfg {
				opts.AutoscalerConfig[k] = v.(string)
			}
		}
	}

	if d.HasChange("cni") {
		if cniI, ok := d.GetOk("cni"); ok {
			cniA := cniI.([]interface{})
			cni := cniA[0].(map[string]interface{})
			opts.CNI = &clusters.CNICreateOpts{Provider: clusters.CNIProvider(cni["provider"].(string))}
			if cni["provider"].(string) == "cilium" {
				if ciliumI, ok := cni["cilium"]; ok {
					ciliumA := ciliumI.([]interface{})
					if len(ciliumA) != 0 {
						cilium := ciliumA[0].(map[string]interface{})
						opts.CNI.Cilium = &clusters.CiliumCreateOpts{
							MaskSize:                 cilium["mask_size"].(int),
							MaskSizeV6:               cilium["mask_size_v6"].(int),
							Tunnel:                   clusters.TunnelType(cilium["tunnel"].(string)),
							Encryption:               cilium["encryption"].(bool),
							LoadBalancerMode:         clusters.LBModeType(cilium["lb_mode"].(string)),
							LoadBalancerAcceleration: cilium["lb_acceleration"].(bool),
							RoutingMode:              clusters.RoutingModeType(cilium["routing_mode"].(string)),
							HubbleRelay:              cilium["hubble_relay"].(bool),
							HubbleUI:                 cilium["hubble_ui"].(bool),
						}
					}
				}
			}
		}
	}

	if d.HasChange("ddos_profile") {
		if ddosProfileI, ok := d.GetOk("ddos_profile"); ok {
			ddosProfile := ddosProfileI.([]interface{})
			if len(ddosProfile) != 0 {
				profile := ddosProfile[0].(map[string]interface{})
				opts.DDoSProfile = &clusters.DDoSProfileCreateOpts{
					Enabled: profile["enabled"].(bool),
				}
				if fields, ok := profile["fields"].([]interface{}); ok {
					if len(fields) != 0 {
						opts.DDoSProfile.Fields = make([]clusters.DDoSProfileField, len(fields))
						for i, field := range fields {
							fieldMap := field.(map[string]interface{})
							opts.DDoSProfile.Fields[i] = clusters.DDoSProfileField{
								BaseField: fieldMap["base_field"].(int),
							}
							if value, ok := fieldMap["value"].(string); ok {
								opts.DDoSProfile.Fields[i].Value = &value
							}
							if fieldValue, ok := fieldMap["field_value"].(string); ok {
								opts.DDoSProfile.Fields[i].FieldValue = &fieldValue
							}
						}
					}
				}
				if profileTemplate, ok := profile["profile_template"].(int); ok {
					opts.DDoSProfile.ProfileTemplate = &profileTemplate
				}
				if profileTemplateName, ok := profile["profile_template_name"].(string); ok {
					opts.DDoSProfile.ProfileTemplateName = &profileTemplateName
				}
			}
		}
	}

	results, err := clusters.Update(client, clusterName, opts).Extract()
	if err != nil {
		return fmt.Errorf("update cluster: %w", err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		return nil, nil
	})
	if err != nil {
		return fmt.Errorf("wait for task %s: %w", taskID, err)
	}

	return nil
}

func resourceK8sV2Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start k8s cluster deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, K8sPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterName := d.Get("name").(string)
	results, err := clusters.Delete(client, clusterName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)

	tasksClient, err := CreateClient(provider, d, tasksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := clusters.Get(client, clusterName).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete k8s cluster with name: %s", clusterName)
		}
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Finish k8s cluster deleting")
	return diags
}

func resourceK8sV2FindClusterPool(list []interface{}, pool interface{}) interface{} {
	if _, ok := pool.(map[string]interface{}); !ok {
		return nil
	}
	for _, item := range list {
		if _, ok := item.(map[string]interface{}); !ok {
			continue
		}
		if item.(map[string]interface{})["name"] == pool.(map[string]interface{})["name"] {
			return item
		}
	}
	return nil
}

func resourceK8sV2ClusterPoolNeedsUpdate(list []interface{}, pool interface{}) bool {
	found := resourceK8sV2FindClusterPool(list, pool)
	if found == nil {
		return false // adding new pool is not an update
	}
	old, new := found.(map[string]interface{}), pool.(map[string]interface{})
	if old["min_node_count"] != new["min_node_count"] {
		return true
	}
	if old["max_node_count"] != new["max_node_count"] {
		return true
	}
	if old["auto_healing_enabled"] != new["auto_healing_enabled"] {
		return true
	}
	if !reflect.DeepEqual(old["labels"], new["labels"]) {
		return true
	}
	if !reflect.DeepEqual(old["taints"], new["taints"]) {
		return true
	}
	return false
}

func resourceK8sV2ClusterPoolNeedsReplace(list []interface{}, pool interface{}) bool {
	found := resourceK8sV2FindClusterPool(list, pool)
	if found == nil {
		return false // adding new pool is not a replace
	}
	old, new := found.(map[string]interface{}), pool.(map[string]interface{})
	if old["flavor_id"] != new["flavor_id"] {
		return true
	}
	if old["boot_volume_type"] != new["boot_volume_type"] {
		return true
	}
	if old["boot_volume_size"] != new["boot_volume_size"] {
		return true
	}
	if old["is_public_ipv4"] != new["is_public_ipv4"] {
		return true
	}
	if old["servergroup_policy"] != new["servergroup_policy"] {
		return true
	}
	if !reflect.DeepEqual(old["crio_config"], new["crio_config"]) {
		return true
	}
	if !reflect.DeepEqual(old["kubelet_config"], new["kubelet_config"]) {
		return true
	}
	return false
}

func resourceK8sV2CreateClusterPool(client, tasksClient *gcorecloud.ServiceClient, clusterName string, data interface{}) error {
	pool := data.(map[string]interface{})
	poolName := pool["name"].(string)
	log.Printf("[DEBUG] Creating cluster pool (%s)", poolName)

	opts := pools.CreateOpts{
		Name:               pool["name"].(string),
		FlavorID:           pool["flavor_id"].(string),
		MinNodeCount:       pool["min_node_count"].(int),
		MaxNodeCount:       pool["max_node_count"].(int),
		BootVolumeSize:     pool["boot_volume_size"].(int),
		BootVolumeType:     volumes.VolumeType(pool["boot_volume_type"].(string)),
		AutoHealingEnabled: pool["auto_healing_enabled"].(bool),
		ServerGroupPolicy:  servergroups.ServerGroupPolicy(pool["servergroup_policy"].(string)),
		IsPublicIPv4:       pool["is_public_ipv4"].(bool),
	}
	if labels, ok := pool["labels"].(map[string]interface{}); ok {
		opts.Labels = map[string]string{}
		for k, v := range labels {
			opts.Labels[k] = v.(string)
		}
	}
	if taints, ok := pool["taints"].(map[string]interface{}); ok {
		opts.Taints = map[string]string{}
		for k, v := range taints {
			opts.Taints[k] = v.(string)
		}
	}
	if crioCfg, ok := pool["crio_config"].(map[string]interface{}); ok {
		opts.CrioConfig = map[string]string{}
		for k, v := range crioCfg {
			opts.CrioConfig[k] = v.(string)
		}
	}
	if kubeletCfg, ok := pool["kubelet_config"].(map[string]interface{}); ok {
		opts.KubeletConfig = map[string]string{}
		for k, v := range kubeletCfg {
			opts.KubeletConfig[k] = v.(string)
		}
	}
	results, err := pools.Create(client, clusterName, opts).Extract()
	if err != nil {
		return fmt.Errorf("create cluster pool: %w", err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		return nil, nil
	})
	if err != nil {
		return fmt.Errorf("wait for task %s: %w", taskID, err)
	}

	log.Printf("[DEBUG] Created cluster pool (%s)", poolName)
	return nil
}

func resourceK8sV2DeleteClusterPool(client, tasksClient *gcorecloud.ServiceClient, clusterName string, data interface{}) error {
	pool, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	poolName, ok := pool["name"].(string)
	if !ok || poolName == "" {
		return nil
	}

	log.Printf("[DEBUG] Deleting cluster pool (%s)", poolName)
	results, err := pools.Delete(client, clusterName, poolName).Extract()
	if err != nil {
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil
		default:
			return fmt.Errorf("delete cluster pool: %w", err)
		}
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		return nil, nil
	})
	if err != nil {
		return fmt.Errorf("wait for task %s: %w", taskID, err)
	}

	log.Printf("[DEBUG] Deleted cluster pool (%s)", poolName)
	return nil
}

func resourceK8sV2UpdateClusterPool(client *gcorecloud.ServiceClient, clusterName string, data interface{}) error {
	pool := data.(map[string]interface{})
	poolName := pool["name"].(string)
	log.Printf("[DEBUG] Updating cluster pool (%s)", poolName)

	opts := pools.UpdateOpts{
		MinNodeCount: pool["min_node_count"].(int),
		MaxNodeCount: pool["max_node_count"].(int),
	}
	if v, ok := pool["auto_healing_enabled"].(bool); ok {
		opts.AutoHealingEnabled = &v
	}
	if labels, ok := pool["labels"].(map[string]interface{}); ok && len(labels) > 0 {
		result := map[string]string{}
		for k, v := range labels {
			result[k] = v.(string)
		}
		opts.Labels = &result
	}
	if taints, ok := pool["taints"].(map[string]interface{}); ok && len(taints) > 0 {
		result := map[string]string{}
		for k, v := range taints {
			result[k] = v.(string)
		}
		opts.Taints = &result
	}
	_, err := pools.Update(client, clusterName, poolName, opts).Extract()
	if err != nil {
		return fmt.Errorf("update cluster pool: %w", err)
	}

	log.Printf("[DEBUG] Updated cluster pool (%s)", poolName)
	return nil
}

func resourceK8sV2PoolDataFromPool(pool pools.ClusterPool) interface{} {
	return map[string]interface{}{
		"name":                 pool.Name,
		"flavor_id":            pool.FlavorID,
		"min_node_count":       pool.MinNodeCount,
		"max_node_count":       pool.MaxNodeCount,
		"node_count":           pool.NodeCount,
		"boot_volume_type":     pool.BootVolumeType.String(),
		"boot_volume_size":     pool.BootVolumeSize,
		"auto_healing_enabled": pool.AutoHealingEnabled,
		"is_public_ipv4":       pool.IsPublicIPv4,
		"labels":               resourceK8sV2FilteredPoolLabels(pool.Labels),
		"taints":               pool.Taints,
		"crio_config":          pool.CrioConfig,
		"kubelet_config":       pool.KubeletConfig,
		"servergroup_policy":   pool.ServerGroupPolicy,
		"servergroup_name":     pool.ServerGroupName,
		"servergroup_id":       pool.ServerGroupID,
		"status":               pool.Status,
		"created_at":           pool.CreatedAt.Format(time.RFC850),
	}
}

func resourceK8sV2FilteredPoolLabels(labels map[string]string) map[string]string {
	result := map[string]string{}
	for k, v := range labels {
		// filter out system labels to hide them from state file and diffs
		if strings.HasPrefix(k, "gcorecluster.x-k8s.io") {
			continue
		}
		result[k] = v
	}
	return result
}

func resourceK8sV2IsVMFlavor(flavor string) bool {
	return strings.HasPrefix(flavor, "g") || strings.HasPrefix(flavor, "a")
}

func resourceK8sV2CheckLimits(client *gcorecloud.ServiceClient, old, new []interface{}) error {
	log.Printf("[DEBUG] Checking quota limits")

	opts := &clusters.CheckLimitsOpts{}
	for _, n := range new {
		newPool, ok := n.(map[string]interface{})
		if !ok || len(newPool) == 0 {
			continue
		}
		if resourceK8sV2FindClusterPool(old, newPool) == nil || resourceK8sV2ClusterPoolNeedsReplace(old, newPool) {
			poolOpts := clusters.CheckLimitsPoolOpts{
				Name:              newPool["name"].(string),
				FlavorID:          newPool["flavor_id"].(string),
				MinNodeCount:      newPool["min_node_count"].(int),
				MaxNodeCount:      newPool["max_node_count"].(int),
				BootVolumeSize:    newPool["boot_volume_size"].(int),
				ServerGroupPolicy: servergroups.ServerGroupPolicy(newPool["servergroup_policy"].(string)),
			}
			opts.Pools = append(opts.Pools, poolOpts)
		} else if resourceK8sV2ClusterPoolNeedsUpdate(old, newPool) {
			oldPool := resourceK8sV2FindClusterPool(old, newPool).(map[string]interface{})
			minCount := newPool["min_node_count"].(int) - oldPool["min_node_count"].(int)
			maxCount := newPool["max_node_count"].(int) - oldPool["max_node_count"].(int)
			if minCount <= 0 {
				continue
			}
			poolOpts := clusters.CheckLimitsPoolOpts{
				Name:              newPool["name"].(string),
				FlavorID:          newPool["flavor_id"].(string),
				MinNodeCount:      minCount,
				MaxNodeCount:      maxCount,
				BootVolumeSize:    newPool["boot_volume_size"].(int),
				ServerGroupPolicy: servergroups.ServerGroupPolicy(newPool["servergroup_policy"].(string)),
			}
			opts.Pools = append(opts.Pools, poolOpts)
		}
	}

	if len(opts.Pools) > 0 {
		quota, err := clusters.CheckLimits(client, opts).Extract()
		if err != nil {
			return fmt.Errorf("check limits: %w", err)
		}
		if len(*quota) > 0 {
			b, _ := json.Marshal(quota) // pretty print map
			return fmt.Errorf("quota limits exceeded for this operation: %s", string(b))
		}
	}
	return nil
}
