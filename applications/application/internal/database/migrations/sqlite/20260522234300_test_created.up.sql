PRAGMA foreign_keys = ON;

-- Create table "users"
CREATE TABLE IF NOT EXISTS "users" (
    "uuid"         TEXT PRIMARY KEY,
    "nickname"     TEXT UNIQUE NOT NULL,
    "password_hash" TEXT NOT NULL,
    "is_teacher"   INTEGER NOT NULL DEFAULT 0,

    "deleted_at"   TEXT DEFAULT NULL,
    "updated_at"   TEXT,
    "created_at"   TEXT
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

-- Create table "quizzes"
CREATE TABLE IF NOT EXISTS "quizzes" (
    "uuid"             TEXT PRIMARY KEY,
    "path"             TEXT UNIQUE,
    "checksum"         TEXT NOT NULL,          
    "score"            INTEGER NOT NULL,
    "correct_answer"   TEXT NOT NULL
);

-- Create table "tests"
CREATE TABLE IF NOT EXISTS "tests" (
    "uuid"       TEXT PRIMARY KEY,
    "name"       TEXT NOT NULL UNIQUE,
    "created_at" TEXT NOT NULL,
    "updated_at" TEXT NOT NULL
);

-- Create table "tests_quizzes"
CREATE TABLE IF NOT EXISTS "tests_quizzes" (
    "test_uuid" TEXT REFERENCES tests(uuid) ON DELETE CASCADE,
    "quiz_uuid" TEXT REFERENCES quizzes("uuid") ON DELETE CASCADE,
    "position" INTEGER NOT NULL,

    PRIMARY KEY ("test_uuid", "quiz_uuid")
);

-- Create table "users_groups_tests"
CREATE TABLE IF NOT EXISTS "users_groups_tests" (
    "test_uuid"  TEXT REFERENCES tests(uuid) ON DELETE CASCADE,
    "group_uuid" TEXT REFERENCES groups(uuid) ON DELETE CASCADE,
    "user_uuid"  TEXT REFERENCES users(uuid) ON DELETE CASCADE,

    "score"      REAL NOT NULL,

    PRIMARY KEY ("test_uuid", "group_uuid", "user_uuid")
);

-- Create table "users_groups_tests_quiz_answers"
CREATE TABLE IF NOT EXISTS "users_groups_tests_quiz_answers" (
    "test_uuid"    TEXT REFERENCES tests(uuid) ON DELETE CASCADE,
    "group_uuid"   TEXT REFERENCES groups(uuid) ON DELETE CASCADE,
    "user_uuid"    TEXT REFERENCES users(uuid) ON DELETE CASCADE,
    "quiz_uuid"    TEXT REFERENCES quizzes(uuid) ON DELETE CASCADE,
    "score"        REAL NOT NULL, -- for real?
    "answer_value" TEXT NOT NULL,

    "answered_at" TEXT NOT NULL,

    PRIMARY KEY ("test_uuid", "group_uuid", "user_uuid", "quiz_uuid", "answered_at")
);