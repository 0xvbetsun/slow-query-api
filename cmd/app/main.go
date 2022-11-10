package main

import (
	"log"

	"github.com/vbetsun/slow-query-api/configs"
	"github.com/vbetsun/slow-query-api/internal/storage/psql"
	"github.com/vbetsun/slow-query-api/internal/transport/rest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := configs.New()
	if err != nil {
		log.Fatal("failed to make config: ", err)
	}

	dbConf := postgres.New(postgres.Config{
		DSN:                  cfg.PostgresDSN,
		PreferSimpleProtocol: true,
	})

	dbConn, err := gorm.Open(dbConf)
	if err != nil {
		log.Fatal("DB connection ERROR: ", err)
	}

	repo := psql.NewPgStatStatementRepo(dbConn)
	handler := rest.NewHandler(repo, cfg.SlowQueryDuration)
	app := new(rest.Server)

	log.Println("Exit", app.Run(":"+cfg.Port, handler))
}
