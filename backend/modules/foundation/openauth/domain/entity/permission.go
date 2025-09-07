package entity

type (
	ResourceType int
	Decision     int
)

type ResourceIdentifier struct {
	Type ResourceType
	ID   string
}

type ActionAndResource struct {
	Action             string
	ResourceIdentifier ResourceIdentifier
}

type CheckPermissionRequest struct {
	IdentityTicket     string
	ActionAndResources []ActionAndResource
}

type CheckPermissionResponse struct {
	Decision Decision
}
