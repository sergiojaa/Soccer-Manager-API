-- =========================
-- ENUMS
-- =========================

CREATE TYPE player_position AS ENUM ('GK', 'DEF', 'MID', 'ATT');

CREATE TYPE transfer_listing_status AS ENUM ('ACTIVE', 'SOLD', 'CANCELLED');


-- =========================
-- USERS
-- =========================

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_email_not_blank_chk
        CHECK (BTRIM(email) <> '')
);


-- =========================
-- TEAMS
-- =========================

CREATE TABLE teams (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    name TEXT NOT NULL DEFAULT 'My Team',
    country TEXT NOT NULL DEFAULT 'Unknown',
    budget BIGINT NOT NULL DEFAULT 5000000,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT teams_user_fk
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT teams_name_not_blank_chk
        CHECK (BTRIM(name) <> ''),

    CONSTRAINT teams_country_not_blank_chk
        CHECK (BTRIM(country) <> ''),

    CONSTRAINT teams_budget_nonnegative_chk
        CHECK (budget >= 0)
);


-- =========================
-- PLAYERS
-- =========================

CREATE TABLE players (
    id BIGSERIAL PRIMARY KEY,
    team_id BIGINT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    country TEXT NOT NULL,
    age SMALLINT NOT NULL,
    position player_position NOT NULL,
    market_value BIGINT NOT NULL DEFAULT 1000000,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT players_team_fk
        FOREIGN KEY (team_id)
        REFERENCES teams(id)
        ON DELETE CASCADE,

    CONSTRAINT players_first_name_not_blank_chk
        CHECK (BTRIM(first_name) <> ''),

    CONSTRAINT players_last_name_not_blank_chk
        CHECK (BTRIM(last_name) <> ''),

    CONSTRAINT players_country_not_blank_chk
        CHECK (BTRIM(country) <> ''),

    CONSTRAINT players_age_range_chk
        CHECK (age BETWEEN 18 AND 40),

    CONSTRAINT players_market_value_positive_chk
        CHECK (market_value > 0)
);


-- =========================
-- TRANSFER LISTINGS
-- =========================

CREATE TABLE transfer_listings (
    id BIGSERIAL PRIMARY KEY,
    player_id BIGINT NOT NULL,
    seller_team_id BIGINT NOT NULL,
    asking_price BIGINT NOT NULL,
    status transfer_listing_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMPTZ,

    CONSTRAINT transfer_listings_player_fk
        FOREIGN KEY (player_id)
        REFERENCES players(id)
        ON DELETE CASCADE,

    CONSTRAINT transfer_listings_seller_team_fk
        FOREIGN KEY (seller_team_id)
        REFERENCES teams(id)
        ON DELETE CASCADE,

    CONSTRAINT transfer_listings_asking_price_positive_chk
        CHECK (asking_price > 0),

    CONSTRAINT transfer_listings_status_closed_chk
        CHECK (
            (status = 'ACTIVE' AND closed_at IS NULL)
            OR
            (status IN ('SOLD', 'CANCELLED') AND closed_at IS NOT NULL)
        )
);

-- One active listing per player
CREATE UNIQUE INDEX transfer_listings_one_active_per_player_uidx
    ON transfer_listings (player_id)
    WHERE status = 'ACTIVE';


-- =========================
-- TRANSFERS (HISTORY)
-- =========================

CREATE TABLE transfers (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL UNIQUE,
    player_id BIGINT NOT NULL,
    seller_team_id BIGINT NOT NULL,
    buyer_team_id BIGINT NOT NULL,
    sale_price BIGINT NOT NULL,
    market_value_before BIGINT NOT NULL,
    market_value_after BIGINT NOT NULL,
    transferred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT transfers_listing_fk
        FOREIGN KEY (listing_id)
        REFERENCES transfer_listings(id)
        ON DELETE CASCADE,

    CONSTRAINT transfers_player_fk
        FOREIGN KEY (player_id)
        REFERENCES players(id)
        ON DELETE CASCADE,

    CONSTRAINT transfers_seller_team_fk
        FOREIGN KEY (seller_team_id)
        REFERENCES teams(id)
        ON DELETE CASCADE,

    CONSTRAINT transfers_buyer_team_fk
        FOREIGN KEY (buyer_team_id)
        REFERENCES teams(id)
        ON DELETE CASCADE,

    CONSTRAINT transfers_buyer_not_seller_chk
        CHECK (buyer_team_id <> seller_team_id),

    CONSTRAINT transfers_sale_price_positive_chk
        CHECK (sale_price > 0),

    CONSTRAINT transfers_market_value_before_positive_chk
        CHECK (market_value_before > 0),

    CONSTRAINT transfers_market_value_after_positive_chk
        CHECK (market_value_after > 0),

    CONSTRAINT transfers_market_value_growth_chk
        CHECK (market_value_after > market_value_before)
);


-- =========================
-- BASIC INDEXES
-- =========================

CREATE INDEX players_team_idx ON players(team_id);
CREATE INDEX players_team_position_idx ON players(team_id, position);
CREATE INDEX transfers_player_idx ON transfers(player_id);