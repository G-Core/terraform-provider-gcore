// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_ai_task

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type StreamingAITaskDataSourceModel struct {
	TaskID         types.String                                                           `tfsdk:"task_id" path:"task_id,required"`
	Progress       types.Int64                                                            `tfsdk:"progress" json:"progress,computed"`
	Status         types.String                                                           `tfsdk:"status" json:"status,computed"`
	TaskName       types.String                                                           `tfsdk:"task_name" json:"task_name,computed"`
	ProcessingTime customfield.NestedObject[StreamingAITaskProcessingTimeDataSourceModel] `tfsdk:"processing_time" json:"processing_time,computed"`
	Result         customfield.NestedObject[StreamingAITaskResultDataSourceModel]         `tfsdk:"result" json:"result,computed"`
	TaskData       customfield.NestedObject[StreamingAITaskTaskDataDataSourceModel]       `tfsdk:"task_data" json:"task_data,computed"`
}

type StreamingAITaskProcessingTimeDataSourceModel struct {
	CompletedAt  types.String  `tfsdk:"completed_at" json:"completed_at,computed"`
	StartedAt    types.String  `tfsdk:"started_at" json:"started_at,computed"`
	TotalTimeSec types.Float64 `tfsdk:"total_time_sec" json:"total_time_sec,computed"`
}

type StreamingAITaskResultDataSourceModel struct {
	ConcatenatedText         types.String                                                                `tfsdk:"concatenated_text" json:"concatenated_text,computed"`
	Languages                customfield.List[types.String]                                              `tfsdk:"languages" json:"languages,computed"`
	SpeechDetected           types.Bool                                                                  `tfsdk:"speech_detected" json:"speech_detected,computed"`
	Subtitles                customfield.NestedObjectList[StreamingAITaskResultSubtitlesDataSourceModel] `tfsdk:"subtitles" json:"subtitles,computed"`
	VttContent               types.String                                                                `tfsdk:"vtt_content" json:"vttContent,computed"`
	DetectionResults         customfield.List[types.String]                                              `tfsdk:"detection_results" json:"detection_results,computed"`
	Frames                   customfield.NestedObjectList[StreamingAITaskResultFramesDataSourceModel]    `tfsdk:"frames" json:"frames,computed"`
	SportDetected            types.Bool                                                                  `tfsdk:"sport_detected" json:"sport_detected,computed"`
	WeaponDetected           types.Bool                                                                  `tfsdk:"weapon_detected" json:"weapon_detected,computed"`
	NsfwDetected             types.Bool                                                                  `tfsdk:"nsfw_detected" json:"nsfw_detected,computed"`
	PornDetected             types.Bool                                                                  `tfsdk:"porn_detected" json:"porn_detected,computed"`
	ChildPornographyDetected types.Bool                                                                  `tfsdk:"child_pornography_detected" json:"child_pornography_detected,computed"`
	Error                    types.String                                                                `tfsdk:"error" json:"error,computed"`
}

type StreamingAITaskResultSubtitlesDataSourceModel struct {
	EndTime   types.String `tfsdk:"end_time" json:"end_time,computed"`
	StartTime types.String `tfsdk:"start_time" json:"start_time,computed"`
	Text      types.String `tfsdk:"text" json:"text,computed"`
}

type StreamingAITaskResultFramesDataSourceModel struct {
	Confidence  types.Float64 `tfsdk:"confidence" json:"confidence,computed"`
	FrameNumber types.Int64   `tfsdk:"frame_number" json:"frame-number,computed"`
	Label       types.String  `tfsdk:"label" json:"label,computed"`
}

type StreamingAITaskTaskDataDataSourceModel struct {
	TaskName          types.String `tfsdk:"task_name" json:"task_name,computed"`
	URL               types.String `tfsdk:"url" json:"url,computed"`
	AudioLanguage     types.String `tfsdk:"audio_language" json:"audio_language,computed"`
	ClientEntityData  types.String `tfsdk:"client_entity_data" json:"client_entity_data,computed"`
	ClientUserID      types.String `tfsdk:"client_user_id" json:"client_user_id,computed"`
	SubtitlesLanguage types.String `tfsdk:"subtitles_language" json:"subtitles_language,computed"`
	Category          types.String `tfsdk:"category" json:"category,computed"`
	StopObjects       types.String `tfsdk:"stop_objects" json:"stop_objects,computed"`
}
