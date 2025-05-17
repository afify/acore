CREATE OR REPLACE FUNCTION public.get_user_password_hash(
	p_identity TEXT
	) RETURNS TABLE (
	id            UUID,
	password_hash TEXT
)
LANGUAGE SQL AS $$
SELECT id, password_hash
FROM users
WHERE email = p_identity
OR username = p_identity;
$$;
