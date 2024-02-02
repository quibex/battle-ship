package main

import (
	"battlship/internal/config"
	"battlship/internal/rabbitmq"
	"battlship/internal/service/auth"
	"battlship/internal/ui"
)

func main() {
	cfg := config.MustLoad("config/config.yaml")

	rmq := rabbitmq.New(cfg.RmqURL)

	auth := auth.New(rmq)

	ui := ui.New(auth)

	ui.MustRun()
}
