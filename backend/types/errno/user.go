package errno

import (
	"github.com/kiosk404/airi-go/backend/pkg/errorx/code"
)

// Passport: 700 000 000 ~ 700 999 999
const (
	ErrUserAuthenticationFailed = 700012006 // Don't change this code. It is used in the frontend.

	ErrUserAccountAlreadyExistCode    = 700000001
	ErrUserUniqueNameAlreadyExistCode = 700000002
	ErrUserInfoInvalidateCode         = 700000003
	ErrUserSessionInvalidateCode      = 700000004
	ErrUserResourceNotFound           = 700000005
	ErrUserInvalidParamCode           = 700000006
	ErrUserPermissionCode             = 700000007
	ErrNotAllowedRegisterCode         = 700000008
)

func init() {

	code.Register(
		ErrNotAllowedRegisterCode,
		"The user registration has been disabled by the administrator. Please contact the administrator!",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrUserPermissionCode,
		"unauthorized access : {msg}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrUserInvalidParamCode,
		"invalid parameter : {msg}",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrUserResourceNotFound,
		"{type} not found: {id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrUserInfoInvalidateCode,
		"invalid email or password, please try again.",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrUserUniqueNameAlreadyExistCode,
		"unique name already exist : {name}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrUserEmailAlreadyExistCode,
		"email already exist : {email}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrUserAuthenticationFailed,
		"authentication failed: {reason}",
		code.WithAffectStability(false),
	)
}
