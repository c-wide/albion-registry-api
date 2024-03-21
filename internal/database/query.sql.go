// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const findPlayersByNameAndRegion = `-- name: FindPlayersByNameAndRegion :many
SELECT player_id AS id, name
FROM players
WHERE name ILIKE $1 || '%' AND region = $2
`

type FindPlayersByNameAndRegionParams struct {
	Name   pgtype.Text `json:"name"`
	Region RegionEnum  `json:"region"`
}

type FindPlayersByNameAndRegionRow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) FindPlayersByNameAndRegion(ctx context.Context, arg FindPlayersByNameAndRegionParams) ([]FindPlayersByNameAndRegionRow, error) {
	rows, err := q.db.Query(ctx, findPlayersByNameAndRegion, arg.Name, arg.Region)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FindPlayersByNameAndRegionRow{}
	for rows.Next() {
		var i FindPlayersByNameAndRegionRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCountsOfEntities = `-- name: GetCountsOfEntities :one
SELECT
  (SELECT COUNT(*) FROM players) AS total_players,
  (SELECT COUNT(*) FROM guilds) AS total_guilds,
  (SELECT COUNT(*) FROM alliances) AS total_alliances
`

type GetCountsOfEntitiesRow struct {
	TotalPlayers   int64 `json:"total_players"`
	TotalGuilds    int64 `json:"total_guilds"`
	TotalAlliances int64 `json:"total_alliances"`
}

func (q *Queries) GetCountsOfEntities(ctx context.Context) (GetCountsOfEntitiesRow, error) {
	row := q.db.QueryRow(ctx, getCountsOfEntities)
	var i GetCountsOfEntitiesRow
	err := row.Scan(&i.TotalPlayers, &i.TotalGuilds, &i.TotalAlliances)
	return i, err
}

const getPlayerAllianceMemberships = `-- name: GetPlayerAllianceMemberships :many
SELECT a.alliance_id AS id, a.tag, pam.first_seen, pam.last_seen
FROM player_alliance_memberships pam
JOIN alliances a ON pam.alliance_id = a.alliance_id AND pam.region = a.region
WHERE pam.player_id = $1 AND pam.region = $2
`

type GetPlayerAllianceMembershipsParams struct {
	PlayerID string     `json:"player_id"`
	Region   RegionEnum `json:"region"`
}

type GetPlayerAllianceMembershipsRow struct {
	ID        string    `json:"id"`
	Tag       string    `json:"tag"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

func (q *Queries) GetPlayerAllianceMemberships(ctx context.Context, arg GetPlayerAllianceMembershipsParams) ([]GetPlayerAllianceMembershipsRow, error) {
	rows, err := q.db.Query(ctx, getPlayerAllianceMemberships, arg.PlayerID, arg.Region)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPlayerAllianceMembershipsRow{}
	for rows.Next() {
		var i GetPlayerAllianceMembershipsRow
		if err := rows.Scan(
			&i.ID,
			&i.Tag,
			&i.FirstSeen,
			&i.LastSeen,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPlayerByIdAndRegion = `-- name: GetPlayerByIdAndRegion :one
SELECT player_id AS id, name, first_seen, last_seen
FROM players
WHERE player_id = $1 AND region = $2
`

type GetPlayerByIdAndRegionParams struct {
	PlayerID string     `json:"player_id"`
	Region   RegionEnum `json:"region"`
}

type GetPlayerByIdAndRegionRow struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

func (q *Queries) GetPlayerByIdAndRegion(ctx context.Context, arg GetPlayerByIdAndRegionParams) (GetPlayerByIdAndRegionRow, error) {
	row := q.db.QueryRow(ctx, getPlayerByIdAndRegion, arg.PlayerID, arg.Region)
	var i GetPlayerByIdAndRegionRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.FirstSeen,
		&i.LastSeen,
	)
	return i, err
}

const getPlayerGuildMemberships = `-- name: GetPlayerGuildMemberships :many
SELECT g.guild_id AS id, g.name, pgm.first_seen, pgm.last_seen
FROM player_guild_memberships pgm
JOIN guilds g ON pgm.guild_id = g.guild_id AND pgm.region = g.region
WHERE pgm.player_id = $1 AND pgm.region = $2
`

type GetPlayerGuildMembershipsParams struct {
	PlayerID string     `json:"player_id"`
	Region   RegionEnum `json:"region"`
}

type GetPlayerGuildMembershipsRow struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

func (q *Queries) GetPlayerGuildMemberships(ctx context.Context, arg GetPlayerGuildMembershipsParams) ([]GetPlayerGuildMembershipsRow, error) {
	rows, err := q.db.Query(ctx, getPlayerGuildMemberships, arg.PlayerID, arg.Region)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPlayerGuildMembershipsRow{}
	for rows.Next() {
		var i GetPlayerGuildMembershipsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.FirstSeen,
			&i.LastSeen,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
