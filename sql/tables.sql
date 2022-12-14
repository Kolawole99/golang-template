-- CREATE DATABASE "graphreach.com" WITH OWNER = postgres ENCODING = 'UTF8' CONNECTION
-- LIMIT = -1;
CREATE TABLE IF NOT EXISTS organization (
    id SERIAL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    created_on TIMESTAMP,
    updated_on TIMESTAMP,
    PRIMARY KEY (id),
);
CREATE INDEX org_id ON organization(id);
CREATE TABLE IF NOT EXISTS user (
    id SERIAL,
    org_id INT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    active_organization INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    last_authenticated TIMESTAMP,
    last_logged_in TIMESTAMP,
    created_on TIMESTAMP,
    updated_on TIMESTAMP,
    PRIMARY KEY (id, email),
    FOREIGN KEY (org_id) REFERENCES organization (id),
);
CREATE INDEX user_id ON user(id);
CREATE INDEX email ON user(email);
CREATE INDEX org_id ON user(org_id);
CREATE TABLE IF NOT EXISTS organization_users (
    id SERIAL,
    org_id INT NOT NULL,
    user_id INT NOT NULL,
    role VARCHAR(15) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    created_on TIMESTAMP,
    updated_on TIMESTAMP,
    PRIMARY KEY (id, org_id, user_id),
    FOREIGN KEY (org_id) REFERENCES organization (id),
    FOREIGN KEY (user_id) REFERENCES user (id)
);
CREATE INDEX org_id_user_id ON organization_users(org_id, user_id);