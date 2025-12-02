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

func (p *ModelAbility) GetCotDisplay() (v bool) {
	v = p.CotDisplay
	return
}

func (p *ModelAbility) GetFunctionCall() (v bool) {
	v = p.FunctionCall
	return
}

func (p *ModelAbility) GetImageUnderstanding() (v bool) {
	v = p.ImageUnderstanding
	return
}

func (p *ModelAbility) GetVideoUnderstanding() (v bool) {
	v = p.VideoUnderstanding
	return
}

func (p *ModelAbility) GetAudioUnderstanding() (v bool) {
	v = p.AudioUnderstanding
	return
}

func (p *ModelAbility) GetSupportMultiModal() (v bool) {
	v = p.SupportMultiModal
	return
}

func (p *ModelAbility) GetPrefillResp() (v bool) {
	v = p.PrefillResp
	return
}
