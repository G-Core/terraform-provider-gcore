// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_cluster_kubeconfig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudK8SClusterKubeconfigDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Kubeconfig provides the necessary configuration and credentials to access a Kubernetes cluster using kubectl or other Kubernetes clients.",
		Attributes: map[string]schema.Attribute{
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
			"client_certificate": schema.StringAttribute{
				Description: "String in base64 format. Cluster client certificate",
				Computed:    true,
			},
			"client_key": schema.StringAttribute{
				Description: "String in base64 format. Cluster client key",
				Computed:    true,
			},
			"cluster_ca_certificate": schema.StringAttribute{
				Description: "String in base64 format. Cluster ca certificate",
				Computed:    true,
			},
			"config": schema.StringAttribute{
				Description: "Cluster kubeconfig",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Kubeconfig creation date",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"expires_at": schema.StringAttribute{
				Description: "Kubeconfig expiration date",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"host": schema.StringAttribute{
				Description: "Cluster host",
				Computed:    true,
			},
		},
	}
}

func (d *CloudK8SClusterKubeconfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudK8SClusterKubeconfigDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
