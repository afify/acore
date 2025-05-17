CREATE OR REPLACE FUNCTION public.get_user_by_id(
	p_id UUID
	) RETURNS TABLE (
	id            UUID,
	username      VARCHAR,
	email         VARCHAR,
	created_at    TIMESTAMPTZ,
	updated_at    TIMESTAMPTZ
)
LANGUAGE SQL AS $$
SELECT
id,
username,
email,
created_at,
updated_at
FROM users
WHERE id = p_id;
$$;
