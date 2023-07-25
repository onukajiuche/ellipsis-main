package postgres

import (
	"brief/internal/model"
	"context"
	"strings"

	"gorm.io/gorm/clause"
)

// CreateUser stores 'user' in the database
func (p *Postgres) CreateUser(ctx context.Context, user *model.User) error {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()
	return db.Create(user).Error
}

// GetUser fetches a user from the database using its 'id' or 'email'
func (p *Postgres) GetUser(ctx context.Context, idOrEmail string) (*model.User, error) {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()

	// check if 'idOrEmail' is an email or an id
	var cond string
	if strings.Contains(idOrEmail, "@") {
		cond = "email = ?"
	} else {
		cond = "id = ?"
	}

	var user model.User
	err := db.First(&user, cond, idOrEmail).Error
	return &user, err
}

// GetAllUsers gets all users in the database
func (p *Postgres) GetAllUsers(ctx context.Context) ([]model.User, error) {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()

	var users []model.User
	err := db.Find(&users).Error
	return users, err
}

// UpdateUser updates some specific user fields
func (p *Postgres) UpdateUser(ctx context.Context, id string, user *model.User) error {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()

	// Ensure 'is_verified', 'is_locked', 'password', 'email' and 'salt' cannot be updated using this function
	return db.Model(user).Clauses(clause.Returning{}).
		Omit("is_locked", "password", "salt", "email", "role").
		Where("id = ?", id).Updates(user).Error
}

// ResetPassword resets the 'password' and 'salt' of a user with 'id'
func (p *Postgres) ResetPassword(ctx context.Context, id string, rp *model.ResetPassword) (*model.User, error) {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()

	var user model.User
	user.ID = id
	err := db.Model(&user).Clauses(clause.Returning{}).Select("Password", "Salt", "UpdatedAt").
		Updates(model.User{Password: rp.Password, Salt: rp.Salt}).Error

	return &user, err
}

// LockUnlock sets the 'is_locked' field of a user to 'false' or 'false'
func (p *Postgres) LockUnlock(ctx context.Context, idOrEmail string, isLocked bool) (*model.User, error) {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()

	// check if 'idOrEmail' is an email or an id
	var cond string
	if strings.Contains(idOrEmail, "@") {
		cond = "email = ?"
	} else {
		cond = "id = ?"
	}

	var user model.User
	err := db.Model(&user).Clauses(clause.Returning{}).Where(cond, idOrEmail).
		Update("is_locked", isLocked).Error

	return &user, err
}
