// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry_artifact

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudRegistryArtifactsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudRegistryArtifactsItemsDataSourceModel] `json:"results,computed"`
}

type CloudRegistryArtifactsDataSourceModel struct {
	RegistryID     types.Int64                                                              `tfsdk:"registry_id" path:"registry_id,required"`
	RepositoryName types.String                                                             `tfsdk:"repository_name" path:"repository_name,required"`
	ProjectID      types.Int64                                                              `tfsdk:"project_id" path:"project_id,optional"`
	RegionID       types.Int64                                                              `tfsdk:"region_id" path:"region_id,optional"`
	MaxItems       types.Int64                                                              `tfsdk:"max_items"`
	Items          customfield.NestedObjectList[CloudRegistryArtifactsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudRegistryArtifactsDataSourceModel) toListParams(_ context.Context) (params cloud.RegistryArtifactListParams, diags diag.Diagnostics) {
	params = cloud.RegistryArtifactListParams{
		RegistryID: m.RegistryID.ValueInt64(),
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudRegistryArtifactsItemsDataSourceModel struct {
	ID           types.Int64                                                             `tfsdk:"id" json:"id,computed"`
	Digest       types.String                                                            `tfsdk:"digest" json:"digest,computed"`
	PulledAt     timetypes.RFC3339                                                       `tfsdk:"pulled_at" json:"pulled_at,computed" format:"date-time"`
	PushedAt     timetypes.RFC3339                                                       `tfsdk:"pushed_at" json:"pushed_at,computed" format:"date-time"`
	RegistryID   types.Int64                                                             `tfsdk:"registry_id" json:"registry_id,computed"`
	RepositoryID types.Int64                                                             `tfsdk:"repository_id" json:"repository_id,computed"`
	Size         types.Int64                                                             `tfsdk:"size" json:"size,computed"`
	Tags         customfield.NestedObjectList[CloudRegistryArtifactsTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
}

type CloudRegistryArtifactsTagsDataSourceModel struct {
	ID           types.Int64       `tfsdk:"id" json:"id,computed"`
	ArtifactID   types.Int64       `tfsdk:"artifact_id" json:"artifact_id,computed"`
	Name         types.String      `tfsdk:"name" json:"name,computed"`
	PulledAt     timetypes.RFC3339 `tfsdk:"pulled_at" json:"pulled_at,computed" format:"date-time"`
	PushedAt     timetypes.RFC3339 `tfsdk:"pushed_at" json:"pushed_at,computed" format:"date-time"`
	RepositoryID types.Int64       `tfsdk:"repository_id" json:"repository_id,computed"`
}
