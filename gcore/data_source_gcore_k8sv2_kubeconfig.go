package gcore

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceK8sV2KubeConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceK8sV2KubeConfigRead,
		Description: "Represent k8s cluster's kubeconfig.",
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
			"cluster_name": {
				Type:        schema.TypeString,
				Description: "Cluster name to fetch kubeconfig",
				Required:    true,
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Description: "Raw kubeconfig file",
				Computed:    true,
			},
			"host": {
				Type:        schema.TypeString,
				Description: "Cluster host",
				Computed:    true,
			},
			"cluster_ca_certificate": {
				Type:        schema.TypeString,
				Description: "String in base64 format. Cluster ca certificate",
				Computed:    true,
			},
			"client_certificate": {
				Type:        schema.TypeString,
				Description: "String in base64 format. Cluster client certificate",
				Computed:    true,
			},
			"client_key": {
				Type:        schema.TypeString,
				Description: "String in base64 format. Cluster client key",
				Computed:    true,
			},
		},
	}
}

func dataSourceK8sV2KubeConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start K8s kubeconfig reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, K8sPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterName := d.Get("cluster_name").(string)
	cluster, err := clusters.Get(client, clusterName).Extract()
	if err != nil {
		return diag.FromErr(fmt.Errorf("cant get cluster: %s", err.Error()))
	}

	kubeconfig, err := clusters.GetConfig(client, clusterName).Extract()
	if err != nil {
		return diag.FromErr(fmt.Errorf("cant get kubeconfig: %s", err.Error()))
	}

	d.SetId(cluster.Name)
	d.Set("kubeconfig", kubeconfig.Config)
	d.Set("cluster_ca_certificate", kubeconfig.ClusterCACertificate)
	d.Set("client_certificate", kubeconfig.ClientCertificate)
	d.Set("client_key", kubeconfig.ClientKey)
	d.Set("host", kubeconfig.Host)

	log.Println("[DEBUG] Finish K8s kubeconfig reading")
	return diags
}
