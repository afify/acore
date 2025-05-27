CREATE OR REPLACE FUNCTION get_user_by_provider(
	p_provider_id SMALLINT,
	p_sub TEXT
)
RETURNS TABLE(user_id UUID)
AS $$
SELECT up.user_id
FROM user_providers up
WHERE up.auth_provider_id = p_provider_id
AND up.provider_sub       = p_sub;
$$ LANGUAGE sql STRICT;
