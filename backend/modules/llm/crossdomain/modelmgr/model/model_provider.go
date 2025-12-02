package model

type ModelProvider struct {
	Name        *I18nText  `thrift:"name,1" form:"name" json:"name" query:"name"`
	IconURI     string     `thrift:"icon_uri,2" form:"icon_uri" json:"icon_uri" query:"icon_uri"`
	IconURL     string     `thrift:"icon_url,3" form:"icon_url" json:"icon_url" query:"icon_url"`
	Description *I18nText  `thrift:"description,4" form:"description" json:"description" query:"description"`
	ModelClass  ModelClass `thrift:"model_class,5" form:"model_class" json:"model_class" query:"model_class"`
}
