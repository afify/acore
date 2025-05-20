DROP FUNCTION IF EXISTS create_user_session;
CREATE FUNCTION create_user_session(
	p_user_id        UUID,
	p_session_token  TEXT,
	p_ip_address     VARCHAR(145),
	p_user_agent     TEXT,
	p_expires_at     TIMESTAMP WITHOUT TIME ZONE
)
RETURNS user_sessions
LANGUAGE plpgsql AS $$
DECLARE
new_row user_sessions;
BEGIN
	INSERT INTO user_sessions (
		user_id,
		session_token,
		ip_address,
		user_agent,
		expires_at
	)
	VALUES (
		p_user_id,
		p_session_token,
		p_ip_address,
		p_user_agent,
		p_expires_at
	)
	RETURNING
	id,
	user_id,
	session_token,
	ip_address,
	user_agent,
	created_at,
	expires_at
	INTO new_row;

	RETURN new_row;
END;
$$;
