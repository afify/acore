CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
	id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	username        VARCHAR(50)  NOT NULL UNIQUE,
	email           VARCHAR(255) NOT NULL UNIQUE,
	password_hash   VARCHAR(128)     NULL,
	created_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);
-------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS session_types (
	id   SMALLSERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL UNIQUE
);
-------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS auth_providers (
	id   SMALLSERIAL       PRIMARY KEY,
	name VARCHAR(50)       NOT NULL UNIQUE
);
-------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_sessions (
	id               UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id          UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	session_type_id  SMALLINT     NOT NULL REFERENCES session_types(id) ON DELETE RESTRICT,
	auth_provider_id SMALLINT     NOT NULL REFERENCES auth_providers(id) ON DELETE RESTRICT,
	session_token    VARCHAR(300) NOT NULL UNIQUE,
	ip_address       VARCHAR(128)     NULL,
	user_agent       TEXT             NULL,
	created_at       TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	expires_at       TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
-------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_providers (
	user_id            UUID            NOT NULL REFERENCES users(id)              ON DELETE CASCADE,
	auth_provider_id   SMALLINT        NOT NULL REFERENCES auth_providers(id)     ON DELETE CASCADE,
	provider_sub       VARCHAR(255)    NOT NULL,
	is_email_verified  BOOLEAN         NOT NULL DEFAULT FALSE,
	is_private_email   BOOLEAN         NOT NULL DEFAULT FALSE,
	linked_at          TIMESTAMP       NOT NULL DEFAULT NOW(),
	PRIMARY KEY (user_id, auth_provider_id)
);
