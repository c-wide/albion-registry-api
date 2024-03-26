-- name: GetCountsOfEntities :one
SELECT
  (SELECT COUNT(*) FROM players) AS total_players,
  (SELECT COUNT(*) FROM guilds) AS total_guilds,
  (SELECT COUNT(*) FROM alliances) AS total_alliances;

-- name: SearchPlayers :many
SELECT player_id AS id, name
FROM players
WHERE name ILIKE @name || '%' AND region = @region
LIMIT $1 OFFSET $2;

-- name: SearchGuilds :many
SELECT guild_id AS id, name
FROM guilds
WHERE name ILIKE @name || '%' AND region = @region
LIMIT $1 OFFSET $2;

-- name: SearchAlliances :many
SELECT alliance_id AS id, COALESCE(name, ''), tag
FROM alliances
WHERE (name ILIKE @name || '%' OR tag ILIKE @name || '%') AND region = @region
LIMIT $1 OFFSET $2;

-- name: SearchEntities :many
SELECT * FROM (
    SELECT 'player' AS type, player_id AS id, name, NULL AS tag
    FROM players AS p
    WHERE p.name ILIKE @name || '%' AND p.region = @region
    LIMIT $1
) AS players_subquery
UNION ALL
SELECT * FROM (
    SELECT 'guild' AS type, guild_id AS id, name, NULL AS tag
    FROM guilds AS g
    WHERE g.name ILIKE @name || '%' AND g.region = @region
    LIMIT $1
) AS guilds_subquery
UNION ALL
SELECT * FROM (
    SELECT 'alliance' AS type, alliance_id AS id, COALESCE(name, ''), tag 
    FROM alliances AS a
    WHERE (a.name ILIKE @name || '%' OR a.tag ILIKE @name || '%') AND a.region = @region
    LIMIT $1
) AS alliances_subquery;

-- name: GetPlayerDetails :one
SELECT player_id AS id, name, first_seen, last_seen
FROM players
WHERE player_id = $1 AND region = $2;

-- name: GetGuildDetails :one
SELECT guild_id AS id, name, first_seen, last_seen
FROM guilds
WHERE guild_id = $1 AND region = $2;

-- name: GetAllianceDetails :one
SELECT alliance_id AS id, COALESCE(name, ''), tag, first_seen, last_seen
FROM alliances
WHERE alliance_id = $1 AND region = $2;

-- name: DoesPlayerExist :one
SELECT EXISTS(SELECT 1 FROM players WHERE player_id = $1 AND region = $2);

-- name: DoesGuildExist :one
SELECT EXISTS(SELECT 1 FROM guilds WHERE guild_id = $1 AND region = $2);

-- name: DoesAllianceExist :one
SELECT EXISTS(SELECT 1 FROM alliances WHERE alliance_id = $1 AND region = $2);

-- name: GetPlayerHistory :many
(SELECT 'guild' AS type, g.guild_id AS id, g.name, NULL AS tag, pgm.first_seen, pgm.last_seen
FROM player_guild_memberships pgm
JOIN guilds g ON pgm.guild_id = g.guild_id AND pgm.region = g.region
WHERE pgm.player_id = $1 AND pgm.region = $2
LIMIT $3)
UNION ALL
(SELECT 'alliance' AS type, a.alliance_id AS id, COALESCE(a.name, ''), a.tag, pam.first_seen, pam.last_seen
FROM player_alliance_memberships pam
JOIN alliances a ON pam.alliance_id = a.alliance_id AND pam.region = a.region
WHERE pam.player_id = $1 AND pam.region = $2
LIMIT $3);

-- name: GetGuildHistory :many
(SELECT 'player' AS type, p.player_id AS id, p.name, NULL AS tag, pgm.first_seen, pgm.last_seen
FROM player_guild_memberships pgm
JOIN players p ON pgm.player_id = p.player_id AND pgm.region = p.region
WHERE pgm.guild_id = $1 AND pgm.region = $2
LIMIT $3)
UNION ALL
(SELECT 'alliance' AS type, a.alliance_id AS id, COALESCE(a.name, ''), a.tag, gam.first_seen, gam.last_seen
FROM guild_alliance_memberships gam
JOIN alliances a ON gam.alliance_id = a.alliance_id AND gam.region = a.region
WHERE gam.guild_id = $1 AND gam.region = $2
LIMIT $3);

-- name: GetAllianceHistory :many
(SELECT 'player' AS type, p.player_id AS id, p.name, pam.first_seen, pam.last_seen
FROM player_alliance_memberships pam
JOIN players p ON pam.player_id = p.player_id AND pam.region = p.region
WHERE pam.alliance_id = $1 AND pam.region = $2
LIMIT $3)
UNION ALL
(SELECT 'guild' AS type, g.guild_id AS id, g.name, gam.first_seen, gam.last_seen
FROM guild_alliance_memberships gam
JOIN guilds g ON gam.guild_id = g.guild_id AND gam.region = g.region
WHERE gam.alliance_id = $1 AND gam.region = $2
LIMIT $3);

-- name: GetPlayerGuilds :many
SELECT g.guild_id AS id, g.name, pgm.first_seen, pgm.last_seen
FROM player_guild_memberships pgm
JOIN guilds g ON pgm.guild_id = g.guild_id AND pgm.region = g.region
WHERE pgm.player_id = $1 AND pgm.region = $2
LIMIT $3 OFFSET $4;

-- name: GetPlayerAlliances :many
SELECT a.alliance_id AS id, COALESCE(a.name, ''), a.tag, pam.first_seen, pam.last_seen
FROM player_alliance_memberships pam
JOIN alliances a ON pam.alliance_id = a.alliance_id AND pam.region = a.region
WHERE pam.player_id = $1 AND pam.region = $2
LIMIT $3 OFFSET $4;

-- name: GetGuildMembers :many
SELECT p.player_id AS id, p.name, pgm.first_seen, pgm.last_seen
FROM player_guild_memberships pgm
JOIN players p ON pgm.player_id = p.player_id AND pgm.region = p.region
WHERE pgm.guild_id = $1 AND pgm.region = $2
LIMIT $3 OFFSET $4;

-- name: GetGuildAlliances :many
SELECT a.alliance_id AS id, COALESCE(a.name, ''), a.tag, gam.first_seen, gam.last_seen
FROM guild_alliance_memberships gam
JOIN alliances a ON gam.alliance_id = a.alliance_id AND gam.region = a.region
WHERE gam.guild_id = $1 AND gam.region = $2
LIMIT $3 OFFSET $4;

-- name: GetAllianceGuilds :many
SELECT g.guild_id AS id, g.name, gam.first_seen, gam.last_seen
FROM guild_alliance_memberships gam
JOIN guilds g ON gam.guild_id = g.guild_id AND gam.region = g.region
WHERE gam.alliance_id = $1 AND gam.region = $2
LIMIT $3 OFFSET $4;
