-- name: GetCountsOfEntities :one
SELECT
  (SELECT COUNT(*) FROM players) AS total_players,
  (SELECT COUNT(*) FROM guilds) AS total_guilds,
  (SELECT COUNT(*) FROM alliances) AS total_alliances;

-- name: FindPlayersByNameAndRegion :many
SELECT player_id AS id, name
FROM players
WHERE name ILIKE @name || '%' AND region = @region;

-- name: GetPlayerByIdAndRegion :one
SELECT player_id AS id, name, first_seen, last_seen
FROM players
WHERE player_id = $1 AND region = $2;

-- name: GetPlayerGuildMemberships :many
SELECT g.guild_id AS id, g.name, pgm.first_seen, pgm.last_seen
FROM player_guild_memberships pgm
JOIN guilds g ON pgm.guild_id = g.guild_id AND pgm.region = g.region
WHERE pgm.player_id = $1 AND pgm.region = $2;

-- name: GetPlayerAllianceMemberships :many
SELECT a.alliance_id AS id, a.tag, pam.first_seen, pam.last_seen
FROM player_alliance_memberships pam
JOIN alliances a ON pam.alliance_id = a.alliance_id AND pam.region = a.region
WHERE pam.player_id = $1 AND pam.region = $2;
