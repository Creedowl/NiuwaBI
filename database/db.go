package database

import (
	"fmt"
	"github.com/Creedowl/NiuwaBI/utils"
	lru "github.com/hashicorp/golang-lru"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

var DB *gorm.DB

var Pool *lru.Cache

func InitDB() {
	logrus.Infoln("init db")
	dbCfg := utils.Cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DefaultDB)
	logrus.Debugf("mysql dsn: %s", dsn)
	var err error
	customLogger := logger.New(
		logrus.StandardLogger(),
		logger.Config{
			SlowThreshold:             0,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Info,
		})
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: customLogger})
	if err != nil {
		logrus.Fatalf("failed to init db: %+v", err)
	}

	Pool, err = lru.New(10)
	if err != nil {
		logrus.Fatalf("failed to create lru pool: %+v", err)
	}
}

func GetDB() *gorm.DB {
	db := DB.Session(&gorm.Session{
		NewDB: true,
	})
	if utils.Cfg.Debug {
		return db.Debug()
	}
	return db
}

func (c *DBConfig) TestConn() (*gorm.DB, error) {
	switch c.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Username, c.Password, c.Host, c.Port, c.Database)
		logrus.Debugf("dsn: %s", dsn)
		db, err := gorm.Open(mysql.Open(dsn))
		if err != nil {
			return nil, err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		return db, sqlDB.Ping()
	default:
		return nil, fmt.Errorf("unsupport database type: %s", c.Type)
	}
}

func GetCachedDB(key uint, c *DBConfig) (*gorm.DB, error) {
	db, ok := Pool.Get(key)
	if !ok {
		var err error
		db, err = c.TestConn()
		if err != nil {
			return nil, err
		}
		Pool.Add(key, db)
	}
	return db.(*gorm.DB), nil
}
