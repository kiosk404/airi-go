package entity

// Reuse DomainName and OpType in AppDomainEvent

type ResourceDomainEvent struct {
	OpType   OpType            `json:"op_type"`
	Resource *ResourceDocument `json:"resource_document,omitempty"`
	Meta     *EventMeta        `json:"meta,omitempty"`
	Extra    map[string]any    `json:"extra"`
}
