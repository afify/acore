DROP FUNCTION IF EXISTS get_user_sessions;
CREATE OR REPLACE FUNCTION get_user_sessions(
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
