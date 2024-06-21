-- name: GetCountsOfEntities :one
SELECT
  (SELECT COUNT(*) FROM players) AS total_players,
  (SELECT COUNT(*) FROM guilds) AS total_guilds,
  (SELECT COUNT(*) FROM alliances) AS total_alliances;

-- name: DoesPlayerExist :one
SELECT EXISTS(SELECT 1 FROM players WHERE player_id = $1 AND region = $2);

-- name: DoesGuildExist :one
SELECT EXISTS(SELECT 1 FROM guilds WHERE guild_id = $1 AND region = $2);

-- name: DoesAllianceExist :one
SELECT EXISTS(SELECT 1 FROM alliances WHERE alliance_id = $1 AND region = $2);

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
	pgm.first_seen DESC;

-- name: GetAllianceHistory :many
SELECT
	g.name,
	g.guild_id,
	gam.first_seen,
	gam.last_seen,
	COALESCE(
    (
	SELECT
		array_to_json(array_agg(row_to_json(p)))
	FROM
		(
		SELECT
			p.player_id,
			p.name,
			GREATEST(gam.first_seen, pgm.first_seen) AS first_seen,
			LEAST(gam.last_seen, pgm.last_seen) AS last_seen
		FROM
			player_guild_memberships pgm
		JOIN 
            players p ON
			pgm.player_id = p.player_id
			AND pgm.region = p.region
		WHERE
			pgm.guild_id = g.guild_id
			AND pgm.region = g.region
			AND pgm.first_seen <= gam.last_seen
			AND pgm.last_seen >= gam.first_seen
		ORDER BY
			pgm.first_seen DESC
        ) p
    ),
	'[]'::JSON
    ) AS players
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
	gam.first_seen DESC;
