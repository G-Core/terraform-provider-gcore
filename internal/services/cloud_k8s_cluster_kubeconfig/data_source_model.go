// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_cluster_kubeconfig

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudK8SClusterKubeconfigDataSourceModel struct {
	ClusterName          types.String      `tfsdk:"cluster_name" path:"cluster_name,required"`
	ProjectID            types.Int64       `tfsdk:"project_id" path:"project_id,optional"`
	RegionID             types.Int64       `tfsdk:"region_id" path:"region_id,optional"`
	ClientCertificate    types.String      `tfsdk:"client_certificate" json:"client_certificate,computed"`
	ClientKey            types.String      `tfsdk:"client_key" json:"client_key,computed"`
	ClusterCaCertificate types.String      `tfsdk:"cluster_ca_certificate" json:"cluster_ca_certificate,computed"`
	Config               types.String      `tfsdk:"config" json:"config,computed"`
	CreatedAt            timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	ExpiresAt            timetypes.RFC3339 `tfsdk:"expires_at" json:"expires_at,computed" format:"date-time"`
	Host                 types.String      `tfsdk:"host" json:"host,computed"`
}

func (m *CloudK8SClusterKubeconfigDataSourceModel) toReadParams(_ context.Context) (params cloud.K8SClusterKubeconfigGetParams, diags diag.Diagnostics) {
	params = cloud.K8SClusterKubeconfigGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}
