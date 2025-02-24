package main

import (
	"database/sql"
	"go.uber.org/zap"
	"postery/internal/db"
	"postery/internal/env"
	"postery/internal/store"
	"time"
)

const version = "0.0.1"

//	@title			PosterySocial API
//	@version		1.0
//	@description	API for PosterySocial, a social network for posters
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/v2
func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         "postgres://admin:adminpassword@localhost/social?sslmode=disable",
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  "15m",
		},
		env:    env.GetString("ENV", "development"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080/api"),
		mail: mailConfig{
			exp: time.Hour * 24 * 3, // 3 days
		},
	}

	// logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Fatal("Error when closing DB: ", err)
		}
	}(db)

	logger.Info("DB connection pool established")

	appStore := store.NewStorage(db)

	app := &application{
		store:  appStore,
		config: cfg,
		logger: logger,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
