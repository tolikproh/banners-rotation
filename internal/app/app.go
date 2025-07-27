package app

import (
	"context"
	"fmt"
	"time"

	"github.com/tolikproh/banners-rotation/internal/algorithm"
	"github.com/tolikproh/banners-rotation/internal/cnst"
	"github.com/tolikproh/banners-rotation/internal/config"
	"github.com/tolikproh/banners-rotation/internal/logger"
	"github.com/tolikproh/banners-rotation/internal/model"
)

type App struct {
	cfg                  *config.Config
	log                  *logger.Logger
	Storage              Storage
	NotificationReceiver NotificationReceiver
}

type NotificationReceiver interface {
	Add(model.Event) error
}

func New(cfg *config.Config, log *logger.Logger, storage Storage, rcv NotificationReceiver) *App {
	return &App{
		cfg:                  cfg,
		log:                  log,
		Storage:              storage,
		NotificationReceiver: rcv,
	}
}

func (a *App) GetBanner(ctx context.Context, slotID, socialGroupID int) (int, error) {
	banners, err := a.Storage.GetBannersInfo(ctx, slotID, socialGroupID)

	var bannerID int
	if err == nil {
		if bannerID, err = algorithm.Bandit(banners); err != nil {
			return 0, err
		}
	} else {
		bannerID, err = a.Storage.GetRandomBanner(ctx, slotID)
		if err != nil {
			a.log.Error("get banner for slot", "error", err)
			return 0, err
		}
	}

	if bannerID == 0 {
		return 0, cnst.ErrGetBanner
	}

	err = a.Storage.IncrementBannerView(ctx, bannerID, slotID, socialGroupID)
	if err != nil {
		a.log.Error("error add banner view", "error", err)
	}

	// event
	statEvent := model.Event{
		EventType:     "view",
		SlotID:        slotID,
		BannerID:      bannerID,
		SocialGroupID: socialGroupID,
		DateTime:      time.Now(),
	}

	if err := a.NotificationReceiver.Add(statEvent); err != nil {
		a.log.Error("add statEvent for event 'view'", "error", err)
	} else {
		a.log.Info("Event 'view' on banner", "id", bannerID)
	}

	return bannerID, nil
}

func (a *App) AddBannerToSlot(ctx context.Context, bannerID, slotID int) error {
	if err := a.Storage.AddBannerToSlot(ctx, bannerID, slotID); err != nil {
		a.log.Error("add banner to slot", "error", err)
		return err
	}

	return nil
}

func (a *App) ClickBanner(ctx context.Context, bannerID, slotID, socialGroupID int) error {
	if err := a.Storage.ClickBanner(ctx, bannerID, slotID, socialGroupID); err != nil {
		a.log.Error("click banner from slot", "error", err)
		return err
	}

	statEvent := model.Event{
		EventType:     "click",
		SlotID:        slotID,
		BannerID:      bannerID,
		SocialGroupID: socialGroupID,
		DateTime:      time.Now(),
	}

	if err := a.NotificationReceiver.Add(statEvent); err != nil {
		return fmt.Errorf("error add statEvent for event 'click':  %w", err)
	}

	a.log.Info("Event 'click' on banner", "id", bannerID)

	return nil
}

func (a *App) DeleteBannerFromSlot(ctx context.Context, bannerID, slotID int) error {
	if err := a.Storage.DeleteBannerFromSlot(ctx, bannerID, slotID); err != nil {
		a.log.Error("delete banner from slot", "error", err)
		return err
	}

	return nil
}

func (a *App) CreateBanner(ctx context.Context, desc string) (int, error) {
	return a.Storage.CreateBanner(ctx, desc)
}

func (a *App) CreateSlot(ctx context.Context, desc string) (int, error) {
	return a.Storage.CreateSlot(ctx, desc)
}

func (a *App) CreateSocialGroup(ctx context.Context, desc string) (int, error) {
	return a.Storage.CreateSocialGroup(ctx, desc)
}

func (a *App) DeleteBanner(ctx context.Context, id int) error {
	return a.Storage.DeleteBanner(ctx, id)
}

func (a *App) DeleteSlot(ctx context.Context, id int) error {
	return a.Storage.DeleteSlot(ctx, id)
}

func (a *App) DeleteSocialGroup(ctx context.Context, id int) error {
	return a.Storage.DeleteSocialGroup(ctx, id)
}
