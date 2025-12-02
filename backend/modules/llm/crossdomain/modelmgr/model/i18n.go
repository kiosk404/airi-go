package model

type I18nText struct {
	ZhCn string `json:"zh_cn" query:"zh_cn"`
	EnUs string `json:"en_us" query:"en_us"`
}

func NewI18nText() *I18nText {
	return &I18nText{}
}

func (p *I18nText) InitDefault() {
}

func (p *I18nText) GetZhCn() (v string) {
	return p.ZhCn
}

func (p *I18nText) GetEnUs() (v string) {
	return p.EnUs
}

var fieldIDToName_I18nText = map[int16]string{
	1: "zh_cn",
	2: "en_us",
}
