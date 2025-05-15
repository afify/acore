CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS user_role_mappings CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS user_activity_logs CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS user_settings CASCADE;
DROP TABLE IF EXISTS authentication CASCADE;
DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	username VARCHAR(50) UNIQUE NOT NULL,
	email VARCHAR(100) UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS authentication (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	auth_provider VARCHAR(50),
	auth_token TEXT,
	password_reset_token TEXT,
	token_expires_at TIMESTAMP,
	created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_settings (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	is_dark_mode BOOLEAN DEFAULT FALSE,
	preferred_language VARCHAR(50) DEFAULT 'en',
	notification_settings JSONB,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW()
);

-- Create sessions table for storing user sessions
CREATE TABLE IF NOT EXISTS sessions (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	session_token TEXT UNIQUE NOT NULL,  -- Unique session token for the user
	ip_address VARCHAR(45),  -- IPv6 compatible
	user_agent TEXT,  -- Browser or app information
	expires_at TIMESTAMP,  -- Session expiration
	created_at TIMESTAMP DEFAULT NOW()
);

-- Create a table to log all user activities (for auditing and tracking)
CREATE TABLE IF NOT EXISTS user_activity_logs (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	activity_type VARCHAR(100),  -- Type of activity (e.g., 'login', 'update_profile', 'logout')
	activity_details TEXT,  -- Detailed description of the activity
	created_at TIMESTAMP DEFAULT NOW()
);

-- Create a table for user roles and permissions
CREATE TABLE IF NOT EXISTS user_roles (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	role_name VARCHAR(50) UNIQUE NOT NULL  -- Role (e.g., 'admin', 'user', 'moderator')
);

-- Create a table to assign roles to users
CREATE TABLE IF NOT EXISTS user_role_mappings (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	role_id UUID REFERENCES user_roles(id) ON DELETE CASCADE,
	assigned_at TIMESTAMP DEFAULT NOW()
);

-- Create an audit table for tracking changes to user data
CREATE TABLE IF NOT EXISTS audit_logs (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	change_type VARCHAR(50),  -- What type of change (e.g., 'update', 'delete')
	old_data JSONB,  -- Previous data
	new_data JSONB,  -- Updated data
	changed_at TIMESTAMP DEFAULT NOW()
);
