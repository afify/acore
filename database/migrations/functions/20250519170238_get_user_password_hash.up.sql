DROP FUNCTION IF EXISTS get_user_password_hash;
CREATE FUNCTION get_user_password_hash(
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
