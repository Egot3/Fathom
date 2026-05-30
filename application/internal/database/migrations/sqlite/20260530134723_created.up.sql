PRAGMA foreign_keys = ON;

-- Create table "users"
CREATE TABLE IF NOT EXISTS "users" (
    "uuid"          TEXT PRIMARY KEY,
    "nickname"      TEXT UNIQUE NOT NULL,
    "password_hash" TEXT NOT NULL
);

-- Create table "groups" 
CREATE TABLE IF NOT EXISTS "groups" (
    "uuid" TEXT PRIMARY KEY,
    "name" TEXT UNIQUE
);

-- Create table "groups_users" 
CREATE TABLE IF NOT EXISTS "groups_users" (
    "group_uuid" TEXT REFERENCES groups(uuid) ON DELETE CASCADE,
    "user_uuid"  TEXT REFERENCES users(uuid) ON DELETE CASCADE,
    PRIMARY KEY ("group_uuid", "user_uuid")
);

-- Create table "quizes"
CREATE TABLE IF NOT EXISTS "quizes" (
    "uuid" TEXT PRIMARY KEY,
    "path" TEXT UNIQUE
);

-- Create table "tests"
CREATE TABLE IF NOT EXISTS "tests" (
    "uuid"      TEXT PRIMARY KEY,
    "name"      TEXT NOT NULL UNIQUE,
    "max_score" INTEGER NOT NULL
);

-- Create table "tests_quizes"
CREATE TABLE IF NOT EXISTS "tests_quizes" (
    "test_uuid" TEXT REFERENCES tests(uuid) ON DELETE CASCADE,
    "quiz_uuid" TEXT REFERENCES quizes(uuid) ON DELETE CASCADE,
    PRIMARY KEY ("test_uuid", "quiz_uuid")
);

-- Create table "users_groups_tests"
CREATE TABLE IF NOT EXISTS "users_groups_tests" (
    "test_uuid"  TEXT REFERENCES tests(uuid) ON DELETE CASCADE,
    "group_uuid" TEXT REFERENCES groups(uuid) ON DELETE CASCADE,
    "user_uuid"  TEXT REFERENCES users(uuid) ON DELETE CASCADE,
    "score"      INTEGER NOT NULL,
    PRIMARY KEY ("test_uuid", "group_uuid", "user_uuid")
);