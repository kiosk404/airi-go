package entity

func MergeStreamMsgs(msgs []*Message) *Message {
	if len(msgs) == 0 {
		return nil
	}
	allInOne := &Message{
		Role:         RoleAssistant,
		ResponseMeta: &ResponseMeta{},
	}
	for _, msg := range msgs {
		if msg == nil {
			continue
		}
		allInOne.Content += msg.Content
		allInOne.ReasoningContent += msg.ReasoningContent
		if msg.MultiModalContent != nil {
			allInOne.MultiModalContent = append(allInOne.MultiModalContent, msg.MultiModalContent...)
		}
		if msg.ResponseMeta != nil {
			if msg.ResponseMeta.Usage != nil {
				allInOne.ResponseMeta.Usage = msg.ResponseMeta.Usage
			}
			if msg.ResponseMeta.FinishReason != "" {
				allInOne.ResponseMeta.FinishReason = msg.ResponseMeta.FinishReason
			}
		}
		// 流式tool call聚合
		for _, tc := range msg.ToolCalls {
			if tc == nil {
				continue
			}
			if len(allInOne.ToolCalls) == 0 || isFirstStreamPkgToolCall(allInOne.ToolCalls[len(allInOne.ToolCalls)-1], tc) {
				allInOne.ToolCalls = append(allInOne.ToolCalls, &ToolCall{
					Index:    tc.Index,
					ID:       tc.ID,
					Type:     tc.Type,
					Function: tc.Function,
					Extra:    tc.Extra,
				})
			} else {
				allInOne.ToolCalls[len(allInOne.ToolCalls)-1].Function.Arguments += tc.Function.Arguments
				allInOne.ToolCalls[len(allInOne.ToolCalls)-1].Extra = tc.Extra
			}
		}
	}
	return allInOne
}

func isFirstStreamPkgToolCall(tc1, tc2 *ToolCall) bool {
	if tc1 == nil || tc2 == nil {
		return false
	}
	if tc2.ID == "" {
		return false
	}
	return true
}
