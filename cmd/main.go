package main

import (
	"context"
	"flag"
	"fmt"
	app2 "github.com/BAN1ce/skyTree/app"
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title           SkyTree API
// @version         1.0

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func main() {
	flag.Parse()
	if err := config.LoadConfig(); err != nil {
		panic(fmt.Errorf("load config error %v", err))
	}
	logger.Load()
	logger.Logger.Debug().Any("config", config.GetConfig()).Msg("config loaded")

	var (
		app         = app2.NewApp()
		ctx, cancel = context.WithCancel(context.Background())
	)

	fmt.Println("Starting with version: ", app.Version())

	app.Start(ctx)

	fmt.Println(logo())
	fmt.Println("Wish you enjoy it!")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	exitTimeout(32 * time.Second)
	cancel()

	if err := app.Close(); err != nil {
		logger.Logger.Error().Err(err).Msg("app close error")
	}
	logger.Logger.Info().Msg("exit successfully")
	os.Exit(0)
}

func exitTimeout(t time.Duration) {
	go func() {
		time.Sleep(t)
		logger.Logger.Error().Msg("exit timeout")
		os.Exit(1)
	}()

}
func logo() string {
	return `
   _____ _       _______
  / ____| |     |__   __|
 | (___ | | ___   _| |_ __ ___  ___
  \___ \| |/ / | | | | '__/ _ \/ _ \
  ____) |   <| |_| | | | |  __/  __/
 |_____/|_|\_\\__, |_|_|  \___|\___|
               __/ |
              |___/
`
}
