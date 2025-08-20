// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_ai_task

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*StreamingAITasksDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"date_created": schema.StringAttribute{
				Description: "Time when task was created. Datetime in ISO 8601 format.",
				Optional:    true,
			},
			"ordering": schema.StringAttribute{
				Description: "Which field to use when ordering the results: `task_id`, status, and `task_name`.\nSorting is done in ascending (ASC) order.\nIf parameter is omitted then \"`started_at` DESC\" is used for ordering by default.\nAvailable values: \"task_id\", \"status\", \"task_name\", \"started_at\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"task_id",
						"status",
						"task_name",
						"started_at",
					),
				},
			},
			"search": schema.StringAttribute{
				Description: "This is an field for combined text search in the following fields: `task_id`, `task_name`, status, and `task_data`.\nBoth full and partial searches are possible inside specified above fields. For example, you can filter tasks of a certain category, or tasks by a specific original file.\nExample:\n- To filter tasks of Content Moderation NSFW method: ```GET /streaming/ai/tasks?search=nsfw```\n- To filter tasks of processing video from a specific origin: ```GET /streaming/ai/tasks?search=s3.eu-west-1.amazonaws.com```",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Task status\nAvailable values: \"FAILURE\", \"PENDING\", \"RECEIVED\", \"RETRY\", \"REVOKED\", \"STARTED\", \"SUCCESS\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"FAILURE",
						"PENDING",
						"RECEIVED",
						"RETRY",
						"REVOKED",
						"STARTED",
						"SUCCESS",
					),
				},
			},
			"task_id": schema.StringAttribute{
				Description: "The task unique identifier to fiund",
				Optional:    true,
			},
			"task_name": schema.StringAttribute{
				Description: "Type of the AI task. Reflects the original API method that was used to create the AI task.\nAvailable values: \"transcription\", \"content-moderation\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("transcription", "content-moderation"),
				},
			},
			"limit": schema.Int64Attribute{
				Description: "Number of results to return per page.",
				Computed:    true,
				Optional:    true,
			},
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[StreamingAITasksItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"progress": schema.Int64Attribute{
							Description: "Percentage of task completed. A value greater than 0 means that it has been taken into operation and is being processed.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Status of processing the AI task. See GET /ai/results method for description.\nAvailable values: \"PENDING\", \"STARTED\", \"SUCCESS\", \"FAILURE\", \"REVOKED\", \"RETRY\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"PENDING",
									"STARTED",
									"SUCCESS",
									"FAILURE",
									"REVOKED",
									"RETRY",
								),
							},
						},
						"task_data": schema.SingleNestedAttribute{
							Description: "The object will correspond to the task type that was specified in the original request. There will be one object for transcription, another for searching for nudity, and so on.",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[StreamingAITasksTaskDataDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"task_name": schema.StringAttribute{
									Description: "Name of the task to be performed\nAvailable values: \"transcription\", \"content-moderation\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("transcription", "content-moderation"),
									},
								},
								"url": schema.StringAttribute{
									Description: "URL to the MP4 file to analyse. File must be publicly accessible via HTTP/HTTPS.",
									Computed:    true,
								},
								"audio_language": schema.StringAttribute{
									Description: "Language in original audio (transcription only). This value is used to determine the language from which to transcribe.\nIf this is not set, the system will run auto language identification and the subtitles will be in the detected language. The method also works based on AI analysis. It's fairly accurate, but if it's wrong, then set the language explicitly.\nAdditionally, when this is not set, we also support recognition of alternate languages in the video (language code-switching).\nLanguage is set by 3-letter language code according to ISO-639-2 (bibliographic code).\nWe can process languages:\n- 'afr': Afrikaans\n- 'alb': Albanian\n- 'amh': Amharic\n- 'ara': Arabic\n- 'arm': Armenian\n- 'asm': Assamese\n- 'aze': Azerbaijani\n- 'bak': Bashkir\n- 'baq': Basque\n- 'bel': Belarusian\n- 'ben': Bengali\n- 'bos': Bosnian\n- 'bre': Breton\n- 'bul': Bulgarian\n- 'bur': Myanmar\n- 'cat': Catalan\n- 'chi': Chinese\n- 'cze': Czech\n- 'dan': Danish\n- 'dut': Nynorsk\n- 'eng': English\n- 'est': Estonian\n- 'fao': Faroese\n- 'fin': Finnish\n- 'fre': French\n- 'geo': Georgian\n- 'ger': German\n- 'glg': Galician\n- 'gre': Greek\n- 'guj': Gujarati\n- 'hat': Haitian creole\n- 'hau': Hausa\n- 'haw': Hawaiian\n- 'heb': Hebrew\n- 'hin': Hindi\n- 'hrv': Croatian\n- 'hun': Hungarian\n- 'ice': Icelandic\n- 'ind': Indonesian\n- 'ita': Italian\n- 'jav': Javanese\n- 'jpn': Japanese\n- 'kan': Kannada\n- 'kaz': Kazakh\n- 'khm': Khmer\n- 'kor': Korean\n- 'lao': Lao\n- 'lat': Latin\n- 'lav': Latvian\n- 'lin': Lingala\n- 'lit': Lithuanian\n- 'ltz': Luxembourgish\n- 'mac': Macedonian\n- 'mal': Malayalam\n- 'mao': Maori\n- 'mar': Marathi\n- 'may': Malay\n- 'mlg': Malagasy\n- 'mlt': Maltese\n- 'mon': Mongolian\n- 'nep': Nepali\n- 'dut': Dutch\n- 'nor': Norwegian\n- 'oci': Occitan\n- 'pan': Punjabi\n- 'per': Persian\n- 'pol': Polish\n- 'por': Portuguese\n- 'pus': Pashto\n- 'rum': Romanian\n- 'rus': Russian\n- 'san': Sanskrit\n- 'sin': Sinhala\n- 'slo': Slovak\n- 'slv': Slovenian\n- 'sna': Shona\n- 'snd': Sindhi\n- 'som': Somali\n- 'spa': Spanish\n- 'srp': Serbian\n- 'sun': Sundanese\n- 'swa': Swahili\n- 'swe': Swedish\n- 'tam': Tamil\n- 'tat': Tatar\n- 'tel': Telugu\n- 'tgk': Tajik\n- 'tgl': Tagalog\n- 'tha': Thai\n- 'tib': Tibetan\n- 'tuk': Turkmen\n- 'tur': Turkish\n- 'ukr': Ukrainian\n- 'urd': Urdu\n- 'uzb': Uzbek\n- 'vie': Vietnamese\n- 'wel': Welsh\n- 'yid': Yiddish\n- 'yor': Yoruba",
									Computed:    true,
								},
								"client_entity_data": schema.StringAttribute{
									Description: "Meta parameter, designed to store your own extra information about a video entity: video source, video id, etc. It is not used in any way in video processing.\nFor example, if an AI-task was created automatically when you uploaded a video with the AI auto-processing option (transcribing, translationing), then the ID of the associated video for which the task was performed will be explicitly indicated here.",
									Computed:    true,
								},
								"client_user_id": schema.StringAttribute{
									Description: "Meta parameter, designed to store your own identifier. Can be used by you to tag requests from different end-users. It is not used in any way in video processing.",
									Computed:    true,
								},
								"subtitles_language": schema.StringAttribute{
									Description: "Indicates which language it is clearly necessary to translate into.\nIf this is not set, the original language will be used from attribute \"`audio_language`\".\nPlease note that:\n- transcription into the original language is a free procedure,\n- and translation from the original language into any other languages is a \"translation\" procedure and is paid. More details in [POST /ai/tasks#transcribe](https://api.gcore.com/docs/streaming/docs/api-reference/streaming/ai/create-ai-asr-task).\nLanguage is set by 3-letter language code according to ISO-639-2 (bibliographic code).",
									Computed:    true,
								},
								"category": schema.StringAttribute{
									Description: "AI content moderation with NSFW detection algorithm\nAvailable values: \"nsfw\", \"sport\", \"weapon\", \"hard_nudity\", \"soft_nudity\", \"child_pornography\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive(
											"nsfw",
											"sport",
											"weapon",
											"hard_nudity",
											"soft_nudity",
											"child_pornography",
										),
									},
								},
								"stop_objects": schema.StringAttribute{
									Description: "Comma separated objects, and probabilities, that will cause the processing to stop immediatelly after finding.\nAvailable values: \"ANUS_EXPOSED\", \"BUTTOCKS_EXPOSED\", \"FEMALE_BREAST_EXPOSED\", \"FEMALE_GENITALIA_EXPOSED\", \"MALE_BREAST_EXPOSED\", \"MALE_GENITALIA_EXPOSED\", \"ANUS_COVERED\", \"ARMPITS_COVERED\", \"ARMPITS_EXPOSED\", \"BELLY_COVERED\", \"BELLY_EXPOSED\", \"BUTTOCKS_COVERED\", \"FACE_FEMALE\", \"FACE_MALE\", \"FEET_COVERED\", \"FEET_EXPOSED\", \"FEMALE_BREAST_COVERED\", \"FEMALE_GENITALIA_COVERED\".",
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive(
											"ANUS_EXPOSED",
											"BUTTOCKS_EXPOSED",
											"FEMALE_BREAST_EXPOSED",
											"FEMALE_GENITALIA_EXPOSED",
											"MALE_BREAST_EXPOSED",
											"MALE_GENITALIA_EXPOSED",
											"ANUS_COVERED",
											"ARMPITS_COVERED",
											"ARMPITS_EXPOSED",
											"BELLY_COVERED",
											"BELLY_EXPOSED",
											"BUTTOCKS_COVERED",
											"FACE_FEMALE",
											"FACE_MALE",
											"FEET_COVERED",
											"FEET_EXPOSED",
											"FEMALE_BREAST_COVERED",
											"FEMALE_GENITALIA_COVERED",
										),
									},
								},
							},
						},
						"task_id": schema.StringAttribute{
							Description: "ID of the AI task",
							Computed:    true,
						},
						"task_name": schema.StringAttribute{
							Description: "Type of AI task\nAvailable values: \"content-moderation\", \"transcription\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("content-moderation", "transcription"),
							},
						},
					},
				},
			},
		},
	}
}

func (d *StreamingAITasksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *StreamingAITasksDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
