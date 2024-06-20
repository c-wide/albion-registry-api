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

-- name: GetPlayerHistory :one
SELECT 
  (
    SELECT array_to_json(array_agg(row_to_json(g)))
    FROM (
      SELECT 
        pgm.guild_id,
        g.name,
        pgm.first_seen,
        pgm.last_seen,
        (
          SELECT array_to_json(array_agg(row_to_json(a)))
          FROM (
            SELECT 
              gam.alliance_id,
              a.name,
              a.tag,
              GREATEST(pgm.first_seen, gam.first_seen) AS first_seen,
              LEAST(pgm.last_seen, gam.last_seen) AS last_seen
            FROM 
              guild_alliance_memberships gam
            JOIN 
              alliances a ON gam.alliance_id = a.alliance_id AND gam.region = a.region
            WHERE 
              gam.guild_id = g.guild_id AND gam.region = g.region
              AND gam.first_seen <= pgm.last_seen
              AND gam.last_seen >= pgm.first_seen
          ) a
        ) AS alliances
      FROM 
        player_guild_memberships pgm
      JOIN 
        guilds g ON pgm.guild_id = g.guild_id AND pgm.region = g.region
      WHERE 
        pgm.player_id = $1 AND pgm.region = $2
    ) g
  ) AS history;
