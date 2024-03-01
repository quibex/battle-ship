package main

import (
	"battlship/internal/adapters/rabbitmq"
	"battlship/internal/config"
	"battlship/internal/service/authSrvc"
	gameSrvs "battlship/internal/service/game"
	terminalUI "battlship/internal/ui/terminal"
	authUI "battlship/internal/ui/terminal/auth"
	gameUI "battlship/internal/ui/terminal/game"
)

func main() {
	cfg := config.MustLoad("config/config.yaml")

	//fmt.Println(cfg.TimeOut)

	rmq := rabbitmq.New(cfg.RmqURL, cfg.TimeOut)

	auth := authSrvc.New(rmq)
	game := gameSrvs.New(rmq)

	authUI1 := authUI.New(auth)
	gameUI1 := gameUI.New(game)

	ui := terminalUI.New(authUI1, gameUI1)

	ui.MustRun()
}
