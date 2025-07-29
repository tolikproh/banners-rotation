package app

import (
	"context"

	"github.com/tolikproh/banners-rotation/internal/model"
)

type Storage interface {
	GetBannersInfo(ctx context.Context, slotID, socialGroupID int) ([]model.Banner, error)
	AddBannerToSlot(ctx context.Context, bannerID, slotID int) error
	ClickBanner(ctx context.Context, bannerID, slotID, socialGroupID int) error
	DeleteBannerFromSlot(ctx context.Context, bannerID, slotID int) error
	IncrementBannerView(ctx context.Context, bannerID, slotID, socialGroupID int) error
	GetRandomBanner(ctx context.Context, slotID int) (int, error)
	CreateBanner(ctx context.Context, desc string) (int, error)
	CreateSlot(ctx context.Context, desc string) (int, error)
	CreateSocialGroup(ctx context.Context, desc string) (int, error)
	DeleteBanner(ctx context.Context, id int) error
	DeleteSlot(ctx context.Context, id int) error
	DeleteSocialGroup(ctx context.Context, id int) error
}
