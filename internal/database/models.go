// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type RegionEnum string

const (
	RegionEnumAmericas RegionEnum = "americas"
	RegionEnumEurope   RegionEnum = "europe"
	RegionEnumAsia     RegionEnum = "asia"
)

func (e *RegionEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RegionEnum(s)
	case string:
		*e = RegionEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for RegionEnum: %T", src)
	}
	return nil
}

type NullRegionEnum struct {
	RegionEnum RegionEnum `json:"region_enum"`
	Valid      bool       `json:"valid"` // Valid is true if RegionEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRegionEnum) Scan(value interface{}) error {
	if value == nil {
		ns.RegionEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RegionEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRegionEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RegionEnum), nil
}

type Alliance struct {
	Tag        string     `json:"tag"`
	AllianceID string     `json:"alliance_id"`
	Region     RegionEnum `json:"region"`
	FirstSeen  time.Time  `json:"first_seen"`
	LastSeen   time.Time  `json:"last_seen"`
}

type Guild struct {
	Name      string     `json:"name"`
	GuildID   string     `json:"guild_id"`
	Region    RegionEnum `json:"region"`
	FirstSeen time.Time  `json:"first_seen"`
	LastSeen  time.Time  `json:"last_seen"`
}

type GuildAllianceMembership struct {
	ID         pgtype.UUID `json:"id"`
	GuildID    string      `json:"guild_id"`
	AllianceID string      `json:"alliance_id"`
	Region     RegionEnum  `json:"region"`
	FirstSeen  time.Time   `json:"first_seen"`
	LastSeen   time.Time   `json:"last_seen"`
}

type Player struct {
	Name      string     `json:"name"`
	PlayerID  string     `json:"player_id"`
	Region    RegionEnum `json:"region"`
	FirstSeen time.Time  `json:"first_seen"`
	LastSeen  time.Time  `json:"last_seen"`
}

type PlayerAllianceMembership struct {
	ID         pgtype.UUID `json:"id"`
	PlayerID   string      `json:"player_id"`
	AllianceID string      `json:"alliance_id"`
	Region     RegionEnum  `json:"region"`
	FirstSeen  time.Time   `json:"first_seen"`
	LastSeen   time.Time   `json:"last_seen"`
}

type PlayerGuildMembership struct {
	ID        pgtype.UUID `json:"id"`
	PlayerID  string      `json:"player_id"`
	GuildID   string      `json:"guild_id"`
	Region    RegionEnum  `json:"region"`
	FirstSeen time.Time   `json:"first_seen"`
	LastSeen  time.Time   `json:"last_seen"`
}
