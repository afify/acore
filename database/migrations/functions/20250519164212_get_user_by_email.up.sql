DROP FUNCTION IF EXISTS get_user_by_email;
CREATE FUNCTION get_user_by_email(
	p_username_email TEXT
	) RETURNS TABLE (
	id            UUID,
	username      VARCHAR(50),
	email         VARCHAR(100),
	password_hash TEXT,
	created_at    TIMESTAMP WITHOUT TIME ZONE,
	updated_at    TIMESTAMP WITHOUT TIME ZONE
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
WHERE email = p_username_email
OR username = p_username_email;
$$;
