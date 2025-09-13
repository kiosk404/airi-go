package entity

import (
	"github.com/kiosk404/airi-go/backend/api/crossdomain/singleagent"
)

// SingleAgent Use composition instead of aliasing for domain entities to enhance extensibility
type SingleAgent struct {
	*singleagent.SingleAgent
}

type AgentIdentity = singleagent.AgentIdentity

type ExecuteRequest = singleagent.ExecuteRequest

type DuplicateInfo struct {
	UserID     int64
	NewAgentID int64
	DraftAgent *SingleAgent
}
