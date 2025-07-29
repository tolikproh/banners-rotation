package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/tolikproh/banners-rotation/internal/app"
	"github.com/tolikproh/banners-rotation/internal/brokers/rabbitmq"
	"github.com/tolikproh/banners-rotation/internal/config"
	"github.com/tolikproh/banners-rotation/internal/logger"
	ihttp "github.com/tolikproh/banners-rotation/internal/server/http"
	"github.com/tolikproh/banners-rotation/internal/storage"
)

func run() {
	pathConf := filepath.Dir(configFile)
	fileConf := filepath.Base(configFile)

	cfg := config.NewConfig(pathConf, fileConf)
	log := logger.New(cfg.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	rcv, err := rabbitmq.NewRabbit(ctx, cfg, log)
	if err != nil {
		log.Error("create rabbit client", "error", err)
		return
	}

	storage, err := storage.New(ctx, cfg, log)
	if err != nil {
		log.Error("connect to database", "error", err)
		return
	}

	bannerRotation := app.New(cfg, log, storage, rcv)
	serverHTTP := ihttp.New(cfg, log, bannerRotation)

	log.Debug("set loger level: " + cfg.Logger.Level)
	log.Debug("storage connection: " + cfg.Storage.Conn)

	go func() {
		if err := serverHTTP.Start(ctx); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info("http server stopped")
				return
			}
			log.Error("failed start http server", "error", err)
		}
	}()

	log.Info("banner rotation is running...")

	<-ctx.Done()

	if err := serverHTTP.Stop(ctx); err != nil {
		log.Error("http server stop", "error", err.Error())
	}
}
