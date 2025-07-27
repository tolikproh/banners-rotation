package internalhttp

type BannerSlot struct {
	BannerID int `json:"bannerId"`
	SlotID   int `json:"slotId"`
}

type BannerSlotSocialGroup struct {
	BannerID      int `json:"bannerId"`
	SlotID        int `json:"slotId"`
	SocialGroupID int `json:"socialGroupId"`
}

type SlotSocialGroup struct {
	SlotID        int `json:"slotId"`
	SocialGroupID int `json:"socialGroupId"`
}

type BannerResponse struct {
	BannerID int `json:"bannerId"`
}

type SlotResponse struct {
	SlotID int `json:"slotId"`
}

type SocialGroupResponse struct {
	SocialGroupID int `json:"socialGroupId"`
}

type DescRequest struct {
	Desc string `json:"desc"`
}

type errorHTTP struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Error  string `json:"error"`
}
