package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"homapage-i18n/config"
	"homapage-i18n/log"
	"homapage-i18n/mongodb"
	"homapage-i18n/quit"
	"homapage-i18n/routes"
	"homapage-i18n/server"
	"homapage-i18n/token"
)

func main() {
	defer logrus.Info("Main exiting")
	config.InitConfig()
	log.InitLogger()
	mongodb.ConnectDB()
	defer mongodb.DisconnectDB()
	token.InitVerifyKey(viper.GetString("token.publicKeyPath"))

	r := routes.SetupRouter()
	srv := server.StartServer(r)

	quit.WaitForQuitSignal()

	server.ShutdownServer(srv)

}
