package repository

import (
	"context"

	"github.com/kiosk404/airi-go/backend/domain/user/internal/dal"
	"github.com/kiosk404/airi-go/backend/domain/user/internal/dal/model"
	"gorm.io/gorm"
)

func NewUserRepo(db *gorm.DB) UserRepository {
	return dal.NewUserDAO(db)
}

type UserRepository interface {
	GetUsersByAccount(ctx context.Context, account string) (*model.User, bool, error)
	UpdateSessionKey(ctx context.Context, userID int64, sessionKey string) error
	ClearSessionKey(ctx context.Context, userID int64) error
	UpdatePassword(ctx context.Context, account, password string) error
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	UpdateAvatar(ctx context.Context, userID int64, iconURI string) error
	CheckUniqueNameExist(ctx context.Context, uniqueName string) (bool, error)
	UpdateProfile(ctx context.Context, userID int64, updates map[string]any) error
	CheckAccountExist(ctx context.Context, account string) (bool, error)
	CreateUser(ctx context.Context, user *model.User) error
	GetUserBySessionKey(ctx context.Context, sessionKey string) (*model.User, bool, error)
	GetUsersByIDs(ctx context.Context, userIDs []int64) ([]*model.User, error)
}
