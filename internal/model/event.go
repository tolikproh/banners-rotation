package model

import (
	"fmt"
	"strings"
	"time"
)

type Event struct {
	EventType     string
	SlotID        int
	BannerID      int
	SocialGroupID int
	DateTime      time.Time
}

func (s Event) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("New event for statistic. event: %s, slot_id: %d, banner_id: %d, "+
		"social_group_id: %d at %s",
		s.EventType, s.SlotID, s.BannerID, s.SocialGroupID, s.DateTime.Format(time.RFC3339)))
	return builder.String()
}
