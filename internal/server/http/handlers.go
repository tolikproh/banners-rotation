package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/tolikproh/banners-rotation/pkg/httperr"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API rotation-banners")
}

func (s *Server) getBanner(ctx context.Context, r *http.Request) (interface{}, error) {
	slotID, err := strconv.Atoi(r.URL.Query().Get("slotId"))
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	groupID, err := strconv.Atoi(r.URL.Query().Get("socialGroupId"))
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	sSG := SlotSocialGroup{
		SlotID:        slotID,
		SocialGroupID: groupID,
	}

	idBanner, err := s.app.GetBanner(ctx, sSG.SlotID, sSG.SocialGroupID)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	result := BannerResponse{BannerID: idBanner}
	return result, nil
}

func (s *Server) addBannerToSlot(ctx context.Context, r *http.Request) (interface{}, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	var bS BannerSlot
	err = json.Unmarshal(body, &bS)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := s.app.AddBannerToSlot(ctx, bS.BannerID, bS.SlotID); err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil, nil
}

func (s *Server) clickBanner(ctx context.Context, r *http.Request) (interface{}, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	var bSS BannerSlotSocialGroup
	err = json.Unmarshal(body, &bSS)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err = s.app.ClickBanner(ctx, bSS.BannerID, bSS.SlotID, bSS.SocialGroupID); err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil, nil
}

func (s *Server) deleteBannerFromSlot(ctx context.Context, r *http.Request) (interface{}, error) {
	bannerID, err := strconv.Atoi(r.URL.Query().Get("bannerId"))
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	slotID, err := strconv.Atoi(r.URL.Query().Get("slotId"))
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	bS := BannerSlot{
		BannerID: bannerID,
		SlotID:   slotID,
	}

	if err := s.app.DeleteBannerFromSlot(ctx, bS.BannerID, bS.SlotID); err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil, nil
}

func (s *Server) createBanner(ctx context.Context, r *http.Request) (interface{}, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	var req DescRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	id, err := s.app.CreateBanner(ctx, req.Desc)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return &BannerResponse{BannerID: id}, nil
}

func (s *Server) createSlot(ctx context.Context, r *http.Request) (interface{}, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	var req DescRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	id, err := s.app.CreateSlot(ctx, req.Desc)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return &SlotResponse{SlotID: id}, nil
}

func (s *Server) createSocislGroup(ctx context.Context, r *http.Request) (interface{}, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	var req DescRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	id, err := s.app.CreateSocialGroup(ctx, req.Desc)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return &SocialGroupResponse{SocialGroupID: id}, nil
}

func (s *Server) deleteBanner(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := s.app.DeleteBanner(ctx, id); err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil, nil
}

func (s *Server) deleteSlot(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := s.app.DeleteSlot(ctx, id); err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil, nil
}

func (s *Server) deleteSocialGroup(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	if err := s.app.DeleteSocialGroup(ctx, id); err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil, nil
}
