package runtime

import (
	"context"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/kiosk404/airi-go/backend/api/crossdomain/agentrun"
	"github.com/kiosk404/airi-go/backend/api/crossdomain/singleagent"
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	"github.com/kiosk404/airi-go/backend/infra/contract/imagex"
	agentEntity "github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/domain/repo"
	"github.com/kiosk404/airi-go/backend/modules/conversation/agent_run/pkg"
	"github.com/kiosk404/airi-go/backend/modules/conversation/conversation/pkg/errno"
	crossmessage "github.com/kiosk404/airi-go/backend/modules/conversation/message/crossdomain"
	msgEntity "github.com/kiosk404/airi-go/backend/modules/conversation/message/domain/entity"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
)

type AgentRuntime struct {
	RunRecord     *agentEntity.RunRecordMeta
	AgentInfo     *singleagent.SingleAgent
	QuestionMsgID int64
	RunMeta       *agentEntity.AgentRunMeta
	StartTime     time.Time
	Input         *msgEntity.Message
	HistoryMsg    []*msgEntity.Message
	Usage         *agentrun.Usage
	SW            *schema.StreamWriter[*agentEntity.AgentRunResponse]

	RunProcess    *RunProcess
	RunRecordRepo repo.RunRecordRepo
	ImagexClient  imagex.ImageX
	MessageEvent  *Event
}

func (rd *AgentRuntime) SetRunRecord(runRecord *agentEntity.RunRecordMeta) {
	rd.RunRecord = runRecord
}

func (rd *AgentRuntime) GetRunRecord() *agentEntity.RunRecordMeta {
	return rd.RunRecord
}

func (rd *AgentRuntime) SetUsage(usage *agentrun.Usage) {
	rd.Usage = usage
}
func (rd *AgentRuntime) GetUsage() *agentrun.Usage {
	return rd.Usage
}

func (rd *AgentRuntime) SetRunMeta(arm *agentEntity.AgentRunMeta) {
	rd.RunMeta = arm
}
func (rd *AgentRuntime) GetRunMeta() *agentEntity.AgentRunMeta {
	return rd.RunMeta
}
func (rd *AgentRuntime) SetAgentInfo(agentInfo *singleagent.SingleAgent) {
	rd.AgentInfo = agentInfo
}
func (rd *AgentRuntime) GetAgentInfo() *singleagent.SingleAgent {
	return rd.AgentInfo
}
func (rd *AgentRuntime) SetQuestionMsgID(msgID int64) {
	rd.QuestionMsgID = msgID
}
func (rd *AgentRuntime) GetQuestionMsgID() int64 {
	return rd.QuestionMsgID
}
func (rd *AgentRuntime) SetStartTime(t time.Time) {
	rd.StartTime = t
}
func (rd *AgentRuntime) GetStartTime() time.Time {
	return rd.StartTime
}
func (rd *AgentRuntime) SetInput(input *msgEntity.Message) {
	rd.Input = input
}
func (rd *AgentRuntime) GetInput() *msgEntity.Message {
	return rd.Input
}

func (rd *AgentRuntime) SetHistoryMsg(histroyMsg []*msgEntity.Message) {
	rd.HistoryMsg = histroyMsg
}

func (rd *AgentRuntime) GetHistory() []*msgEntity.Message {
	return rd.HistoryMsg
}

func (art *AgentRuntime) Run(ctx context.Context) (err error) {

	agentInfo, err := getAgentInfo(ctx, art.GetRunMeta().AgentID, art.GetRunMeta().IsDraft)
	if err != nil {
		return
	}

	art.SetAgentInfo(agentInfo)

	history, err := art.getHistory(ctx)
	if err != nil {
		return
	}

	runRecord, err := art.createRunRecord(ctx)

	if err != nil {
		return
	}

	art.SetRunRecord(runRecord)
	art.SetHistoryMsg(history)

	defer func() {
		srRecord := buildSendRunRecord(ctx, runRecord, agentEntity.RunStatusCompleted)
		if err != nil {
			srRecord.Error = &agentEntity.RunError{
				Code: errno.ErrConversationAgentRunError,
				Msg:  err.Error(),
			}
			art.RunProcess.StepToFailed(ctx, srRecord, art.SW)
			return
		}
		art.RunProcess.StepToComplete(ctx, srRecord, art.SW, art.GetUsage())
	}()
	mh := &MesssageEventHanlder{
		messageEvent: art.MessageEvent,
		sw:           art.SW,
	}
	input, err := mh.HandlerInput(ctx, art)
	if err != nil {
		return
	}
	art.SetInput(input)

	art.SetQuestionMsgID(input.ID)

	if art.GetAgentInfo().BotMode == bot_common.BotMode_WorkflowMode {
		err = art.ChatflowRun(ctx, art.ImagexClient)
	} else {
		err = art.AgentStreamExecute(ctx, art.ImagexClient)
	}
	return
}

func (art *AgentRuntime) getHistory(ctx context.Context) ([]*msgEntity.Message, error) {

	conversationTurns := getAgentHistoryRounds(art.GetAgentInfo())

	runRecordList, err := art.RunRecordRepo.List(ctx, &agentEntity.ListRunRecordMeta{
		ConversationID: art.GetRunMeta().ConversationID,
		SectionID:      art.GetRunMeta().SectionID,
		Limit:          conversationTurns,
	})
	if err != nil {
		return nil, err
	}

	if len(runRecordList) == 0 {
		return nil, nil
	}
	runIDS := concactRunID(runRecordList)
	history, err := crossmessage.DefaultSVC().GetByRunIDs(ctx, art.GetRunMeta().ConversationID, runIDS)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func concactRunID(rr []*agentEntity.RunRecordMeta) []int64 {
	ids := make([]int64, 0, len(rr))
	for _, c := range rr {
		ids = append(ids, c.ID)
	}

	return ids
}

func (art *AgentRuntime) createRunRecord(ctx context.Context) (*agentEntity.RunRecordMeta, error) {
	runPoData, err := art.RunRecordRepo.Create(ctx, art.GetRunMeta())
	if err != nil {
		logs.Error(pkg.ModelName, "RunRecordRepo.Create error: %v", err)
		return nil, err
	}

	srRecord := buildSendRunRecord(ctx, runPoData, agentEntity.RunStatusCreated)

	art.RunProcess.StepToCreate(ctx, srRecord, art.SW)

	err = art.RunProcess.StepToInProgress(ctx, srRecord, art.SW)
	if err != nil {
		logs.ErrorX(pkg.ModelName, "runProcess.StepToInProgress error: %v", err)
		return nil, err
	}
	return runPoData, nil
}
