DROP FUNCTION IF EXISTS get_user_session_by_id;
CREATE FUNCTION get_user_session_by_id(
	p_user_id UUID,
	p_token   VARCHAR(300)
)
RETURNS TABLE(id UUID)
LANGUAGE sql
STABLE
AS $$
SELECT us.id
FROM user_sessions us
WHERE us.user_id    = p_user_id
AND us.session_token = p_token
AND us.expires_at > NOW()
LIMIT 1;
$$;
