package main

import (
	"battle-ship_server/internal/config"
	"battle-ship_server/internal/rabbitmq"
	"battle-ship_server/internal/service/auth"
	"battle-ship_server/internal/service/game"
	"battle-ship_server/internal/storage/postgres"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad(os.Getenv("CONFIG_PATH"))

	log := setupLogger(cfg.Env)

	log.Info("Starting server")

	postgresUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName) 
	storage, err := postgres.New(postgresUrl)
	if err != nil {
		panic(err)
	}

	auth := auth.New(storage, log)
	game := game.New(log) 

	rmqUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
	rmq := rabbitmq.New(rmqUrl, log, auth, game)

	go rmq.Login()
	go rmq.Register()

	log.Info("Server started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	rmq.Close()
	storage.Close()
	log.Info("Gracefully stopped")
	
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
