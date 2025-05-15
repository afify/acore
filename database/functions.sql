-- Function to get a user by ID
CREATE OR REPLACE FUNCTION get_user(_user_id UUID) RETURNS TABLE (
	id UUID,
	username VARCHAR,
	email VARCHAR,
	created_at TIMESTAMP
) AS $$
BEGIN
	RETURN QUERY SELECT id, username, email, created_at
	FROM users WHERE id = _user_id;
END;
$$ LANGUAGE plpgsql;

-- Function to get all sessions for a user by user ID
CREATE OR REPLACE FUNCTION get_user_sessions(_user_id UUID) RETURNS TABLE (
	session_id UUID,
	session_token TEXT,
	created_at TIMESTAMP,
	expires_at TIMESTAMP
) AS $$
BEGIN
	RETURN QUERY SELECT id, session_token, created_at, expires_at
	FROM sessions WHERE user_id = _user_id;
END;
$$ LANGUAGE plpgsql;

-- Function to get user settings by user ID
CREATE OR REPLACE FUNCTION get_user_settings(_user_id UUID) RETURNS TABLE (
	id UUID,
	is_dark_mode BOOLEAN,
	preferred_language VARCHAR,
	notification_settings JSONB
) AS $$
BEGIN
	RETURN QUERY SELECT id, is_dark_mode, preferred_language, notification_settings
	FROM user_settings WHERE user_id = _user_id;
END;
$$ LANGUAGE plpgsql;

-- Function to retrieve all users
CREATE OR REPLACE FUNCTION get_all_users() RETURNS TABLE (
	id UUID,
	username VARCHAR,
	email VARCHAR,
	created_at TIMESTAMP
) AS $$
BEGIN
	RETURN QUERY SELECT id, username, email, created_at FROM users;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION check_user_exists(_username VARCHAR, _email VARCHAR) 
RETURNS TABLE (username_exists BOOLEAN, email_exists BOOLEAN) AS $$
BEGIN
	-- Check if the username exists
	SELECT EXISTS(SELECT 1 FROM users WHERE username = _username) INTO username_exists;

	-- Check if the email exists
	SELECT EXISTS(SELECT 1 FROM users WHERE email = _email) INTO email_exists;

	RETURN;
END;
$$ LANGUAGE plpgsql;
