// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_ai_task

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/streaming"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingAITasksResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[StreamingAITasksItemsDataSourceModel] `json:"results,computed"`
}

type StreamingAITasksDataSourceModel struct {
	DateCreated types.String                                                       `tfsdk:"date_created" query:"date_created,optional"`
	Ordering    types.String                                                       `tfsdk:"ordering" query:"ordering,optional"`
	Search      types.String                                                       `tfsdk:"search" query:"search,optional"`
	Status      types.String                                                       `tfsdk:"status" query:"status,optional"`
	TaskID      types.String                                                       `tfsdk:"task_id" query:"task_id,optional"`
	TaskName    types.String                                                       `tfsdk:"task_name" query:"task_name,optional"`
	Limit       types.Int64                                                        `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems    types.Int64                                                        `tfsdk:"max_items"`
	Items       customfield.NestedObjectList[StreamingAITasksItemsDataSourceModel] `tfsdk:"items"`
}

func (m *StreamingAITasksDataSourceModel) toListParams(_ context.Context) (params streaming.AITaskListParams, diags diag.Diagnostics) {
	params = streaming.AITaskListParams{}

	if !m.DateCreated.IsNull() {
		params.DateCreated = param.NewOpt(m.DateCreated.ValueString())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = streaming.AITaskListParamsOrdering(m.Ordering.ValueString())
	}
	if !m.Search.IsNull() {
		params.Search = param.NewOpt(m.Search.ValueString())
	}
	if !m.Status.IsNull() {
		params.Status = streaming.AITaskListParamsStatus(m.Status.ValueString())
	}
	if !m.TaskID.IsNull() {
		params.TaskID = param.NewOpt(m.TaskID.ValueString())
	}
	if !m.TaskName.IsNull() {
		params.TaskName = streaming.AITaskListParamsTaskName(m.TaskName.ValueString())
	}

	return
}

type StreamingAITasksItemsDataSourceModel struct {
	Progress types.Int64                                                       `tfsdk:"progress" json:"progress,computed"`
	Status   types.String                                                      `tfsdk:"status" json:"status,computed"`
	TaskData customfield.NestedObject[StreamingAITasksTaskDataDataSourceModel] `tfsdk:"task_data" json:"task_data,computed"`
	TaskID   types.String                                                      `tfsdk:"task_id" json:"task_id,computed"`
	TaskName types.String                                                      `tfsdk:"task_name" json:"task_name,computed"`
}

type StreamingAITasksTaskDataDataSourceModel struct {
	TaskName          types.String `tfsdk:"task_name" json:"task_name,computed"`
	URL               types.String `tfsdk:"url" json:"url,computed"`
	AudioLanguage     types.String `tfsdk:"audio_language" json:"audio_language,computed"`
	ClientEntityData  types.String `tfsdk:"client_entity_data" json:"client_entity_data,computed"`
	ClientUserID      types.String `tfsdk:"client_user_id" json:"client_user_id,computed"`
	SubtitlesLanguage types.String `tfsdk:"subtitles_language" json:"subtitles_language,computed"`
	Category          types.String `tfsdk:"category" json:"category,computed"`
	StopObjects       types.String `tfsdk:"stop_objects" json:"stop_objects,computed"`
}
