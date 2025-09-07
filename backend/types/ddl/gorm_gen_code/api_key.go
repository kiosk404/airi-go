package main

type ApiKey struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ApiKey      string `json:"api_key"`
	ConnectorID int64  `json:"connector_id"`
	UserID      int64  `json:"user_id"`
	AkType      string `json:"ak_type"`
	LastUsedAt  int64  `json:"last_used_at"`
	ExpiredAt   int64  `json:"expired_at"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}
