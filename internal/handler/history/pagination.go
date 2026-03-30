package history

import (
	"fmt"
	"strings"
	"time"
)

type CursorParams struct {
	UseCursor  bool
	BeforeTime time.Time
	BeforeID   string
	Offset     int32
}

func defaultLimit(limit int32) int32 {
	if limit == 0 {
		return 10
	}

	return limit
}

func parseCursorParams(beforeFirstSeen string, beforeID string, offset int32) (*CursorParams, error) {
	beforeFirstSeen = strings.TrimSpace(beforeFirstSeen)
	beforeID = strings.TrimSpace(beforeID)

	if beforeFirstSeen == "" && beforeID == "" {
		return &CursorParams{Offset: offset}, nil
	}

	if beforeFirstSeen == "" || beforeID == "" {
		return nil, fmt.Errorf("before_first_seen and before_id must be provided together")
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, beforeFirstSeen)
	if err != nil {
		return nil, fmt.Errorf("before_first_seen must be a valid RFC3339 timestamp")
	}

	return &CursorParams{
		UseCursor:  true,
		BeforeTime: parsedTime,
		BeforeID:   beforeID,
		Offset:     0,
	}, nil
}
