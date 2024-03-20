CREATE TYPE region_enum AS ENUM ('americas', 'europe', 'asia');

CREATE TABLE players (
    name VARCHAR(25) NOT NULL,
    player_id VARCHAR(50) NOT NULL,
    region region_enum NOT NULL,
    first_seen TIMESTAMPTZ NOT NULL,
    last_seen TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (player_id, region)
);

CREATE TABLE guilds (
    name VARCHAR(50) NOT NULL,
    guild_id VARCHAR(50) NOT NULL,
    region region_enum NOT NULL,
    first_seen TIMESTAMPTZ NOT NULL,
    last_seen TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (guild_id, region)
);

CREATE TABLE alliances (
    tag VARCHAR(5) NOT NULL,
    alliance_id VARCHAR(50) NOT NULL,
    region region_enum NOT NULL,
    first_seen TIMESTAMPTZ NOT NULL,
    last_seen TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (alliance_id, region)
);

CREATE TABLE player_guild_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id VARCHAR(50) NOT NULL,
    guild_id VARCHAR(50) NOT NULL,
    region region_enum NOT NULL,
    first_seen TIMESTAMPTZ NOT NULL,
    last_seen TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (player_id, region) REFERENCES players(player_id, region),
    FOREIGN KEY (guild_id, region) REFERENCES guilds(guild_id, region)
);

CREATE TABLE guild_alliance_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    guild_id VARCHAR(50) NOT NULL,
    alliance_id VARCHAR(50) NOT NULL,
    region region_enum NOT NULL,
    first_seen TIMESTAMPTZ NOT NULL,
    last_seen TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (guild_id, region) REFERENCES guilds(guild_id, region),
    FOREIGN KEY (alliance_id, region) REFERENCES alliances(alliance_id, region)
);

CREATE TABLE player_alliance_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id VARCHAR(50) NOT NULL,
    alliance_id VARCHAR(50) NOT NULL,
    region region_enum NOT NULL,
    first_seen TIMESTAMPTZ NOT NULL,
    last_seen TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (player_id, region) REFERENCES players(player_id, region),
    FOREIGN KEY (alliance_id, region) REFERENCES alliances(alliance_id, region)
);
