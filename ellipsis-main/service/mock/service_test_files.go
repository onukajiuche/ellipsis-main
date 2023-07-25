package mock

import (
	"brief/internal/model"
	"context"
	"fmt"
)

type Repo struct{}

// User

func (r *Repo) CreateUser(ctx context.Context, user *model.User) error {
	fmt.Println("Hit CreateUser repo function...")
	return nil
}

func (r *Repo) GetUser(ctx context.Context, idOrEmail string) (*model.User, error) {
	fmt.Println("Hit GetUser repo function...")
	return &model.User{ID: idOrEmail, Email: idOrEmail}, nil
}

func (r *Repo) GetAllUsers(ctx context.Context) ([]model.User, error) {
	fmt.Println("Hit GetAllUsers repo function...")
	return []model.User{}, nil
}

func (r *Repo) UpdateUser(ctx context.Context, id string, user *model.User) error {
	fmt.Println("Hit UpdateUser repo function...")
	return nil
}

func (r *Repo) ResetPassword(ctx context.Context, id string, rp *model.ResetPassword) (*model.User, error) {
	fmt.Println("Hit ResetPassword repo function...")
	return &model.User{ID: id, Password: rp.Password, Salt: rp.Salt}, nil
}

func (r *Repo) LockUnlock(ctx context.Context, idOrEmail string, isLocked bool) (*model.User, error) {
	fmt.Println("Hit LockUnlock repo function...")
	return &model.User{ID: idOrEmail, Email: idOrEmail, IsLocked: isLocked}, nil
}

// URL

func (r *Repo) CreateURL(ctx context.Context, url *model.URL) error {
	fmt.Println("Hit CreateURL repo function...")
	return nil
}

func (r *Repo) GetURL(ctx context.Context, hash string) (*model.URL, error) {
	fmt.Println("Hit GetURL repo function...")
	return &model.URL{Hash: hash}, nil
}

func (r *Repo) GetURLById(ctx context.Context, id string) (*model.URL, error) {
	fmt.Println("Hit GetURLById repo function...")
	return &model.URL{ID: id, UserID: id}, nil
}

func (r *Repo) GetUrls(ctx context.Context, userID string) ([]model.URL, error) {
	fmt.Println("Hit GetUrls repo function...")
	return []model.URL{{UserID: userID}}, nil
}

func (r *Repo) GetAll(ctx context.Context) ([]model.URL, error) {
	fmt.Println("Hit GetAll repo function...")
	return []model.URL{}, nil
}

func (r *Repo) DeleteUrl(ctx context.Context, id string) (*model.URL, error) {
	fmt.Println("Hit DeleteUrl repo function...")
	return &model.URL{ID: id}, nil
}
