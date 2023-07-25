package user

import (
	"brief/internal/config"
	"brief/internal/constant"
	"brief/internal/model"
	"brief/pkg/repository/storage/postgres"
	"brief/utility"
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// CreateAdminUser creates an admin user if one doesn't exist
func (u *userService) CreateAdminUser(logger *log.Logger) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	getConfig := config.GetConfig()
	password, salt, err := utility.HashPassword(getConfig.AdminPassword)
	if err != nil {
		return fmt.Errorf("could not hash admin password, got error: %w", err)
	}

	db := postgres.GetDB()

	// Check if admin user already exists
	if _, err = db.GetUser(ctx, getConfig.AdminID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create admin user if it doesn't exist
			err := db.CreateUser(ctx, &model.User{
				ID:       getConfig.AdminID,
				Password: password,
				Salt:     salt,
				Email:    getConfig.AdminID + "@email.com",
				Role:     constant.Roles[constant.Admin],
			})
			if err != nil {
				return fmt.Errorf("could not create admin user, got error: %w", err)
			}
		} else {
			return fmt.Errorf("could not check admin user, got error: %w", err)
		}
	} else {
		logger.Info("ADMIN USER ALREADY EXISTS")
		return nil
	}

	logger.Info("CREATED ADMIN USER SUCCESSFULLY")
	return nil
}
