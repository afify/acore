DROP FUNCTION IF EXISTS create_user;
CREATE FUNCTION create_user(
	p_username      VARCHAR(50),
	p_email         VARCHAR(50),
	p_password_hash VARCHAR(128)
)
RETURNS TABLE (
	id          UUID,
	username    VARCHAR(50),
	email       VARCHAR(50),
	created_at  TIMESTAMP WITHOUT TIME ZONE,
	updated_at  TIMESTAMP WITHOUT TIME ZONE
)
LANGUAGE plpgsql
AS $$
BEGIN
	RETURN QUERY
	INSERT INTO users(username, email, password_hash, created_at, updated_at)
	VALUES (p_username, p_email, p_password_hash, NOW(), NOW())
	RETURNING
		users.id,
		users.username,
		users.email,
		users.created_at,
		users.updated_at;
END;
$$;
