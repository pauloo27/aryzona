package db

import (
	"fmt"

	// postgres driver
	_ "github.com/lib/pq"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var (
	DB *DBConn
)

type DBConn struct {
	*xorm.Engine
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
	engine.Logger().SetLevel(log.LOG_ERR)
	return &DBConn{
		Engine: engine,
	}, nil
}

func (db *DBConn) Close() error {
	return db.Engine.Close()
}
