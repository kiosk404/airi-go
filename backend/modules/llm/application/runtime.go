package application

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/cloudwego/eino/components/model"
	druntime "github.com/kiosk404/airi-go/backend/api/model/llm/domain/runtime"
	"github.com/kiosk404/airi-go/backend/api/model/llm/runtime"
	"github.com/kiosk404/airi-go/backend/infra/contract/limiter"
	"github.com/kiosk404/airi-go/backend/modules/llm/application/convertor"
	llmmodel "github.com/kiosk404/airi-go/backend/modules/llm/crossdomain/modelmgr/model"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/llm/domain/service"
	"github.com/kiosk404/airi-go/backend/modules/llm/pkg"
	llmerrorx "github.com/kiosk404/airi-go/backend/modules/llm/pkg/errno"
	"github.com/kiosk404/airi-go/backend/modules/llm/pkg/traceutil"
	"github.com/kiosk404/airi-go/backend/pkg/errorx"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"github.com/kiosk404/airi-go/backend/pkg/utils/goroutineutil"
	"github.com/pkg/errors"
)

type runtimeApp struct {
	manageSrv   service.IManage
	runtimeSrv  service.IRuntime
	rateLimiter limiter.IRateLimiter
}

func NewRuntimeApplication(
	manageSrv service.IManage,
	runtimeSrv service.IRuntime,
	factory limiter.IRateLimiterFactory,
) runtime.LLMRuntimeService {
	return &runtimeApp{
		manageSrv:   manageSrv,
		runtimeSrv:  runtimeSrv,
		rateLimiter: factory.NewRateLimiter(),
	}
}

func (r *runtimeApp) Chat(ctx context.Context, req *runtime.ChatRequest) (resp *runtime.ChatResponse, err error) {
	resp = runtime.NewChatResponse()
	if err = r.validateChatReq(ctx, req); err != nil {
		return resp, errorx.NewByCode(llmerrorx.RequestNotValidCode, errorx.WithExtraMsg(err.Error()))
	}
	// 1. 模型信息获取
	model, err := r.manageSrv.GetModelByID(ctx, req.GetModelConfig().GetModelID())
	if err != nil {
		return resp, err
	}
	// 2. model参数校验
	if err = model.Valid(); err != nil {
		return resp, errorx.NewByCode(llmerrorx.ModelInvalidCode, errorx.WithExtraMsg(err.Error()))
	}
	// 3. 限流
	if err = r.rateLimitAllow(ctx, req, model); err != nil {
		return resp, err
	}
	// 4. 格式转换
	msgs := convertor.MessagesDTO2DO(req.GetMessages())
	msgs, err = r.runtimeSrv.HandleMsgsPreCallModel(ctx, model, msgs)
	if err != nil {
		return resp, errorx.NewByCode(llmerrorx.RequestNotValidCode, errorx.WithExtraMsg(err.Error()))
	}
	options := convertor.ModelAndTools2OptionDOs(req.GetModelConfig(), req.GetTools())
	var respMsg *entity.Message

	// 5. 调用llm.generate or llm.stream方法, 并解析流式返回
	defer func() {
		// 异步记录本次模型请求
		r.recordModelRequest(ctx, &recordModelRequestParam{
			bizParam: req.BizParam,
			model:    model,
			input:    msgs,
			lastMsg:  respMsg,
			err:      err,
		})
	}()
	respMsg, err = r.runtimeSrv.Generate(ctx, model, msgs, options...)
	if err != nil {
		return resp, err
	}
	msgDTO := convertor.MessageDO2DTO(respMsg)
	resp.SetMessage(msgDTO)
	return resp, nil
}

func (r *runtimeApp) ChatStream(req *runtime.ChatRequest, stream runtime.LLMRuntimeService_ChatStreamServer) (err error) {
	ctx := context.Background()
	// 参数校验
	if err = r.validateChatReq(ctx, req); err != nil {
		return errorx.NewByCode(llmerrorx.RequestNotValidCode, errorx.WithExtraMsg(err.Error()))
	}
	// 1. 模型信息获取
	model, err := r.manageSrv.GetModelByID(ctx, req.GetModelConfig().GetModelID())
	if err != nil {
		return err
	}
	// 对model参数做校验
	if err = model.Valid(); err != nil {
		return errorx.NewByCode(llmerrorx.ModelInvalidCode, errorx.WithExtraMsg(err.Error()))
	}
	// 2. 限流
	if err = r.rateLimitAllow(ctx, req, model); err != nil {
		return err
	}
	// 3. 格式转换
	msgs := convertor.MessagesDTO2DO(req.GetMessages())
	msgs, err = r.runtimeSrv.HandleMsgsPreCallModel(ctx, model, msgs)
	if err != nil {
		return errorx.NewByCode(llmerrorx.RequestNotValidCode, errorx.WithExtraMsg(err.Error()))
	}
	options := convertor.ModelAndTools2OptionDOs(req.GetModelConfig(), req.GetTools())
	// 4. 调用llm.generate or llm.stream方法, 并解析流式返回
	var parseResult entity.StreamRespParseResult
	beginTime := time.Now()
	defer func() {
		// 异步记录本次模型请求
		r.recordModelRequest(ctx, &recordModelRequestParam{
			bizParam: req.BizParam,
			model:    model,
			input:    msgs,
			lastMsg:  parseResult.LastRespMsg,
			err:      err,
		})
	}()
	sr, err := r.runtimeSrv.Stream(ctx, model, msgs, options...)
	if err != nil {
		return err
	}
	if parseResult, err = r.parseChatStreamResp(sr, stream, beginTime); err != nil {
		return errorx.NewByCode(llmerrorx.ParseModelRespFailedCode, errorx.WithExtraMsg(err.Error()))
	}
	return nil
}

func (r *runtimeApp) parseChatStreamResp(streamDO entity.IStreamReader, streamDTO runtime.LLMRuntimeService_ChatStreamServer,
	beginTime time.Time,
) (parseResult entity.StreamRespParseResult, err error) {
	var hasReasoningContent bool
	for {
		msgDO, err := streamDO.Recv()
		if int64(parseResult.FirstTokenLatency) <= int64(0) {
			parseResult.FirstTokenLatency = time.Since(beginTime)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return parseResult, err
		}
		if msgDO.ReasoningContent != "" {
			hasReasoningContent = true
		}
		// 计算reasoning duration
		if hasReasoningContent && msgDO.ReasoningContent == "" && parseResult.ReasoningDuration == 0 {
			parseResult.ReasoningDuration = time.Since(beginTime)
		}
		parseResult.LastRespMsg = msgDO
		parseResult.RespMsgs = append(parseResult.RespMsgs, msgDO)
		msgDTO := convertor.MessageDO2DTO(msgDO)
		respDTO := &runtime.ChatResponse{
			Message: msgDTO,
		}
		if err := streamDTO.Send(respDTO); err != nil {
			return parseResult, err
		}
	}
	return parseResult, nil
}

func (r *runtimeApp) rateLimitAllow(ctx context.Context, req *runtime.ChatRequest, model *entity.Model) error {
	var scenario *entity.Scenario
	if req.GetBizParam() != nil && req.GetBizParam().Scenario != nil {
		scenario = convertor.ScenarioPtrDTO2DTO(req.GetBizParam().Scenario)
	} else {
		scenario = ptr.Of(llmmodel.ScenarioDefault)
	}
	// 获得模型在此场景下的qpm tpm
	sceneCfg := model.GetScenarioConfig(scenario)
	if sceneCfg == nil || sceneCfg.Quota == nil {
		return nil
	}
	qpm := ptr.From(sceneCfg.Quota.Qpm)
	tpm := ptr.From(sceneCfg.Quota.Tpm)
	// qpm
	if qpm >= 0 {
		qpmKey := fmt.Sprintf("%s:%d:%s", "qpm", model.ID, *scenario)
		result, err := r.rateLimiter.AllowN(ctx, qpmKey, 1, limiter.WithLimit(&limiter.Limit{
			Rate:   int(qpm),
			Burst:  int(qpm),
			Period: time.Minute,
		}))
		if err == nil && result != nil && !result.Allowed {
			return errorx.NewByCode(llmerrorx.ModelQPMLimitCode)
		}
	}
	// tpm
	if tpm >= 0 {
		tpmKey := fmt.Sprintf("%s:%d:%s", "tpm", model.ID, *scenario)
		result, err := r.rateLimiter.AllowN(ctx, tpmKey, int(req.GetModelConfig().GetMaxTokens()), limiter.WithLimit(&limiter.Limit{
			Rate:   int(tpm),
			Burst:  int(tpm),
			Period: time.Minute,
		}))
		if err == nil && result != nil && !result.Allowed {
			return errorx.NewByCode(llmerrorx.ModelTPMLimitCode)
		}
	}
	return nil
}

type recordModelRequestParam struct {
	bizParam *druntime.BizParam
	model    *entity.Model
	input    []*entity.Message
	lastMsg  *entity.Message
	err      error
}

func (r *runtimeApp) recordModelRequest(ctx context.Context, param *recordModelRequestParam) {
	goroutineutil.GoWithDefaultRecovery(ctx, func() {
		record := &entity.ModelRequestRecord{
			UserID:              param.bizParam.GetUserID(),
			UsageScene:          llmmodel.Scenario(param.bizParam.GetScenario()),
			UsageSceneEntityID:  param.bizParam.GetScenarioEntityID(),
			Protocol:            ptr.FromPtrConvert[llmmodel.Protocol, entity.Protocol](param.model.Protocol),
			ModelIdentification: ptr.From(param.model.ProtocolConfig.Model),
			ModelAk:             ptr.From(param.model.ProtocolConfig.APIKey),
			ModelID:             strconv.FormatInt(param.model.ID, 10),
			ModelName:           param.model.Name,
			InputToken:          int64(param.lastMsg.GetInputToken()),
			OutputToken:         int64(param.lastMsg.GetOutputToken()),
			LogId:               logs.GetLogID(ctx),
		}
		if param.err != nil {
			record.ErrorCode = strconv.FormatInt(int64(traceutil.GetTraceStatusCode(param.err)), 10)
			record.ErrorMsg = ptr.Of(param.err.Error())
		}
		if err := r.runtimeSrv.CreateModelRequestRecord(ctx, record); err != nil {
			logs.WarnX(pkg.ModelName, "[recordModelRequest] failed, err:%v", err)
		}
	})
}

type setSpanParam struct {
	stream     bool
	inputMsgs  []*entity.Message
	toolInfos  []*entity.ToolInfo
	toolChoice *entity.ToolChoice
	options    []model.Option
	model      *entity.Model
	bizParam   *druntime.BizParam

	firstTokenLatency time.Duration
	reasoningDuration time.Duration
	err               error
	respMsgs          []*entity.Message
}

func (r *runtimeApp) validateChatReq(ctx context.Context, req *runtime.ChatRequest) (err error) {
	if req.GetModelConfig() == nil {
		return errors.Errorf("model config is required")
	}
	if len(req.GetMessages()) == 0 {
		return errors.Errorf("messages is required")
	}
	if req.GetBizParam() == nil {
		return errors.Errorf("bizParam is required")
	}
	if !req.GetBizParam().IsSetScenario() {
		return errors.Errorf("bizParam.scenario is required")
	}
	if !req.GetBizParam().IsSetScenarioEntityID() {
		return errors.Errorf("bizParam.scenario_entity_id is required")
	}
	// if !req.GetBizParam().IsSetUserID() {
	// 	return errors.Errorf("bizParam.user_id is required")
	// }
	return nil
}
