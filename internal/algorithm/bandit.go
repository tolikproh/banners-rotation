package algorithm

import (
	"math"
	"math/rand"

	"github.com/tolikproh/banners-rotation/internal/cnst"
	"github.com/tolikproh/banners-rotation/internal/model"
)

func Bandit(banners []model.Banner) (int, error) {
	var bannerID int
	var bannerIds []int
	var totalViewCount int64
	var maxIncome float64 = -1
	for _, banner := range banners {
		totalViewCount += banner.ViewCount
	}

	for _, banner := range banners {
		bannerIncome := (float64(banner.ClickCount) / float64(banner.ViewCount)) +
			math.Sqrt((2.0*math.Log(float64(totalViewCount)))/float64(banner.ViewCount))
		if bannerIncome > maxIncome {
			maxIncome = bannerIncome
			bannerIds = bannerIds[:0]
			bannerIds = append(bannerIds, int(banner.BannerID))
		} else if bannerIncome == maxIncome {
			bannerIds = append(bannerIds, int(banner.BannerID))
		}
	}

	if len(bannerIds) == 0 {
		return 0, cnst.ErrGetBanner
	}
	if len(bannerIds) == 1 {
		return bannerIds[0], nil
	}

	index := rand.Intn(len(bannerIds)) //nolint: gosec
	bannerID = bannerIds[index]

	return bannerID, nil
}
