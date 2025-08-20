// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_ai_task

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*StreamingAITaskResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"task_id": schema.StringAttribute{
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"task_name": schema.StringAttribute{
				Description: "Name of the task to be performed\nAvailable values: \"transcription\", \"content-moderation\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("transcription", "content-moderation"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"url": schema.StringAttribute{
				Description:   "URL to the MP4 file to analyse. File must be publicly accessible via HTTP/HTTPS.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"audio_language": schema.StringAttribute{
				Description:   "Language in original audio (transcription only). This value is used to determine the language from which to transcribe.\nIf this is not set, the system will run auto language identification and the subtitles will be in the detected language. The method also works based on AI analysis. It's fairly accurate, but if it's wrong, then set the language explicitly.\nAdditionally, when this is not set, we also support recognition of alternate languages in the video (language code-switching).\nLanguage is set by 3-letter language code according to ISO-639-2 (bibliographic code).\nWe can process languages:\n- 'afr': Afrikaans\n- 'alb': Albanian\n- 'amh': Amharic\n- 'ara': Arabic\n- 'arm': Armenian\n- 'asm': Assamese\n- 'aze': Azerbaijani\n- 'bak': Bashkir\n- 'baq': Basque\n- 'bel': Belarusian\n- 'ben': Bengali\n- 'bos': Bosnian\n- 'bre': Breton\n- 'bul': Bulgarian\n- 'bur': Myanmar\n- 'cat': Catalan\n- 'chi': Chinese\n- 'cze': Czech\n- 'dan': Danish\n- 'dut': Nynorsk\n- 'eng': English\n- 'est': Estonian\n- 'fao': Faroese\n- 'fin': Finnish\n- 'fre': French\n- 'geo': Georgian\n- 'ger': German\n- 'glg': Galician\n- 'gre': Greek\n- 'guj': Gujarati\n- 'hat': Haitian creole\n- 'hau': Hausa\n- 'haw': Hawaiian\n- 'heb': Hebrew\n- 'hin': Hindi\n- 'hrv': Croatian\n- 'hun': Hungarian\n- 'ice': Icelandic\n- 'ind': Indonesian\n- 'ita': Italian\n- 'jav': Javanese\n- 'jpn': Japanese\n- 'kan': Kannada\n- 'kaz': Kazakh\n- 'khm': Khmer\n- 'kor': Korean\n- 'lao': Lao\n- 'lat': Latin\n- 'lav': Latvian\n- 'lin': Lingala\n- 'lit': Lithuanian\n- 'ltz': Luxembourgish\n- 'mac': Macedonian\n- 'mal': Malayalam\n- 'mao': Maori\n- 'mar': Marathi\n- 'may': Malay\n- 'mlg': Malagasy\n- 'mlt': Maltese\n- 'mon': Mongolian\n- 'nep': Nepali\n- 'dut': Dutch\n- 'nor': Norwegian\n- 'oci': Occitan\n- 'pan': Punjabi\n- 'per': Persian\n- 'pol': Polish\n- 'por': Portuguese\n- 'pus': Pashto\n- 'rum': Romanian\n- 'rus': Russian\n- 'san': Sanskrit\n- 'sin': Sinhala\n- 'slo': Slovak\n- 'slv': Slovenian\n- 'sna': Shona\n- 'snd': Sindhi\n- 'som': Somali\n- 'spa': Spanish\n- 'srp': Serbian\n- 'sun': Sundanese\n- 'swa': Swahili\n- 'swe': Swedish\n- 'tam': Tamil\n- 'tat': Tatar\n- 'tel': Telugu\n- 'tgk': Tajik\n- 'tgl': Tagalog\n- 'tha': Thai\n- 'tib': Tibetan\n- 'tuk': Turkmen\n- 'tur': Turkish\n- 'ukr': Ukrainian\n- 'urd': Urdu\n- 'uzb': Uzbek\n- 'vie': Vietnamese\n- 'wel': Welsh\n- 'yid': Yiddish\n- 'yor': Yoruba",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"category": schema.StringAttribute{
				Description: "Model for analysis (content-moderation only). Determines what exactly needs to be found in the video.\nAvailable values: \"sport\", \"weapon\", \"nsfw\", \"hard_nudity\", \"soft_nudity\", \"child_pornography\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"sport",
						"weapon",
						"nsfw",
						"hard_nudity",
						"soft_nudity",
						"child_pornography",
					),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"client_entity_data": schema.StringAttribute{
				Description:   "Meta parameter, designed to store your own extra information about a video entity: video source, video id, etc. It is not used in any way in video processing.\nFor example, if an AI-task was created automatically when you uploaded a video with the AI auto-processing option (nudity detection, etc), then the ID of the associated video for which the task was performed will be explicitly indicated here.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"client_user_id": schema.StringAttribute{
				Description:   "Meta parameter, designed to store your own identifier. Can be used by you to tag requests from different end-users. It is not used in any way in video processing.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"subtitles_language": schema.StringAttribute{
				Description:   "Indicates which language it is clearly necessary to translate into.\nIf this is not set, the original language will be used from attribute \"`audio_language`\".\nPlease note that:\n- transcription into the original language is a free procedure,\n- and translation from the original language into any other languages is a \"translation\" procedure and is paid. More details in [POST /ai/tasks#transcribe](https://api.gcore.com/docs/streaming/docs/api-reference/streaming/ai/create-ai-asr-task).\nLanguage is set by 3-letter language code according to ISO-639-2 (bibliographic code).",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
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
			"processing_time": schema.SingleNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectType[StreamingAITaskProcessingTimeModel](ctx),
				Attributes: map[string]schema.Attribute{
					"completed_at": schema.StringAttribute{
						Description: "Video processing end time. Format is date time in ISO 8601",
						Computed:    true,
					},
					"started_at": schema.StringAttribute{
						Description: "Video processing start time. Format is date time in ISO 8601",
						Computed:    true,
					},
					"total_time_sec": schema.Float64Attribute{
						Description: "Duration of video processing in seconds",
						Computed:    true,
					},
				},
			},
			"result": schema.SingleNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectType[StreamingAITaskResultModel](ctx),
				Attributes: map[string]schema.Attribute{
					"concatenated_text": schema.StringAttribute{
						Description: "Full text of the analyzed video. The value is unstructured, unformatted text.",
						Computed:    true,
					},
					"languages": schema.ListAttribute{
						Description: "An array of language codes that were discovered and/or used in transcription. If the audio or subtitle language was explicitly specified in the initial parameters, it will be copied here. For automatic detection the identified languages will be displayed here. Also please note that for multilingual audio, the first 5 languages are displayed in order of frequency of use.",
						Computed:    true,
						CustomType:  customfield.NewListType[types.String](ctx),
						ElementType: types.StringType,
					},
					"speech_detected": schema.BoolAttribute{
						Description: "Determines whether speech was detected or not.\nPlease note: If the task is in \"SUCCESS\" status and speech was not found in the entire file, then \"false\" will be indicated here and the ```subtitles``` field will be empty.",
						Computed:    true,
					},
					"subtitles": schema.ListNestedAttribute{
						Description: `An array of phrases divided into time intervals, in the format "json". Suitable when you need to display the result in chronometric form, or transfer the text for further processing.`,
						Computed:    true,
						CustomType:  customfield.NewNestedObjectListType[StreamingAITaskResultSubtitlesModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"end_time": schema.StringAttribute{
									Description: `End time of the phrase, when it ends in the video. Format is "HH:mm:ss.fff".`,
									Computed:    true,
								},
								"start_time": schema.StringAttribute{
									Description: `Start time of the phrase when it is heard in the video. Format is "HH:mm:ss.fff".`,
									Computed:    true,
								},
								"text": schema.StringAttribute{
									Description: "A complete phrase that sounds during a specified period of time.",
									Computed:    true,
								},
							},
						},
					},
					"vtt_content": schema.StringAttribute{
						Description: "Auto generated subtitles in WebVTT format.",
						Computed:    true,
					},
					"detection_results": schema.ListAttribute{
						Computed: true,
						Validators: []validator.List{
							listvalidator.ValueStringsAre(
								stringvalidator.OneOfCaseInsensitive(
									"archery",
									"arm wrestling",
									"playing badminton",
									"playing baseball",
									"basketball dunk",
									"bowling",
									"boxing punch",
									"boxing speed bag",
									"catching or throwing baseball",
									"catching or throwing softball",
									"cricket",
									"curling",
									"disc golfing",
									"dodgeball",
									"fencing",
									"football",
									"golf chipping",
									"golf driving",
									"golf putting",
									"hitting baseball",
									"hockey stop",
									"ice skating",
									"javelin throw",
									"juggling soccer ball",
									"kayaking",
									"kicking field goal",
									"kicking soccer ball",
									"playing cricket",
									"playing field hockey",
									"playing ice hockey",
									"playing kickball",
									"playing lacrosse",
									"playing ping pong",
									"playing polo",
									"playing squash or racquetball",
									"playing tennis",
									"playing volleyball",
									"pole vault",
									"riding a bike",
									"riding or walking with horse",
									"roller skating",
									"rowing",
									"sailing",
									"shooting goal (soccer)",
									"skateboarding",
									"skiing",
									"gun",
									"heavy weapon",
									"knife",
									"nsfw",
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
									"0-2",
									"3-9",
									"10-19",
								),
							),
						},
						CustomType:  customfield.NewListType[types.String](ctx),
						ElementType: types.StringType,
					},
					"frames": schema.ListNestedAttribute{
						Computed:   true,
						CustomType: customfield.NewNestedObjectListType[StreamingAITaskResultFramesModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"confidence": schema.Float64Attribute{
									Description: "Percentage of probability of identifying the activity",
									Computed:    true,
								},
								"frame_number": schema.Int64Attribute{
									Description: "Video frame number where activity was found",
									Computed:    true,
								},
								"label": schema.StringAttribute{
									Description: "Type of detected activity",
									Computed:    true,
								},
							},
						},
					},
					"sport_detected": schema.BoolAttribute{
						Description: "A boolean value whether any sports were detected",
						Computed:    true,
					},
					"weapon_detected": schema.BoolAttribute{
						Description: "A boolean value whether any weapon was detected",
						Computed:    true,
					},
					"nsfw_detected": schema.BoolAttribute{
						Description: "A boolean value whether any Not Safe For Work content was detected",
						Computed:    true,
					},
					"porn_detected": schema.BoolAttribute{
						Description: "A boolean value whether any nudity was detected",
						Computed:    true,
					},
					"child_pornography_detected": schema.BoolAttribute{
						Description: "A boolean value whether child pornography was detected",
						Computed:    true,
					},
					"error": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"task_data": schema.SingleNestedAttribute{
				Description: "The object will correspond to the task type that was specified in the original request. There will be one object for transcription, another for searching for nudity, and so on.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[StreamingAITaskTaskDataModel](ctx),
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
		},
	}
}

func (r *StreamingAITaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *StreamingAITaskResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
