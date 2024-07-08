-- name: GetCountsOfEntities :one
SELECT
  (SELECT COUNT(*) FROM players) AS players,
  (SELECT COUNT(*) FROM guilds) AS guilds,
  (SELECT COUNT(*) FROM alliances) AS alliances;

-- name: GetPlayer :one
SELECT
	name,
	player_id AS id,
	first_seen,
	last_seen
FROM
	players p
WHERE 
	p.player_id = $1
	AND p.region = $2;

-- name: GetGuild :one
SELECT
	name,
	guild_id AS id,
	first_seen,
	last_seen
FROM
	guilds g
WHERE 
	g.guild_id = $1
	AND g.region = $2;

-- name: GetAlliance :one
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
	AND a.region = $2;

-- name: SearchEntities :many
(
    SELECT 
        'player' AS type,
        player_id AS id,
        name,
        '' AS tag,
        CASE
            WHEN LOWER(name) = LOWER(@searchTerm::text) THEN 0
            ELSE 1
        END AS match_rank
    FROM 
        players p
    WHERE 
        p.region = $1 AND p.name ILIKE (@searchTerm::text || '%')
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
            WHEN LOWER(name) = LOWER(@searchTerm::text) THEN 0
            ELSE 1
        END AS match_rank
    FROM 
        guilds g
    WHERE 
        g.region = $1 AND g.name ILIKE (@searchTerm::text || '%')
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
            WHEN LOWER(name) = LOWER(@searchTerm::text) THEN 0
            WHEN LOWER(tag) = LOWER(@searchTerm::text) THEN 0
            ELSE 1
        END AS match_rank
    FROM 
        alliances a
    WHERE 
        a.region = $1 AND (a.name ILIKE (@searchTerm::text || '%') OR a.tag ILIKE (@searchTerm::text || '%'))
    ORDER BY
        match_rank,
        LENGTH(name)
    LIMIT $2
);

-- name: GetPlayerHistory :many
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
        LIMIT @allianceLimit
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
LIMIT $3 OFFSET $4;

-- name: GetPlayerGuildAlliances :many
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
LIMIT $4 OFFSET $5;

-- name: GetGuildAllianceHistory :many
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
LIMIT $3 OFFSET $4;

-- name: GetGuildPlayerHistory :many
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
LIMIT $3 OFFSET $4;

-- name: GetAllianceGuildHistory :many
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
LIMIT $3 OFFSET $4;
