-- Create table "users"
CREATE TABLE IF NOT EXISTS "users" (
    "uuid" UUID PRIMARY KEY,
    "nickname" VARCHAR(255) UNIQUE NOT NULL,
    "password_hash" CHAR(60) NOT NULL
);

-- Create table "groups" 
CREATE TABLE IF NOT EXISTS "groups" (
    "uuid" UUID PRIMARY KEY,
    "name" VARCHAR(255) UNIQUE
);

-- Create table "groups_users" 
CREATE TABLE IF NOT EXISTS "groups_users" (
    "group_uuid" UUID REFERENCES groups(uuid) ON DELETE CASCADE,
    "user_uuid"  UUID REFERENCES users(uuid) ON DELETE CASCADE,
    PRIMARY KEY("group_uuid","user_uuid")
);

-- Create table "quizes"
CREATE TABLE IF NOT EXISTS "quizes" (
    "uuid" UUID PRIMARY KEY,
    "path" TEXT UNIQUE
);

-- Create table "tests"
CREATE TABLE IF NOT EXISTS "tests" (
    "uuid" UUID         PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "max_score" SMALLINT NOT NULL
);

-- Create table "tests_quizes"
CREATE TABLE IF NOT EXISTS "tests_quizes" (
    "test_uuid" UUID REFERENCES tests(uuid) ON DELETE CASCADE,
    "quiz_uuid" UUID REFERENCES quizes(uuid) ON DELETE CASCADE,
    PRIMARY KEY("test_uuid","quiz_uuid")
);

-- Create table "users_groups_tests"
CREATE TABLE IF NOT EXISTS "users_groups_tests" (
    "test_uuid"  UUID REFERENCES tests(uuid) ON DELETE CASCADE,
    "group_uuid" UUID REFERENCES groups(uuid) ON DELETE CASCADE,
    "user_uuid"  UUID REFERENCES users(uuid) ON DELETE CASCADE,
    "score"      SMALLINT NOT NULL,
    
    PRIMARY KEY("test_uuid","group_uuid","user_uuid")
);
