package convert

import (
	"github.com/kiosk404/airi-go/backend/api/model/component/plugin_develop/common"
	"github.com/kiosk404/airi-go/backend/modules/component/crossdomain/plugin/consts"
)

var assistTypeToFormat = map[consts.APIFileAssistType]string{
	consts.AssistTypeFile:  "file_url",
	consts.AssistTypeImage: "image_url",
	consts.AssistTypeDoc:   "doc_url",
	consts.AssistTypePPT:   "ppt_url",
	consts.AssistTypeCode:  "code_url",
	consts.AssistTypeExcel: "excel_url",
	consts.AssistTypeZIP:   "zip_url",
	consts.AssistTypeVideo: "video_url",
	consts.AssistTypeAudio: "audio_url",
	consts.AssistTypeTXT:   "txt_url",
}

func AssistTypeToFormat(typ consts.APIFileAssistType) (string, bool) {
	format, ok := assistTypeToFormat[typ]
	return format, ok
}

var formatToAssistType = func() map[string]consts.APIFileAssistType {
	types := make(map[string]consts.APIFileAssistType, len(assistTypeToFormat))
	for k, v := range assistTypeToFormat {
		types[v] = k
	}
	return types
}()

func FormatToAssistType(format string) (consts.APIFileAssistType, bool) {
	typ, ok := formatToAssistType[format]
	return typ, ok
}

var assistTypeToThriftFormat = map[consts.APIFileAssistType]common.PluginParamTypeFormat{
	consts.AssistTypeFile:  common.PluginParamTypeFormat_FileUrl,
	consts.AssistTypeImage: common.PluginParamTypeFormat_ImageUrl,
	consts.AssistTypeDoc:   common.PluginParamTypeFormat_DocUrl,
	consts.AssistTypePPT:   common.PluginParamTypeFormat_PptUrl,
	consts.AssistTypeCode:  common.PluginParamTypeFormat_CodeUrl,
	consts.AssistTypeExcel: common.PluginParamTypeFormat_ExcelUrl,
	consts.AssistTypeZIP:   common.PluginParamTypeFormat_ZipUrl,
	consts.AssistTypeVideo: common.PluginParamTypeFormat_VideoUrl,
	consts.AssistTypeAudio: common.PluginParamTypeFormat_AudioUrl,
	consts.AssistTypeTXT:   common.PluginParamTypeFormat_TxtUrl,
}

func AssistTypeToThriftFormat(typ consts.APIFileAssistType) (common.PluginParamTypeFormat, bool) {
	format, ok := assistTypeToThriftFormat[typ]
	return format, ok
}
