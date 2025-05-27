DROP FUNCTION IF EXISTS get_user_sessions;
CREATE OR REPLACE FUNCTION get_user_sessions(
	p_user_id UUID
)
RETURNS TABLE (
	id               UUID,
	user_id          UUID,
	session_type_id  SMALLINT,
	auth_provider_id SMALLINT,
	session_token    TEXT,
	ip_address       VARCHAR(128),
	user_agent       TEXT,
	created_at       TIMESTAMP WITHOUT TIME ZONE,
	expires_at       TIMESTAMP WITHOUT TIME ZONE
)
LANGUAGE SQL
STABLE
AS $$
SELECT
	id,
	user_id,
	session_type_id,
	auth_provider_id,
	session_token,
	ip_address,
	user_agent,
	created_at,
	expires_at
FROM user_sessions
WHERE user_id = p_user_id;
$$;
