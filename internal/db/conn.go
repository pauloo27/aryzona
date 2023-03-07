package db

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/db/model"
	// postgres driver
	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

var (
	models = []any{
		new(model.Guild),
		new(model.User),
	}

	DB *DBConn
)

type DBConn struct {
	engine *xorm.Engine
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSL      bool
}

func NewDB(config *DBConfig) (*DBConn, error) {
	sslMode := "disable"
	if config.SSL {
		sslMode = "enable"
	}
	engine, err := xorm.NewEngine(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.Database, sslMode,
		),
	)
	if err != nil {
		return nil, err
	}
	return &DBConn{
		engine: engine,
	}, nil
}

func (db *DBConn) Migrate() error {
	return db.engine.Sync(models...)
}

func (db *DBConn) Close() error {
	return db.engine.Close()
}
