package user

import (
	"brief/service/user"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	Validate    *validator.Validate
	Logger      *log.Logger
	UserService user.UserService
}

func NewController(validate *validator.Validate, logger *log.Logger, uService user.UserService) *Controller {
	return &Controller{
		validate, logger, uService,
	}
}
