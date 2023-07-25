package router

import (
	"brief/pkg/handler/user"
	mdw "brief/pkg/middleware"
	"brief/pkg/repository/storage/postgres"
	userSrv "brief/service/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

// User registers user paths with router 'e'
func User(r chi.Router, validate *validator.Validate, logger *log.Logger) chi.Router {

	// Use postgres database
	pgDb := postgres.GetDB()
	uService := userSrv.NewUserService(pgDb)
	userCtrl := user.NewController(validate, logger, uService)

	// Create admin user
	uService.CreateAdminUser(logger)

	// Free endpoints
	r.Group(func(r chi.Router) {
		r.Post("/users", userCtrl.Register)
		r.Post("/users/login", userCtrl.Login)
	})

	// User endpoints
	r.Group(func(r chi.Router) {
		r.Use(mdw.Me) // user middleware

		r.Get("/users", userCtrl.GetMe)
		r.Patch("/users", userCtrl.UpdateMe)
		r.Patch("/users/reset-password", userCtrl.ResetPassword)
	})

	// Admin endpoints
	r.Group(func(r chi.Router) {
		r.Use(mdw.Admin) // admin middleware

		r.Get("/users/get-all", userCtrl.GetAll)
		r.Get("/users/{idOrEmail}", userCtrl.GetUserByIdOrEmail)
		r.Patch("/users/lock/{idOrEmail}", userCtrl.LockUser)
		r.Patch("/users/unlock/{idOrEmail}", userCtrl.UnlockUser)
	})

	return r
}
