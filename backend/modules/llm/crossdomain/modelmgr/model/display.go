package model

type DisplayInfo struct {
	Name         string    `json:"name" query:"name"`
	Description  *I18nText `json:"description" query:"description"`
	OutputTokens int64     `json:"output_tokens" query:"output_tokens"`
	MaxTokens    int64     `json:"max_tokens" query:"max_tokens"`
}

func NewDisplayInfo() *DisplayInfo {
	return &DisplayInfo{}
}

func (p *DisplayInfo) InitDefault() {
}

func (p *DisplayInfo) GetName() (v string) {
	return p.Name
}

var DisplayInfo_Description_DEFAULT *I18nText

func (p *DisplayInfo) GetDescription() (v *I18nText) {
	if !p.IsSetDescription() {
		return DisplayInfo_Description_DEFAULT
	}
	return p.Description
}

func (p *DisplayInfo) GetOutputTokens() (v int64) {
	return p.OutputTokens
}

func (p *DisplayInfo) GetMaxTokens() (v int64) {
	return p.MaxTokens
}

func (p *DisplayInfo) IsSetDescription() bool {
	return p.Description != nil
}
