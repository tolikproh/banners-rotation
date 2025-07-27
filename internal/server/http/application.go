package internalhttp

import (
	"context"
)

type Application interface {
	GetBanner(ctx context.Context, slotID, socialGroupID int) (int, error)
	AddBannerToSlot(ctx context.Context, bannerID, slotID int) error
	ClickBanner(ctx context.Context, bannerID, slotID, socialGroupID int) error
	DeleteBannerFromSlot(ctx context.Context, bannerID, slotID int) error
	CreateBanner(ctx context.Context, desc string) (int, error)
	CreateSlot(ctx context.Context, desc string) (int, error)
	CreateSocialGroup(ctx context.Context, desc string) (int, error)
	DeleteBanner(ctx context.Context, id int) error
	DeleteSlot(ctx context.Context, id int) error
	DeleteSocialGroup(ctx context.Context, id int) error
}
