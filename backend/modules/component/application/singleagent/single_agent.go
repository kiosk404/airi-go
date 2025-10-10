package singleagent

import (
	"github.com/bytedance/sonic"
	"github.com/kiosk404/airi-go/backend/api/model/app/bot_common"
	"github.com/kiosk404/airi-go/backend/api/model/playground"
	singleagent "github.com/kiosk404/airi-go/backend/modules/component/agent/domain/service"
	"github.com/kiosk404/airi-go/backend/pkg/lang/ptr"
)

type SingleAgentApplicationService struct {
	appContext *ServiceComponents
	DomainSVC  singleagent.SingleAgent
}

func newApplicationService(s *ServiceComponents, domain singleagent.SingleAgent) *SingleAgentApplicationService {
	return &SingleAgentApplicationService{
		appContext: s,
		DomainSVC:  domain,
	}
}

const onboardingInfoMaxLength = 65535

func (s *SingleAgentApplicationService) generateOnboardingStr(onboardingInfo *bot_common.OnboardingInfo) (string, error) {
	onboarding := playground.OnboardingContent{}
	if onboardingInfo != nil {
		onboarding.Prologue = ptr.Of(onboardingInfo.GetPrologue())
		onboarding.SuggestedQuestions = onboardingInfo.GetSuggestedQuestions()
		onboarding.SuggestedQuestionsShowMode = onboardingInfo.SuggestedQuestionsShowMode
	}

	onboardingInfoStr, err := sonic.MarshalString(onboarding)
	if err != nil {
		return "", err
	}

	return onboardingInfoStr, nil
}
