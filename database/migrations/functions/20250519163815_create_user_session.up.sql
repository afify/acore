DROP FUNCTION IF EXISTS create_user_session;
CREATE FUNCTION create_user_session(
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
