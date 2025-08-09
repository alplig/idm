package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"idm/inner/common"
	"idm/utils"
	"path/filepath"
	"time"
)

type Env string

const (
	DevEnv  Env = ".env"
	TestEnv Env = ".env_test"
)

// ConnectDb получить конфиг и подключиться с ним к базе данных
func ConnectDb(env Env) *sqlx.DB {
	cfg, err := loadConfig(env)
	if err != nil {
		panic(err)
	}
	return ConnectDbWithCfg(cfg)
}

func loadConfig(envType Env) (common.Config, error) {
	pathRoot, _ := utils.FindRoot()
	switch envType {
	case DevEnv:
		{
			path := filepath.Join(pathRoot, string(DevEnv))
			return common.GetConfig(path), nil
		}
	case TestEnv:
		{
			path := filepath.Join(pathRoot, string(TestEnv))
			return common.GetConfig(path), nil
		}
	default:
		return common.Config{}, fmt.Errorf("unknown environment: %s", envType)
	}
}

// ConnectDbWithCfg подключиться к базе данных с переданным конфигом
func ConnectDbWithCfg(cfg common.Config) *sqlx.DB {
	var db = sqlx.MustConnect(cfg.DbDriverName, cfg.Dsn)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(1 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	return db
}
