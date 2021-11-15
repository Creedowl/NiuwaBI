package main

import (
	"fmt"
	"github.com/Creedowl/NiuwaBI/database"
	"github.com/Creedowl/NiuwaBI/utils"
	"github.com/sirupsen/logrus"
)

func init() {
	utils.InitConfig()
	utils.InitLogger()
	database.InitDB()
}

func main() {
	logrus.Infoln("start loading app")
	app := InitApp()

	app.Run(fmt.Sprintf("%s:%d", utils.Cfg.Host, utils.Cfg.Port))
}
