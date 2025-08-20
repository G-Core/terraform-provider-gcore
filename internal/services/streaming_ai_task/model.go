// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_ai_task

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingAITaskModel struct {
	TaskID            types.String                                                 `tfsdk:"task_id" path:"task_id,optional"`
	TaskName          types.String                                                 `tfsdk:"task_name" json:"task_name,required"`
	URL               types.String                                                 `tfsdk:"url" json:"url,required,no_refresh"`
	AudioLanguage     types.String                                                 `tfsdk:"audio_language" json:"audio_language,optional,no_refresh"`
	Category          types.String                                                 `tfsdk:"category" json:"category,optional,no_refresh"`
	ClientEntityData  types.String                                                 `tfsdk:"client_entity_data" json:"client_entity_data,optional,no_refresh"`
	ClientUserID      types.String                                                 `tfsdk:"client_user_id" json:"client_user_id,optional,no_refresh"`
	SubtitlesLanguage types.String                                                 `tfsdk:"subtitles_language" json:"subtitles_language,optional,no_refresh"`
	Progress          types.Int64                                                  `tfsdk:"progress" json:"progress,computed"`
	Status            types.String                                                 `tfsdk:"status" json:"status,computed"`
	ProcessingTime    customfield.NestedObject[StreamingAITaskProcessingTimeModel] `tfsdk:"processing_time" json:"processing_time,computed"`
	Result            customfield.NestedObject[StreamingAITaskResultModel]         `tfsdk:"result" json:"result,computed"`
	TaskData          customfield.NestedObject[StreamingAITaskTaskDataModel]       `tfsdk:"task_data" json:"task_data,computed"`
}

func (m StreamingAITaskModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StreamingAITaskModel) MarshalJSONForUpdate(state StreamingAITaskModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type StreamingAITaskProcessingTimeModel struct {
	CompletedAt  types.String  `tfsdk:"completed_at" json:"completed_at,computed"`
	StartedAt    types.String  `tfsdk:"started_at" json:"started_at,computed"`
	TotalTimeSec types.Float64 `tfsdk:"total_time_sec" json:"total_time_sec,computed"`
}

type StreamingAITaskResultModel struct {
	ConcatenatedText         types.String                                                      `tfsdk:"concatenated_text" json:"concatenated_text,computed"`
	Languages                customfield.List[types.String]                                    `tfsdk:"languages" json:"languages,computed"`
	SpeechDetected           types.Bool                                                        `tfsdk:"speech_detected" json:"speech_detected,computed"`
	Subtitles                customfield.NestedObjectList[StreamingAITaskResultSubtitlesModel] `tfsdk:"subtitles" json:"subtitles,computed"`
	VttContent               types.String                                                      `tfsdk:"vtt_content" json:"vttContent,computed"`
	DetectionResults         customfield.List[types.String]                                    `tfsdk:"detection_results" json:"detection_results,computed"`
	Frames                   customfield.NestedObjectList[StreamingAITaskResultFramesModel]    `tfsdk:"frames" json:"frames,computed"`
	SportDetected            types.Bool                                                        `tfsdk:"sport_detected" json:"sport_detected,computed"`
	WeaponDetected           types.Bool                                                        `tfsdk:"weapon_detected" json:"weapon_detected,computed"`
	NsfwDetected             types.Bool                                                        `tfsdk:"nsfw_detected" json:"nsfw_detected,computed"`
	PornDetected             types.Bool                                                        `tfsdk:"porn_detected" json:"porn_detected,computed"`
	ChildPornographyDetected types.Bool                                                        `tfsdk:"child_pornography_detected" json:"child_pornography_detected,computed"`
	Error                    types.String                                                      `tfsdk:"error" json:"error,computed"`
}

type StreamingAITaskResultSubtitlesModel struct {
	EndTime   types.String `tfsdk:"end_time" json:"end_time,computed"`
	StartTime types.String `tfsdk:"start_time" json:"start_time,computed"`
	Text      types.String `tfsdk:"text" json:"text,computed"`
}

type StreamingAITaskResultFramesModel struct {
	Confidence  types.Float64 `tfsdk:"confidence" json:"confidence,computed"`
	FrameNumber types.Int64   `tfsdk:"frame_number" json:"frame-number,computed"`
	Label       types.String  `tfsdk:"label" json:"label,computed"`
}

type StreamingAITaskTaskDataModel struct {
	TaskName          types.String `tfsdk:"task_name" json:"task_name,computed"`
	URL               types.String `tfsdk:"url" json:"url,computed"`
	AudioLanguage     types.String `tfsdk:"audio_language" json:"audio_language,computed"`
	ClientEntityData  types.String `tfsdk:"client_entity_data" json:"client_entity_data,computed"`
	ClientUserID      types.String `tfsdk:"client_user_id" json:"client_user_id,computed"`
	SubtitlesLanguage types.String `tfsdk:"subtitles_language" json:"subtitles_language,computed"`
	Category          types.String `tfsdk:"category" json:"category,computed"`
	StopObjects       types.String `tfsdk:"stop_objects" json:"stop_objects,computed"`
}
