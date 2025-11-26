package model

type ModelAbility struct {
	CotDisplay         bool `json:"cot_display,omitempty"`
	FunctionCall       bool `json:"function_call,omitempty"`
	ImageUnderstanding bool `json:"image_understanding,omitempty"`
	VideoUnderstanding bool `json:"video_understanding,omitempty"`
	AudioUnderstanding bool `json:"audio_understanding,omitempty"`
	SupportMultiModal  bool `json:"support_multi_modal,omitempty"`
	PrefillResp        bool `json:"prefill_resp,omitempty"`
}
