-- name: GetCountsOfEntities :one
SELECT
  (SELECT COUNT(*) FROM players) AS total_players,
  (SELECT COUNT(*) FROM guilds) AS total_guilds,
  (SELECT COUNT(*) FROM alliances) AS total_alliances;

-- name: FindPlayersByNameAndRegion :many
SELECT player_id, name
FROM players
WHERE name ILIKE @name || '%' AND region = @region;
