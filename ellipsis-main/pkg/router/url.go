package router

import (
	"brief/pkg/handler/url"
	mdw "brief/pkg/middleware"
	"brief/pkg/repository/storage/postgres"
	urlSrv "brief/service/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

// Redirect registers the redirect path for a short url
func Redirect(r chi.Router, validate *validator.Validate, logger *log.Logger) chi.Router {
	// Use postgres database
	pgDb := postgres.GetDB()
	uService := urlSrv.NewUrlService(pgDb)
	urlCtrl := url.NewController(validate, logger, uService)

	r.Group(func(r chi.Router) {
		r.Get("/{hash}", urlCtrl.Redirect)
	})

	return r
}

// Url registers url paths with router 'r'
func Url(r chi.Router, validate *validator.Validate, logger *log.Logger) chi.Router {

	// Use postgres database
	pgDb := postgres.GetDB()
	uService := urlSrv.NewUrlService(pgDb)
	urlCtrl := url.NewController(validate, logger, uService)

	// Shorten endpoint
	r.Group(func(r chi.Router) {
		r.Use(mdw.Shorten)
		r.Post("/url/shorten", urlCtrl.Shorten)
	})

	// User endpoints
	r.Group(func(r chi.Router) {
		r.Use(mdw.Me) // user middleware

		r.Get("/url", urlCtrl.GetUrls)
		r.Delete("/url/{id}", urlCtrl.Delete)
	})

	// Admin endpoints
	r.Group(func(r chi.Router) {
		r.Use(mdw.Admin) // admin middleware

		r.Get("/url/get-all", urlCtrl.GetAll)
		r.Get("/url/{user-id}", urlCtrl.GetUrlsByUserID)
	})

	return r
}
