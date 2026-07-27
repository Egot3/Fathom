-- Create table "users"
CREATE TABLE IF NOT EXISTS "users" (
    "uuid" UUID PRIMARY KEY,
    "nickname" VARCHAR(255) UNIQUE NOT NULL,
    "password_hash" CHAR(60) NOT NULL,
    "is_teacher" BOOLEAN DEFAULT FALSE,

    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "updated_at" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ
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

-- Create table "quizzes"
CREATE TABLE IF NOT EXISTS "quizzes" (
    "uuid" UUID primary key,
    "path" TEXT,
    "checksum" CHAR(32) NOT NULL, --de-facto for "is cached" check
    "score" SMALLINT NOT NULL,
    "correct_answer" TEXT NOT NULL
);

-- Create table "tests"
CREATE TABLE IF NOT EXISTS "tests" (
    "uuid" UUID         PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL
);

-- Create table "tests_quizzes"
CREATE TABLE IF NOT EXISTS "tests_quizzes" (
    "test_uuid" UUID REFERENCES tests(uuid) ON DELETE CASCADE,
    "quiz_uuid" UUID REFERENCES quizzes("uuid") ON DELETE CASCADE,
    "position" INTEGER,

    PRIMARY KEY ("test_uuid", "quiz_uuid")
);

-- Create table "users_groups_tests"
CREATE TABLE IF NOT EXISTS "users_groups_tests" (
    "test_uuid"  UUID REFERENCES tests(uuid) ON DELETE CASCADE,
    "group_uuid" UUID REFERENCES groups(uuid) ON DELETE CASCADE,
    "user_uuid"  UUID REFERENCES users(uuid) ON DELETE CASCADE,
    "finalized_at" TEXT NOT NULL,

    "score"      FLOAT NOT NULL,
    
    PRIMARY KEY("test_uuid","group_uuid","user_uuid")
);

-- Create table "users_groups_tests_quiz_answers"
CREATE TABLE IF NOT EXISTS "users_groups_tests_quiz_answers" (
    "test_uuid"    UUID REFERENCES tests(uuid) ON DELETE CASCADE,
    "group_uuid"   UUID REFERENCES groups(uuid) ON DELETE CASCADE,
    "user_uuid"    UUID REFERENCES users(uuid) ON DELETE CASCADE,
    "quiz_uuid"    UUID REFERENCES quizzes(uuid) ON DELETE CASCADE,
    "score"        FLOAT NOT NULL,
    "answer_value" TEXT NOT NULL,

    "answered_at" TIMESTAMPTZ NOT NULL,

    PRIMARY KEY("test_uuid","group_uuid","user_uuid","quiz_uuid","answered_at")
);