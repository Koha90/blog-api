package main

import (
	"github.com/koha90/blog-api/config"
	"github.com/koha90/blog-api/pkg/apiserver"
	"github.com/koha90/blog-api/pkg/logger"
	"github.com/koha90/blog-api/pkg/storage"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err.Error())
	}

	logger := logger.GetLogger()

	store, err := storage.NewPostgresStore(cfg.Postgres.URL)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Close()

	if err := store.Init(); err != nil {
		logger.Fatal(err)
	}

	logger.Infof("APIServer run in port: %s", cfg.HTTP.Port)

	apiServer := apiserver.New(cfg.HTTP.Port, store)
	if err = apiServer.Run(); err != nil {
		logger.Panic(err)
	}
}
