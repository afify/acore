CALL create_new_user('user1', 'user1@example.com', crypt('password1', gen_salt('bf')));
CALL create_new_user('user2', 'user2@example.com', crypt('password2', gen_salt('bf')));
CALL create_new_user('user3', 'user3@example.com', crypt('password3', gen_salt('bf')));
CALL create_new_user('user4', 'user4@example.com', crypt('password4', gen_salt('bf')));
CALL create_new_user('user5', 'user5@example.com', crypt('password5', gen_salt('bf')));
CALL create_new_user('user6', 'user6@example.com', crypt('password6', gen_salt('bf')));
CALL create_new_user('user7', 'user7@example.com', crypt('password7', gen_salt('bf')));
CALL create_new_user('user8', 'user8@example.com', crypt('password8', gen_salt('bf')));
CALL create_new_user('user9', 'user9@example.com', crypt('password9', gen_salt('bf')));
CALL create_new_user('user10', 'user10@example.com', crypt('password10', gen_salt('bf')));

DO $$
	DECLARE
	user_id UUID;
	BEGIN
	FOR user_id IN
		SELECT id FROM users
		LOOP
			INSERT INTO user_settings (user_id, is_dark_mode, preferred_language, notification_settings, created_at, updated_at)
			VALUES (user_id, FALSE, 'en', '{"email": true, "sms": false}', NOW(), NOW());
		END LOOP;
	END $$;

DO $$
	DECLARE
	user_id UUID;
	BEGIN
	FOR user_id IN
		SELECT id FROM users WHERE username IN ('user1', 'user2', 'user3', 'user4', 'user5')
		LOOP
			CALL create_new_session(user_id, md5(random()::text)::text, '127.0.0.1'::varchar, 'Mozilla/5.0'::varchar, NOW() + INTERVAL '7 days');
		END LOOP;
END $$;
