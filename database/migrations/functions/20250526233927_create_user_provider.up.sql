DROP FUNCTION IF EXISTS create_user_provider;
CREATE FUNCTION create_user_provider(
	p_user_id        UUID,
	p_provider_id    SMALLINT,
	p_sub            TEXT
)
RETURNS TABLE(user_id UUID)
LANGUAGE sql
AS $$
INSERT INTO user_providers (
	user_id,
	auth_provider_id,
	provider_sub
)
VALUES (
	p_user_id,
	p_provider_id,
	p_sub
)
ON CONFLICT (user_id, auth_provider_id) DO NOTHING;

SELECT p_user_id AS user_id;
$$;
