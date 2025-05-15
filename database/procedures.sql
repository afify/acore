-- Procedure to create a new user
CREATE OR REPLACE PROCEDURE create_new_user(
	IN _username VARCHAR,
	IN _email VARCHAR,
	IN _password_hash TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
	INSERT INTO users (username, email, password_hash, created_at, updated_at)
	VALUES (_username, _email, _password_hash, NOW(), NOW());
	RAISE NOTICE 'User % created', _username;
END;
$$;

-- Procedure to update a user's email by ID
CREATE OR REPLACE PROCEDURE update_user_email(
	IN _user_id UUID,
	IN _new_email VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
	UPDATE users
	SET email = _new_email, updated_at = NOW()
	WHERE id = _user_id;
	RAISE NOTICE 'User % email updated to %', _user_id, _new_email;
END;
$$;

-- Procedure to delete a user by ID
CREATE OR REPLACE PROCEDURE delete_user_by_id(
	IN _user_id UUID
)
LANGUAGE plpgsql
AS $$
BEGIN
	DELETE FROM users WHERE id = _user_id;
	RAISE NOTICE 'User % deleted', _user_id;
END;
$$;

-- Procedure to create a new session
CREATE OR REPLACE PROCEDURE create_new_session(
    IN _user_id UUID,
    IN _session_token TEXT,
    IN _ip_address VARCHAR,
    IN _user_agent TEXT,
    IN _expires_at TIMESTAMPTZ  -- Use TIMESTAMPTZ for timestamp with timezone
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO sessions (user_id, session_token, ip_address, user_agent, expires_at, created_at)
    VALUES (_user_id, _session_token, _ip_address, _user_agent, _expires_at, NOW());
    RAISE NOTICE 'Session % created for user %', _session_token, _user_id;
END;
$$;

-- Procedure to delete a session by session token
CREATE OR REPLACE PROCEDURE delete_session_by_token(
	IN _session_token TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
	DELETE FROM sessions WHERE session_token = _session_token;
	RAISE NOTICE 'Session % deleted', _session_token;
END;
$$;
