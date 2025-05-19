DROP FUNCTION IF EXISTS create_user;
CREATE FUNCTION create_user(
	p_username      VARCHAR(50),
	p_email         VARCHAR(100),
	p_password_hash TEXT
) RETURNS UUID
LANGUAGE plpgsql AS $$
DECLARE
new_id UUID;
BEGIN
	INSERT INTO users (username, email, password_hash, created_at, updated_at)
	VALUES (p_username, p_email, p_password_hash, NOW(), NOW())
	RETURNING id INTO new_id;
	RETURN new_id;
END;
$$;
