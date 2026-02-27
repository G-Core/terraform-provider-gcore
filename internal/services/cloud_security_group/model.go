// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudSecurityGroupModel struct {
	ID             types.String                                                `tfsdk:"id" json:"id,computed"`
	ProjectID      types.Int64                                                 `tfsdk:"project_id" path:"project_id,optional"`
	RegionID       types.Int64                                                 `tfsdk:"region_id" path:"region_id,optional"`
	Name           types.String                                                `tfsdk:"name" json:"name,required"`
	Description    types.String                                                `tfsdk:"description" json:"description,computed_optional"`
	Tags           *map[string]types.String                                    `tfsdk:"tags" json:"tags,optional,no_refresh"`
	CreatedAt      timetypes.RFC3339                                           `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Region         types.String                                                `tfsdk:"region" json:"region,computed"`
	RevisionNumber types.Int64                                                 `tfsdk:"revision_number" json:"revision_number,computed"`
	UpdatedAt      timetypes.RFC3339                                           `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	TagsV2         customfield.NestedObjectList[CloudSecurityGroupTagsV2Model] `tfsdk:"tags_v2" json:"tags_v2,computed"`
}

func (m CloudSecurityGroupModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudSecurityGroupModel) MarshalJSONForUpdate(state CloudSecurityGroupModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudSecurityGroupTagsV2Model struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}
