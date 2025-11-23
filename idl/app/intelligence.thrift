include "../base.thrift"
include "search.thrift"
include  "common_struct/intelligence_common_struct.thrift"
include  "common_struct/common_struct.thrift"

namespace go app.intelligence

service IntelligenceService {
    search.GetDraftIntelligenceListResponse GetDraftIntelligenceList(1: search.GetDraftIntelligenceListRequest req) (api.post='/api/intelligence_api/search/get_draft_intelligence_list', api.category="search",agw.preserve_base="true")
    search.GetDraftIntelligenceInfoResponse GetDraftIntelligenceInfo(1: search.GetDraftIntelligenceInfoRequest req) (api.post='/api/intelligence_api/search/get_draft_intelligence_info', api.category="search",agw.preserve_base="true")
    search.GetUserRecentlyEditIntelligenceResponse GetUserRecentlyEditIntelligence(1: search.GetUserRecentlyEditIntelligenceRequest req) (api.post='/api/intelligence_api/search/get_recently_edit_intelligence', api.category="search",agw.preserve_base="true")
}