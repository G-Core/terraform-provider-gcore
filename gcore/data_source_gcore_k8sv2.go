package gcore

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceK8sV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceK8sV2Read,
		Description: "Represent k8s cluster with one default pool.",
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
				Type:     schema.TypeString,
				Required: true,
			},
			"authentication": {
				Type:        schema.TypeList,
				Description: "Cluster authentication configuration.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oidc": {
							Type:        schema.TypeList,
							Description: "OpenID Connect configuration settings.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client_id": {
										Type:        schema.TypeString,
										Description: "A client id that all tokens must be issued for.",
										Computed:    true,
									},
									"groups_claim": {
										Type:        schema.TypeString,
										Description: "JWT claim to use as the user's group.",
										Computed:    true,
									},
									"groups_prefix": {
										Type:        schema.TypeString,
										Description: "Prefix prepended to group claims to prevent clashes with existing names.",
										Computed:    true,
									},
									"issuer_url": {
										Type:        schema.TypeString,
										Description: "URL of the provider that allows the API server to discover public signing keys. Only URLs that use the https:// scheme are accepted.",
										Computed:    true,
									},
									"required_claims": {
										Type:        schema.TypeMap,
										Description: "A map describing required claims in the ID Token. Each claim is verified to be present in the ID Token with a matching value.",
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"signing_algs": {
										Type:        schema.TypeSet,
										Description: "Accepted signing algorithms. Supported values are: RS256, RS384, RS512, ES256, ES384, ES512, PS256, PS384, PS512.",
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"username_claim": {
										Type:        schema.TypeString,
										Description: "JWT claim to use as the user name. When not specified, the `sub` claim will be used.",
										Computed:    true,
									},
									"username_prefix": {
										Type:        schema.TypeString,
										Description: "Prefix prepended to username claims to prevent clashes with existing names.",
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
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cni": {
				Type:        schema.TypeList,
				Description: "Cluster CNI configuration.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:        schema.TypeString,
							Description: "CNI provider used by the cluster. Supported values are: calico, cilium.",
							Computed:    true,
						},
						"cilium": {
							Type:        schema.TypeList,
							Description: "Cilium CNI configuration.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mask_size": {
										Type:        schema.TypeInt,
										Description: "The size allocated from pods_ip_pool CIDR to node.ipam.podCIDRs.",
										Computed:    true,
									},
									"mask_size_v6": {
										Type:        schema.TypeInt,
										Description: "The size allocated from pods_ipv6_pool CIDR to node.ipam.podCIDRs.",
										Computed:    true,
									},
									"tunnel": {
										Type:        schema.TypeString,
										Description: "Tunneling protocol used in tunneling mode and for ad-hoc tunnels.",
										Computed:    true,
									},
									"encryption": {
										Type:        schema.TypeBool,
										Description: "Is transparent network encryption enabled or not.",
										Computed:    true,
									},
									"lb_mode": {
										Type:        schema.TypeString,
										Description: "The operation mode of load balancing for remote backends. Supported values are snat, dsr, hybrid.",
										Computed:    true,
									},
									"lb_acceleration": {
										Type:        schema.TypeBool,
										Description: "Is load balancer acceleration via XDP enabled or not.",
										Computed:    true,
									},
									"routing_mode": {
										Type:        schema.TypeString,
										Description: "Is native-routing mode or tunneling mode enabled.",
										Computed:    true,
									},
									"hubble_relay": {
										Type:        schema.TypeBool,
										Description: "Is Hubble Relay enabled or not.",
										Computed:    true,
									},
									"hubble_ui": {
										Type:        schema.TypeBool,
										Description: "Is Hubble UI enabled or not.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"ddos_profile": {
				Type:        schema.TypeList,
				Description: "DDoS profile configuration.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Indicates if the DDoS profile is enabled.",
							Computed:    true,
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
										Computed:    true,
									},
									"field_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Complex value. Only one of 'value' or 'field_value' must be specified.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Basic type value. Only one of 'value' or 'field_value' must be specified.",
									},
								},
							},
						},
						"profile_template": {
							Type:        schema.TypeInt,
							Description: "Profile template ID.",
							Computed:    true,
						},
						"profile_template_name": {
							Type:        schema.TypeString,
							Description: "Profile template name.",
							Computed:    true,
						},
					},
				},
			},
			"fixed_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pods_ip_pool": {
				Type:        schema.TypeString,
				Description: "Pods IPv4 IP pool in CIDR notation.",
				Computed:    true,
			},
			"services_ip_pool": {
				Type:        schema.TypeString,
				Description: "Services IPv4 IP pool in CIDR notation.",
				Computed:    true,
			},
			"pods_ipv6_pool": {
				Type:        schema.TypeString,
				Description: "Pods IPv6 IP pool in CIDR notation.",
				Computed:    true,
			},
			"services_ipv6_pool": {
				Type:        schema.TypeString,
				Description: "Services IPv6 IP pool in CIDR notation.",
				Computed:    true,
			},
			"keypair": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_ipv6": {
				Type:        schema.TypeBool,
				Description: "Enable public IPv6 address.",
				Computed:    true,
			},
			"pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"boot_volume_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available values are 'standard', 'ssd_hiiops', 'cold', 'ultra'.",
						},
						"boot_volume_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_healing_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_public_ipv4": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"labels": {
							Type:        schema.TypeMap,
							Description: "Labels applied to the cluster pool nodes.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"taints": {
							Type:        schema.TypeMap,
							Description: "Taints applied to the cluster pool nodes.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
						"servergroup_policy": {
							Type:        schema.TypeString,
							Description: "Server group policy: anti-affinity, soft-anti-affinity or affinity",
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
						"status": {
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
			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
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
	}
}

func dataSourceK8sV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start K8s reading")
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
		return diag.FromErr(err)
	}

	d.SetId(cluster.Name)

	d.Set("name", cluster.Name)
	d.Set("fixed_network", cluster.FixedNetwork)
	d.Set("fixed_subnet", cluster.FixedSubnet)
	d.Set("keypair", cluster.KeyPair)
	d.Set("version", cluster.Version)
	d.Set("node_count", cluster.NodeCount)
	d.Set("flavor_id", cluster.FlavorID)
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
		var fields []interface{}
		for _, f := range cluster.DDoSProfile.Fields {
			field := map[string]interface{}{
				"base_field": f.BaseField,
				"value":      f.Value,
			}
			if f.FieldValue != nil {
				b, err := json.Marshal(f.FieldValue)
				if err != nil {
					return diag.FromErr(err)
				}
				field["field_value"] = string(b)
			}
			fields = append(fields, field)
		}
		v := map[string]interface{}{
			"enabled":               cluster.DDoSProfile.Enabled,
			"fields":                fields,
			"profile_template":      cluster.DDoSProfile.ProfileTemplate,
			"profile_template_name": cluster.DDoSProfile.ProfileTemplateName,
		}
		if err := d.Set("ddos_profile", []interface{}{v}); err != nil {
			return diag.FromErr(err)
		}
	}

	var ps []map[string]interface{}
	for _, pool := range cluster.Pools {
		ps = append(ps, map[string]interface{}{
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
		})
	}
	if err := d.Set("pools", ps); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish K8s reading")
	return diags
}
