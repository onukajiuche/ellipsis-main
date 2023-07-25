package postgres

import (
	"brief/internal/config"
	"brief/internal/model"
	"brief/pkg/repository/storage"
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

var (
	db                  *gorm.DB
	generalQueryTimeout = 60 * time.Second
)

func GetDB() storage.StorageRepository {
	return &Postgres{db}
}

func ConnectToDB() *gorm.DB {
	logger := log.New()

	database, err := gorm.Open(postgres.Open(dsn()), &gorm.Config{})
	if err != nil {
		logger.Fatalf("could not connect to postgres, got error: %s", err)
	}
	db = database

	if err := migrateDB(logger); err != nil {
		logger.Fatalf("could not run db migrations, got error: %s", err)
	}

	// IF EVERYTHING IS OKAY, THEN CONNECTION IS ESTABLISHED
	logger.Info("POSTGRES CONNECTION ESTABLISHED")

	return db
}

func dsn() string {
	pgHost := config.GetConfig().PGHost
	pgPort := config.GetConfig().PGPort
	pgUser := config.GetConfig().PGUser
	pgDB := config.GetConfig().PGDatabase
	pgPassword := config.GetConfig().PGPassword
	pgSSL := config.GetConfig().PGSSLMode

	dsn := "host=" + pgHost + " user=" + pgUser +
		" password=" + pgPassword + " dbname=" + pgDB + " port=" + pgPort + " sslmode=" + pgSSL

	return dsn
}

// migrateDB creates db schemas
func migrateDB(logger *log.Logger) error {
	err := db.AutoMigrate(
		&model.User{},
		&model.URL{},
	)
	if err != nil {
		return err
	}

	logger.Info("DATABASE MIGRATION SUCCESSFUL")
	return nil
}

// DBWithTimeout returns a database with timeout, and the context's cancel func
func (p *Postgres) DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout)
	return p.db.WithContext(ctx), cancel
}
