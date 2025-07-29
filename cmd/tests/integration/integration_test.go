package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/tolikproh/banners-rotation/internal/cnst"
	modelHTTP "github.com/tolikproh/banners-rotation/internal/server/http"
)

var apiURL = "http://banners:8000"

// var apiURL = "http://localhost:8000"

type testHTTP struct {
	suite.Suite
	ctx           context.Context
	client        *http.Client
	banners       []int
	slotID        int
	socialGroupID int
}

func (s *testHTTP) SetupTest() {
	s.ctx = context.Background()
	s.client = &http.Client{Timeout: 30 * time.Second}
	s.initData()
}

func (s *testHTTP) TearDownTest() {
	for _, id := range s.banners {
		s.deleteBanner(id)
	}
	s.deleteSlot(s.slotID)
	s.deleteSocialGroup(s.socialGroupID)
}

func (s *testHTTP) initData() {
	s.banners = nil
	var err error

	for i := 0; i < 3; i++ {
		desc := fmt.Sprintf("new banner for tests: %d", i)
		id, err := s.createBanner(desc)
		if err != nil {
			panic(err)
		}
		s.banners = append(s.banners, id)
	}

	s.slotID, err = s.createSlot("new slot for tests")
	if err != nil {
		panic(err)
	}

	s.socialGroupID, err = s.createSocialGroup("new social group for tests")
	if err != nil {
		panic(err)
	}

	for _, bannerID := range s.banners {
		err := s.addBannerToSlot(bannerID, s.slotID)
		if err != nil {
			panic(err)
		}
	}
}

func (s *testHTTP) getBanner(slotID, socialGroupID int) (*modelHTTP.BannerResponse, error) {
	url := fmt.Sprintf("%s/banner?slotId=%d&socialGroupId=%d", apiURL, slotID, socialGroupID)

	req, err := http.NewRequestWithContext(s.ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, cnst.ErrStatusCode
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bres := new(modelHTTP.BannerResponse)
	err = json.Unmarshal(b, bres)
	if err != nil {
		return nil, err
	}

	return bres, nil
}

func (s *testHTTP) addBannerToSlot(bannerID, slotID int) error {
	bS := &modelHTTP.BannerSlot{
		BannerID: bannerID,
		SlotID:   slotID,
	}

	data, err := json.Marshal(bS)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/banner", apiURL)
	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return cnst.ErrStatusCode
	}

	return nil
}

func (s *testHTTP) clickBanner(bannerID, slotID, socialGroupID int) error {
	bSSG := &modelHTTP.BannerSlotSocialGroup{
		BannerID:      bannerID,
		SlotID:        slotID,
		SocialGroupID: socialGroupID,
	}

	data, err := json.Marshal(bSSG)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/click", apiURL)
	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return cnst.ErrStatusCode
	}

	return nil
}

func (s *testHTTP) deleteBannerToSlot(bannerID, slotID int) error {
	url := fmt.Sprintf("%s/banner?bannerId=%d&slotId=%d", apiURL, bannerID, slotID)

	req, err := http.NewRequestWithContext(s.ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return cnst.ErrStatusCode
	}

	return nil
}

func (s *testHTTP) createBanner(description string) (int, error) {
	desc := &modelHTTP.DescRequest{
		Desc: description,
	}

	data, err := json.Marshal(desc)
	if err != nil {
		return 0, err
	}

	url := fmt.Sprintf("%s/create/banner", apiURL)
	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, cnst.ErrStatusCode
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	id := new(modelHTTP.BannerResponse)
	err = json.Unmarshal(b, id)
	if err != nil {
		return 0, err
	}

	return id.BannerID, nil
}

func (s *testHTTP) createSlot(description string) (int, error) {
	desc := &modelHTTP.DescRequest{
		Desc: description,
	}

	data, err := json.Marshal(desc)
	if err != nil {
		return 0, err
	}

	url := fmt.Sprintf("%s/create/slot", apiURL)
	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, cnst.ErrStatusCode
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	id := new(modelHTTP.SlotResponse)
	err = json.Unmarshal(b, id)
	if err != nil {
		return 0, err
	}

	return id.SlotID, nil
}

func (s *testHTTP) createSocialGroup(description string) (int, error) {
	desc := &modelHTTP.DescRequest{
		Desc: description,
	}

	data, err := json.Marshal(desc)
	if err != nil {
		return 0, err
	}

	url := fmt.Sprintf("%s/create/group", apiURL)
	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, cnst.ErrStatusCode
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	id := new(modelHTTP.SocialGroupResponse)
	err = json.Unmarshal(b, id)
	if err != nil {
		return 0, err
	}

	return id.SocialGroupID, nil
}

func (s *testHTTP) deleteBanner(id int) error {
	url := fmt.Sprintf("%s/delete/banner?id=%d", apiURL, id)

	req, err := http.NewRequestWithContext(s.ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return cnst.ErrStatusCode
	}

	return nil
}

func (s *testHTTP) deleteSlot(id int) error {
	url := fmt.Sprintf("%s/delete/slot?id=%d", apiURL, id)

	req, err := http.NewRequestWithContext(s.ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return cnst.ErrStatusCode
	}

	return nil
}

func (s *testHTTP) deleteSocialGroup(id int) error {
	url := fmt.Sprintf("%s/delete/group?id=%d", apiURL, id)

	req, err := http.NewRequestWithContext(s.ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return cnst.ErrStatusCode
	}

	return nil
}

// Testing

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(testHTTP))
}

func (s *testHTTP) TestShowAllBanners() {
	var banners []int

	for i := 0; i < 20; i++ {
		result, err := s.getBanner(s.slotID, s.socialGroupID)
		s.Require().NoError(err)
		banners = append(banners, result.BannerID)
	}

	for _, expectedBannerID := range s.banners {
		s.Contains(banners, expectedBannerID)
	}
}

func (s *testHTTP) TestMoreShowsForPopularBanner() {
	for i := 0; i < 20; i++ {
		_, err := s.getBanner(s.slotID, s.socialGroupID)
		s.Require().NoError(err)
	}

	for i := 0; i < 10; i++ {
		err := s.clickBanner(s.banners[0], s.slotID, s.socialGroupID)
		s.Require().NoError(err)
	}

	var banners []int
	for i := 0; i < 20; i++ {
		result, err := s.getBanner(s.slotID, s.socialGroupID)
		s.Require().NoError(err)
		banners = append(banners, result.BannerID)
	}

	freq := make(map[int]int)
	for _, bannerID := range banners {
		freq[bannerID]++
	}

	var (
		maxFreqBannerID = 0
		maxFreqValue    = -1
	)

	for bannerID, freqValue := range freq {
		if freqValue > maxFreqValue {
			maxFreqValue = freqValue
			maxFreqBannerID = bannerID
		}
	}

	s.Require().Equal(s.banners[0], maxFreqBannerID)
}

func (s *testHTTP) TestDeleteBannerToSlot() {
	for i := 0; i < 20; i++ {
		_, err := s.getBanner(s.slotID, s.socialGroupID)
		s.Require().NoError(err)
	}

	for _, bannerID := range s.banners {
		err := s.deleteBannerToSlot(bannerID, s.slotID)
		s.Require().NoError(err)
	}

	for i := 0; i < 20; i++ {
		_, err := s.getBanner(s.slotID, s.socialGroupID)
		s.Require().Error(err)
	}
}
