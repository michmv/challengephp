package storage

import (
	"challengephp/lib"
	"challengephp/src/config"
	"challengephp/src/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Debug  bool
	Logger lib.Logger
	Config config.Config
	Pool   *pgxpool.Pool
}

func New(configPath string) (Storage, lib.Error) {
	res := Storage{}
	conf, err := config.LoadConfig(configPath)
	if err != nil {
		return res, err.Tap()
	}

	logger, err := lib.NewFileLogger(conf.LogFile)
	if err != nil {
		return res, err.Tap()
	}

	pool, err := db.CreateDB(conf.DB)
	if err != nil {
		lib.Exit(err.Tap())
	}

	return Storage{
		Config: conf,
		Debug:  conf.Debug,
		Logger: logger,
		Pool:   pool,
	}, nil
}

func (it Storage) Close() {
	it.Logger.Close()
	it.Pool.Close()
}
