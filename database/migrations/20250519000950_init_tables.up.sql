CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
	id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	username        VARCHAR(50)                 NOT NULL UNIQUE,
	email           VARCHAR(100)                NOT NULL UNIQUE,
	password_hash   TEXT                        NOT NULL,
	created_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);
--------------------------
CREATE TABLE IF NOT EXISTS user_sessions (
	id            UUID                        PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id       UUID                        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	session_token TEXT                        NOT NULL UNIQUE,
	ip_address    VARCHAR(145)                    NULL,
	user_agent    TEXT                            NULL,
	created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	expires_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
