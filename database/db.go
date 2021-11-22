package database

import (
	"fmt"
	"github.com/Creedowl/NiuwaBI/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	logrus.Infoln("init db")
	dbCfg := utils.Cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DefaultDB)
	logrus.Debugf("mysql dsn: %s\n", dsn)
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
		logrus.Fatalf("failed to init db: %+v\n", err)
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
