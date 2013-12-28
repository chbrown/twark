-- $ dropdb twark; createdb twark && psql twark < schema.sql

CREATE EXTENSION citext;

-- PG data types: http://www.postgresql.org/docs/9.3/static/datatype.html

CREATE TABLE users (
  -- Twitter docs: 1) Your username cannot be longer than 15 characters, 2) A username can only contain alphanumeric characters (letters A-Z, numbers 0-9) with the exception of underscores, as noted above.
  -- named... CONSTRAINT valid_twitter_username
  screen_name CITEXT PRIMARY KEY CHECK (screen_name ~* '^[_a-z0-9]{1,15}$'),

  statuses_count                     BIGINT,
  contributors_enabled               BOOLEAN,
  friends_count                      BIGINT,
  geo_enabled                        BOOLEAN,
  description                        TEXT,
  profile_sidebar_border_color       TEXT,
  listed_count                       BIGINT,
  followers_count                    BIGINT,
  location                           TEXT,
  profile_background_image_url       TEXT,
  name                               TEXT,
  default_profile_image              BOOLEAN,
  profile_image_url_https            TEXT,
  notifications                      BOOLEAN,
  protected                          BOOLEAN,
  profile_background_color           TEXT,
  created_at                         TEXT,
  default_profile                    BOOLEAN,
  url                                TEXT,
  id_str                             TEXT,
  id                                 BIGINT,
  verified                           BOOLEAN,
  profile_link_color                 TEXT,
  profile_image_url                  TEXT,
  profile_use_background_image       BOOLEAN,
  favourites_count                   BIGINT,
  profile_background_image_url_https TEXT,
  profile_sidebar_fill_color         TEXT,
  is_translator                      BOOLEAN,
  follow_request_sent                BOOLEAN,
  following                          BOOLEAN,
  profile_background_tile            BOOLEAN,
  show_all_inline_media              BOOLEAN,
  profile_text_color                 TEXT,
  lang                               TEXT,

  inserted TIMESTAMP DEFAULT current_timestamp NOT NULL
);

CREATE TABLE tweets (
  -- author CITEXT REFERENCES users,
  author CITEXT CHECK (author ~* '^[_a-z0-9]{1,15}$'),

  id            BIGINT,
  id_str        TEXT,
  created_at    TIMESTAMP,
  text          TEXT,
  source        TEXT,
  retweeted     BOOLEAN,
  retweet_count BIGINT,
  favorited     BOOLEAN,
  truncated     BOOLEAN,
  -- entities      TwitterEntities

  inserted TIMESTAMP DEFAULT current_timestamp NOT NULL,

  PRIMARY KEY (author, id_str)
);
