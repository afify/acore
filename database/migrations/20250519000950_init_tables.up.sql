CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
	id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	username        VARCHAR(50)  NOT NULL UNIQUE,
	email           VARCHAR(100) NOT NULL UNIQUE,
	password_hash   TEXT         NOT NULL,
	created_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at      TIMESTAMP WITHOUT TIME ZONE    NOT NULL DEFAULT NOW()
);
--------------------------
CREATE TABLE IF NOT EXISTS user_sessions (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	session_token TEXT NOT NULL UNIQUE,
	ip_address VARCHAR(45),
	user_agent TEXT,
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
--------------------------
-- Create the function
CREATE FUNCTION public.create_user(
	p_username      VARCHAR(50),
	p_email         VARCHAR(100),
	p_password_hash TEXT
) RETURNS UUID
LANGUAGE plpgsql AS $$
DECLARE
new_id UUID;
BEGIN
	INSERT INTO users (username, email, password_hash, created_at, updated_at)
	VALUES (p_username, p_email, p_password_hash, NOW(), NOW())
	RETURNING id INTO new_id;

	RETURN new_id;
END;
$$;
--------------------------
CREATE FUNCTION public.create_user_session(
	p_user_id        UUID,
	p_session_token  TEXT,
	p_ip_address     VARCHAR(45),
	p_user_agent     TEXT,
	p_expires_at     TIMESTAMP WITHOUT TIME ZONE
) RETURNS UUID
LANGUAGE plpgsql AS $$
DECLARE
new_id UUID;
BEGIN
	INSERT INTO user_sessions (
		user_id,
		session_token,
		ip_address,
		user_agent,
		expires_at
		) VALUES (
		p_user_id,
		p_session_token,
		p_ip_address,
		p_user_agent,
		p_expires_at
	)
	RETURNING id INTO new_id;

	RETURN new_id;
END;
$$;
--------------------------
CREATE FUNCTION public.get_user_by_email(
	p_email TEXT
	) RETURNS TABLE (
	id            UUID,
	username      VARCHAR(50),
	email         VARCHAR(100),
	password_hash TEXT,
	created_at    TIMESTAMP WITHOUT TIME ZONE,
	updated_at    TIMESTAMP WITHOUT TIME ZONE
)
LANGUAGE SQL AS $$
SELECT
	id,
	username,
	email,
	password_hash,
	created_at,
	updated_at
FROM users
WHERE email = p_email
OR username = p_email;
$$;
--------------------------
CREATE OR REPLACE FUNCTION public.get_user_by_id(
	p_id UUID
	) RETURNS TABLE (
	id            UUID,
	username      VARCHAR,
	email         VARCHAR,
	created_at    TIMESTAMP WITHOUT TIME ZONE,
	updated_at    TIMESTAMP WITHOUT TIME ZONE
)
LANGUAGE SQL AS $$
SELECT
id,
username,
email,
created_at,
updated_at
FROM users
WHERE id = p_id;
$$;
--------------------------
CREATE OR REPLACE FUNCTION get_user_sessions_fn(
	p_user_id UUID
)
RETURNS TABLE (
	id             UUID,
	user_id        UUID,
	session_token  TEXT,
	ip_address     VARCHAR(45),
	user_agent     TEXT,
	expires_at     TIMESTAMP WITHOUT TIME ZONE,
	created_at     TIMESTAMP WITHOUT TIME ZONE
)
LANGUAGE SQL
STABLE
AS $$
SELECT
	id,
	user_id,
	session_token,
	ip_address,
	user_agent,
	expires_at,
	created_at
FROM
	user_sessions
WHERE
	user_id = p_user_id;
$$;
--------------------------
