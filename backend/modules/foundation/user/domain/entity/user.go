package entity

type User struct {
	UserID int64

	Name         string // nickname
	UniqueName   string // unique name
	Account      string // account
	Description  string // user description
	IconURI      string // avatar URI
	IconURL      string // avatar URL
	UserVerified bool   // Is the user authenticated?
	Locale       string
	SessionKey   string // session key

	CreatedAt int64 // creation time
	UpdatedAt int64 // update time
}
