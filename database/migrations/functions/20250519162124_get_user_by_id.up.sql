DROP FUNCTION IF EXISTS get_user_by_id;
CREATE FUNCTION get_user_by_id(
	p_id UUID
)
RETURNS TABLE (
	id         UUID,
	username   VARCHAR(50),
	email      VARCHAR(100),
	created_at TIMESTAMP WITHOUT TIME ZONE,
	updated_at TIMESTAMP WITHOUT TIME ZONE
)
LANGUAGE SQL
AS $$
SELECT
	id,
	username,
	email,
	created_at,
	updated_at
FROM users
WHERE id = p_id;
$$;
