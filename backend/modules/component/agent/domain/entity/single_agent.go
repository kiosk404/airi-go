package entity

import (
	model "github.com/kiosk404/airi-go/backend/modules/component/crossdomain/agent/model"
)

// SingleAgent Use composition instead of aliasing for domain entities to enhance extensibility
type SingleAgent struct {
	*model.SingleAgent
}

type AgentIdentity = model.AgentIdentity

type ExecuteRequest = model.ExecuteRequest

type DuplicateInfo struct {
	UserID     int64
	NewAgentID int64
	DraftAgent *SingleAgent
}
