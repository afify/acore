CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
	id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	username        VARCHAR(50)  NOT NULL UNIQUE,
	email           VARCHAR(100) NOT NULL UNIQUE,
	password_hash   TEXT         NOT NULL,
	created_at      TIMESTAMP    NOT NULL DEFAULT NOW(),
	updated_at      TIMESTAMP    NOT NULL DEFAULT NOW()
);
