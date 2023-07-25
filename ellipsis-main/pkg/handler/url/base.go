package url

import (
	"brief/service/url"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	Validate   *validator.Validate
	Logger     *log.Logger
	UrlService url.UrlService
}

func NewController(validate *validator.Validate, logger *log.Logger, uService url.UrlService) *Controller {
	return &Controller{
		validate, logger, uService,
	}
}
