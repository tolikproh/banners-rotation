package storage

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/tolikproh/banners-rotation/internal/cnst"
	"github.com/tolikproh/banners-rotation/internal/config"
	"github.com/tolikproh/banners-rotation/internal/logger"
	"github.com/tolikproh/banners-rotation/internal/model"
)

type Storage struct {
	cfg *config.Config
	log *logger.Logger
	db  *pgx.Conn
}

func New(ctx context.Context, cfg *config.Config, log *logger.Logger) (*Storage, error) {
	log.Debug("init storage")

	store := new(Storage)

	db, err := pgx.Connect(ctx, cfg.Storage.Conn)
	if err != nil {
		return nil, err
	}

	store.cfg = cfg
	store.log = log
	store.db = db

	return store, nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.log.Debug("close storage")

	return s.db.Close(ctx)
}

func (s *Storage) GetBannersInfo(ctx context.Context, slotID, socialGroupID int) ([]model.Banner, error) {
	s.log.Debug("storage GetBannersInfo")

	args := []interface{}{
		slotID,
		socialGroupID,
	}

	query := `SELECT bs.banner_id, bs.slot_id, bv.social_group_id, 
				   count(distinct bv.id) view_count, count(distinct cl.id) click_count
			FROM banner_slot bs
			LEFT JOIN banner_views bv ON bv.slot_id = bs.slot_id AND bv.banner_id = bs.banner_id
			LEFT JOIN banner_clicks cl ON bv.slot_id = cl.slot_id AND bv.banner_id = cl.banner_id AND 
											   bv.social_group_id = cl.social_group_id
			WHERE bs.slot_id = $1 AND (bv.social_group_id = $2 OR bv.social_group_id is null)
			GROUP BY bs.banner_id, bs.slot_id, bv.social_group_id
			ORDER BY bv.social_group_id`

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, cnst.ErrGetRowsError
	}
	defer rows.Close()

	var banners []model.Banner
	for rows.Next() {
		var banner model.Banner

		err := rows.Scan(
			&banner.BannerID,
			&banner.SlotID,
			&banner.SocialGroupID,
			&banner.ViewCount,
			&banner.ClickCount)
		if err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}

	if len(banners) == 0 {
		return nil, cnst.ErrObjectNotFound
	}

	return banners, nil
}

func (s *Storage) AddBannerToSlot(ctx context.Context, bannerID, slotID int) error {
	s.log.Debug("storage AddBannerToSlot")

	args := []interface{}{
		bannerID,
		slotID,
	}

	query := `INSERT INTO banner_slot (banner_id, slot_id) VALUES ($1, $2)`
	_, err := s.db.Exec(ctx, query, args...)

	return err
}

func (s *Storage) ClickBanner(ctx context.Context, bannerID, slotID, socialGroupID int) error {
	s.log.Debug("storage ClickBanner")

	args := []interface{}{
		bannerID,
		slotID,
		socialGroupID,
	}

	query := `INSERT INTO banner_clicks (banner_id, slot_id, social_group_id, date) VALUES ($1, $2, $3, current_timestamp)`
	_, err := s.db.Exec(ctx, query, args...)

	return err
}

func (s *Storage) DeleteBannerFromSlot(ctx context.Context, bannerID, slotID int) error {
	s.log.Debug("storage DeleteBannerFromSlot")

	args := []interface{}{
		bannerID,
		slotID,
	}

	query := `DELETE FROM banner_slot WHERE banner_id=$1 AND slot_id=$2`
	_, err := s.db.Exec(ctx, query, args...)

	return err
}

func (s *Storage) IncrementBannerView(ctx context.Context, bannerID, slotID, socialGroupID int) error {
	s.log.Debug("storage IncrementBannerView")

	args := []interface{}{
		bannerID,
		slotID,
		socialGroupID,
	}

	query := `INSERT INTO banner_views (banner_id, slot_id, social_group_id, date) VALUES ($1, $2, $3, current_timestamp)`
	_, err := s.db.Exec(ctx, query, args...)

	return err
}

func (s *Storage) GetRandomBanner(ctx context.Context, slotID int) (int, error) {
	s.log.Debug("storage GetRandomBanner")

	query := `SELECT banner_id FROM banner_slot WHERE slot_id = $1 ORDER BY random() LIMIT 1`

	var bannerID int
	err := s.db.QueryRow(ctx, query, slotID).Scan(
		&bannerID,
	)
	if err != nil {
		return 0, cnst.ErrObjectNotFound
	}

	return bannerID, nil
}

func (s *Storage) CreateBanner(ctx context.Context, desc string) (int, error) {
	s.log.Debug("storage create banner")

	query := `INSERT INTO banners (description) VALUES ($1) RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query, desc).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) CreateSlot(ctx context.Context, desc string) (int, error) {
	s.log.Debug("storage create slot")

	query := `INSERT INTO slots (description) VALUES ($1) RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query, desc).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) CreateSocialGroup(ctx context.Context, desc string) (int, error) {
	s.log.Debug("storage create social group")

	query := `INSERT INTO social_groups (description) VALUES ($1) RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query, desc).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) DeleteBanner(ctx context.Context, id int) error {
	s.log.Debug("storage delete banner")

	query := `DELETE FROM banners WHERE id = $1`
	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteSlot(ctx context.Context, id int) error {
	s.log.Debug("storage delete slot")

	query := `DELETE FROM slots WHERE id = $1`
	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteSocialGroup(ctx context.Context, id int) error {
	s.log.Debug("storage delete social groups")

	query := `DELETE FROM social_groups WHERE id = $1`
	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
