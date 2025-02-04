// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package database

import (
	"context"
	"time"
)

const getAlliance = `-- name: GetAlliance :one
SELECT
	name,
    tag,
	alliance_id AS id,
	first_seen,
	last_seen
FROM
	alliances a
WHERE 
	a.alliance_id = $1
	AND a.region = $2
`

type GetAllianceParams struct {
	AllianceID string `json:"alliance_id"`
	Region     string `json:"region"`
}

type GetAllianceRow struct {
	Name      *string   `json:"name"`
	Tag       string    `json:"tag"`
	ID        string    `json:"id"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

func (q *Queries) GetAlliance(ctx context.Context, arg GetAllianceParams) (GetAllianceRow, error) {
	row := q.db.QueryRow(ctx, getAlliance, arg.AllianceID, arg.Region)
	var i GetAllianceRow
	err := row.Scan(
		&i.Name,
		&i.Tag,
		&i.ID,
		&i.FirstSeen,
		&i.LastSeen,
	)
	return i, err
}

const getAllianceGuildHistory = `-- name: GetAllianceGuildHistory :many
SELECT
    g.guild_id,
	g.name,
	gam.first_seen,
	gam.last_seen
FROM
	guild_alliance_memberships gam
JOIN
	guilds g ON
	gam.guild_id = g.guild_id
	AND gam.region = g.region
WHERE
	gam.alliance_id = $1
	AND gam.region = $2
ORDER BY
	gam.first_seen DESC
LIMIT $3 OFFSET $4
`

type GetAllianceGuildHistoryParams struct {
	AllianceID string `json:"alliance_id"`
	Region     string `json:"region"`
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
}

type GetAllianceGuildHistoryRow struct {
	GuildID   string    `json:"guild_id"`
	Name      string    `json:"name"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

func (q *Queries) GetAllianceGuildHistory(ctx context.Context, arg GetAllianceGuildHistoryParams) ([]GetAllianceGuildHistoryRow, error) {
	rows, err := q.db.Query(ctx, getAllianceGuildHistory,
		arg.AllianceID,
		arg.Region,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllianceGuildHistoryRow{}
	for rows.Next() {
		var i GetAllianceGuildHistoryRow
		if err := rows.Scan(
			&i.GuildID,
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

const getCountsOfEntities = `-- name: GetCountsOfEntities :one
SELECT
  (SELECT COUNT(*) FROM players) AS players,
  (SELECT COUNT(*) FROM guilds) AS guilds,
  (SELECT COUNT(*) FROM alliances) AS alliances
`

type GetCountsOfEntitiesRow struct {
	Players   int64 `json:"players"`
	Guilds    int64 `json:"guilds"`
	Alliances int64 `json:"alliances"`
}

func (q *Queries) GetCountsOfEntities(ctx context.Context) (GetCountsOfEntitiesRow, error) {
	row := q.db.QueryRow(ctx, getCountsOfEntities)
	var i GetCountsOfEntitiesRow
	err := row.Scan(&i.Players, &i.Guilds, &i.Alliances)
	return i, err
}

const getGuild = `-- name: GetGuild :one
SELECT
	name,
	guild_id AS id,
	first_seen,
	last_seen
FROM
	guilds g
WHERE 
	g.guild_id = $1
	AND g.region = $2
`

type GetGuildParams struct {
	GuildID string `json:"guild_id"`
	Region  string `json:"region"`
}

type GetGuildRow struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

func (q *Queries) GetGuild(ctx context.Context, arg GetGuildParams) (GetGuildRow, error) {
	row := q.db.QueryRow(ctx, getGuild, arg.GuildID, arg.Region)
	var i GetGuildRow
	err := row.Scan(
		&i.Name,
		&i.ID,
		&i.FirstSeen,
		&i.LastSeen,
	)
	return i, err
}

const getGuildAllianceHistory = `-- name: GetGuildAllianceHistory :many
SELECT
    a.alliance_id,
	a.name,
	a.tag,
	gam.first_seen,
	gam.last_seen
FROM
	guild_alliance_memberships gam
JOIN
	alliances a ON
	gam.alliance_id = a.alliance_id
	AND gam.region = a.region
WHERE
	gam.guild_id = $1
	AND gam.region = $2
ORDER BY
	gam.first_seen DESC
LIMIT $3 OFFSET $4
`

type GetGuildAllianceHistoryParams struct {
	GuildID string `json:"guild_id"`
	Region  string `json:"region"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

type GetGuildAllianceHistoryRow struct {
	AllianceID string    `json:"alliance_id"`
	Name       *string   `json:"name"`
	Tag        string    `json:"tag"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
}

func (q *Queries) GetGuildAllianceHistory(ctx context.Context, arg GetGuildAllianceHistoryParams) ([]GetGuildAllianceHistoryRow, error) {
	rows, err := q.db.Query(ctx, getGuildAllianceHistory,
		arg.GuildID,
		arg.Region,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetGuildAllianceHistoryRow{}
	for rows.Next() {
		var i GetGuildAllianceHistoryRow
		if err := rows.Scan(
			&i.AllianceID,
			&i.Name,
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

const getGuildPlayerHistory = `-- name: GetGuildPlayerHistory :many
SELECT
	p.player_id,
	p.name,
	pgm.first_seen,
	pgm.last_seen
FROM
	player_guild_memberships pgm
JOIN players p ON
	pgm.player_id = p.player_id
	AND pgm.region = p.region
WHERE 
	pgm.guild_id = $1
	AND pgm.region = $2
ORDER BY
	pgm.first_seen DESC
LIMIT $3 OFFSET $4
`

type GetGuildPlayerHistoryParams struct {
	GuildID string `json:"guild_id"`
	Region  string `json:"region"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

type GetGuildPlayerHistoryRow struct {
	PlayerID  string    `json:"player_id"`
	Name      string    `json:"name"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

func (q *Queries) GetGuildPlayerHistory(ctx context.Context, arg GetGuildPlayerHistoryParams) ([]GetGuildPlayerHistoryRow, error) {
	rows, err := q.db.Query(ctx, getGuildPlayerHistory,
		arg.GuildID,
		arg.Region,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetGuildPlayerHistoryRow{}
	for rows.Next() {
		var i GetGuildPlayerHistoryRow
		if err := rows.Scan(
			&i.PlayerID,
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

const getPlayer = `-- name: GetPlayer :one
SELECT
	name,
	player_id AS id,
    avatar,
    avatar_ring,
	first_seen,
	last_seen
FROM
	players p
WHERE 
	p.player_id = $1
	AND p.region = $2
`

type GetPlayerParams struct {
	PlayerID string `json:"player_id"`
	Region   string `json:"region"`
}

type GetPlayerRow struct {
	Name       string    `json:"name"`
	ID         string    `json:"id"`
	Avatar     *string   `json:"avatar"`
	AvatarRing *string   `json:"avatar_ring"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
}

func (q *Queries) GetPlayer(ctx context.Context, arg GetPlayerParams) (GetPlayerRow, error) {
	row := q.db.QueryRow(ctx, getPlayer, arg.PlayerID, arg.Region)
	var i GetPlayerRow
	err := row.Scan(
		&i.Name,
		&i.ID,
		&i.Avatar,
		&i.AvatarRing,
		&i.FirstSeen,
		&i.LastSeen,
	)
	return i, err
}

const getPlayerGuildAlliances = `-- name: GetPlayerGuildAlliances :many
SELECT
    a.alliance_id,
    a.name,
    a.tag,
    GREATEST(pgm.first_seen, gam.first_seen) AS first_seen,
    LEAST(pgm.last_seen, gam.last_seen) AS last_seen
FROM
    player_guild_memberships pgm
JOIN
    guild_alliance_memberships gam ON
    pgm.guild_id = gam.guild_id AND pgm.region = gam.region
JOIN 
    alliances a ON
    gam.alliance_id = a.alliance_id AND gam.region = a.region
WHERE
    pgm.player_id = $1
    AND pgm.region = $2
    AND pgm.guild_id = $3
    AND gam.first_seen <= pgm.last_seen
    AND gam.last_seen >= pgm.first_seen
ORDER BY
    gam.first_seen DESC
LIMIT $4 OFFSET $5
`

type GetPlayerGuildAlliancesParams struct {
	PlayerID string `json:"player_id"`
	Region   string `json:"region"`
	GuildID  string `json:"guild_id"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

type GetPlayerGuildAlliancesRow struct {
	AllianceID string      `json:"alliance_id"`
	Name       *string     `json:"name"`
	Tag        string      `json:"tag"`
	FirstSeen  interface{} `json:"first_seen"`
	LastSeen   interface{} `json:"last_seen"`
}

func (q *Queries) GetPlayerGuildAlliances(ctx context.Context, arg GetPlayerGuildAlliancesParams) ([]GetPlayerGuildAlliancesRow, error) {
	rows, err := q.db.Query(ctx, getPlayerGuildAlliances,
		arg.PlayerID,
		arg.Region,
		arg.GuildID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPlayerGuildAlliancesRow{}
	for rows.Next() {
		var i GetPlayerGuildAlliancesRow
		if err := rows.Scan(
			&i.AllianceID,
			&i.Name,
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

const getPlayerHistory = `-- name: GetPlayerHistory :many
SELECT
	g.name,
	g.guild_id,
	pgm.first_seen,
	pgm.last_seen,
	COALESCE(
	(
	SELECT
		array_to_json(array_agg(row_to_json(a)))
	FROM
		(
		SELECT
			a.alliance_id,
			a.name,
			a.tag,
			GREATEST(pgm.first_seen, gam.first_seen) AS first_seen,
			LEAST(pgm.last_seen, gam.last_seen) AS last_seen
		FROM
			guild_alliance_memberships gam
		JOIN 
            alliances a ON
			gam.alliance_id = a.alliance_id
			AND gam.region = a.region
		WHERE
			gam.guild_id = g.guild_id
			AND gam.region = g.region
			AND gam.first_seen <= pgm.last_seen
			AND gam.last_seen >= pgm.first_seen
		ORDER BY
			gam.first_seen DESC
        LIMIT $5
		) a
	),
	'[]'::JSON
	) AS alliances
FROM
	player_guild_memberships pgm
JOIN
	guilds g ON
	pgm.guild_id = g.guild_id
	AND pgm.region = g.region
WHERE 
	pgm.player_id = $1
	AND pgm.region = $2
ORDER BY 
	pgm.first_seen DESC
LIMIT $3 OFFSET $4
`

type GetPlayerHistoryParams struct {
	PlayerID      string `json:"player_id"`
	Region        string `json:"region"`
	Limit         int32  `json:"limit"`
	Offset        int32  `json:"offset"`
	Alliancelimit int32  `json:"alliancelimit"`
}

type GetPlayerHistoryRow struct {
	Name      string      `json:"name"`
	GuildID   string      `json:"guild_id"`
	FirstSeen time.Time   `json:"first_seen"`
	LastSeen  time.Time   `json:"last_seen"`
	Alliances interface{} `json:"alliances"`
}

func (q *Queries) GetPlayerHistory(ctx context.Context, arg GetPlayerHistoryParams) ([]GetPlayerHistoryRow, error) {
	rows, err := q.db.Query(ctx, getPlayerHistory,
		arg.PlayerID,
		arg.Region,
		arg.Limit,
		arg.Offset,
		arg.Alliancelimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPlayerHistoryRow{}
	for rows.Next() {
		var i GetPlayerHistoryRow
		if err := rows.Scan(
			&i.Name,
			&i.GuildID,
			&i.FirstSeen,
			&i.LastSeen,
			&i.Alliances,
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

const searchEntities = `-- name: SearchEntities :many
(
    SELECT 
        'player' AS type,
        player_id AS id,
        name,
        '' AS tag,
        CASE
            WHEN LOWER(name) = LOWER($3::text) THEN 0
            ELSE 1
        END AS match_rank
    FROM 
        players p
    WHERE 
        p.region = $1 AND p.name ILIKE ($3::text || '%')
    ORDER BY
        match_rank,
        LENGTH(name)
    LIMIT $2
)
UNION ALL
(
    SELECT 
        'guild' AS type,
        guild_id AS id,
        name,
        '' AS tag,
        CASE
            WHEN LOWER(name) = LOWER($3::text) THEN 0
            ELSE 1
        END AS match_rank
    FROM 
        guilds g
    WHERE 
        g.region = $1 AND g.name ILIKE ($3::text || '%')
    ORDER BY
        match_rank,
        LENGTH(name)
    LIMIT $2
)
UNION ALL
(
    SELECT 
        'alliance' AS type,
        alliance_id AS id,
        name,
        tag,
        CASE
            WHEN LOWER(name) = LOWER($3::text) THEN 0
            WHEN LOWER(tag) = LOWER($3::text) THEN 0
            ELSE 1
        END AS match_rank
    FROM 
        alliances a
    WHERE 
        a.region = $1 AND (a.name ILIKE ($3::text || '%') OR a.tag ILIKE ($3::text || '%'))
    ORDER BY
        match_rank,
        LENGTH(name)
    LIMIT $2
)
`

type SearchEntitiesParams struct {
	Region     string `json:"region"`
	Limit      int32  `json:"limit"`
	Searchterm string `json:"searchterm"`
}

type SearchEntitiesRow struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	MatchRank int32  `json:"match_rank"`
}

func (q *Queries) SearchEntities(ctx context.Context, arg SearchEntitiesParams) ([]SearchEntitiesRow, error) {
	rows, err := q.db.Query(ctx, searchEntities, arg.Region, arg.Limit, arg.Searchterm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SearchEntitiesRow{}
	for rows.Next() {
		var i SearchEntitiesRow
		if err := rows.Scan(
			&i.Type,
			&i.ID,
			&i.Name,
			&i.Tag,
			&i.MatchRank,
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
