package storage

import (
	"brief/internal/model"
	"context"
)

// repositories
type StorageRepository interface {
	// User
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, idOrEmail string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
	UpdateUser(ctx context.Context, id string, user *model.User) error
	ResetPassword(ctx context.Context, id string, rp *model.ResetPassword) (*model.User, error)
	LockUnlock(ctx context.Context, idOrEmail string, isLocked bool) (*model.User, error)

	// URL
	CreateURL(ctx context.Context, url *model.URL) error
	GetURL(ctx context.Context, hash string) (*model.URL, error)
	GetURLById(ctx context.Context, id string) (*model.URL, error)
	GetUrls(ctx context.Context, userID string) ([]model.URL, error)
	GetAll(ctx context.Context) ([]model.URL, error)
	DeleteUrl(ctx context.Context, id string) (*model.URL, error)
}

type RedisRepository interface {
	RedisSet(key string, value interface{}) error
	RedisGet(key string) ([]byte, error)
	RedisDelete(key string) (int64, error)
}

// repositories
