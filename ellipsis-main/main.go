package main

import (
	pgdb "brief/pkg/repository/storage/postgres"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"brief/internal/config"
	"brief/pkg/router"

	"github.com/go-playground/validator/v10"
	// rdb "brief/pkg/repository/storage/redis"
)

func init() {
	config.Setup()
	pgdb.ConnectToDB()
	// redis.SetupRedis() uncomment when you need redis
}

//	@title			Brief
//	@version		1.0
//	@description	URL Shortener.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		brief.up.railway.app
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	JWTToken
//	@in							header
//	@name						Authorization
//	@description				JWT token

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func main() {
	//Load config
	logger := log.New()
	getConfig := config.GetConfig()
	validatorRef := validator.New()
	e := router.Setup(validatorRef, logger)

	// The HTTP Server
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", getConfig.ServerPort),
		Handler: e,
	}

	// Server run context
	serverCtx, serverCancel := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Store counter variable in redis
		// redis.StoreCounter()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		shutdownCancel()
		serverCancel()
	}()

	// Run the server
	fmt.Printf("Server is now listening on port: %s\n", getConfig.ServerPort)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
