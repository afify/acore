CREATE FUNCTION public.get_user_by_email(
	p_email TEXT
	) RETURNS TABLE (
	id            UUID,
	username      VARCHAR(50),
	email         VARCHAR(100),
	password_hash TEXT,
	created_at    TIMESTAMPTZ,
	updated_at    TIMESTAMPTZ
)
LANGUAGE SQL AS $$
SELECT
	id,
	username,
	email,
	password_hash,
	created_at,
	updated_at
FROM users
WHERE email = p_email
OR username = p_email;
$$;
