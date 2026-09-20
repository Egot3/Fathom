PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE "tmp_users_groups_tests_quiz_answers" (
    "test_uuid"    TEXT REFERENCES tests(uuid) ON DELETE CASCADE,
    "group_uuid"   TEXT REFERENCES groups(uuid) ON DELETE CASCADE,
    "user_uuid"    TEXT REFERENCES users(uuid) ON DELETE CASCADE,
    "quiz_uuid"    TEXT REFERENCES quizzes(uuid) ON DELETE CASCADE,
    "score"        REAL NOT NULL, -- for real?
    "answer_value" TEXT NOT NULL CHECK (json_valid("answer_value")),

    "answered_at" TEXT NOT NULL,

    PRIMARY KEY ("test_uuid", "group_uuid", "user_uuid", "quiz_uuid", "answered_at")
);

INSERT INTO tmp_users_groups_tests_quiz_answers
    (test_uuid, group_uuid, user_uuid,
    quiz_uuid, score, answer_value, answered_at)
    SELECT test_uuid, group_uuid, user_uuid,
    quiz_uuid, score, answer_value, answered_at FROM users_groups_tests_quiz_answers;

DROP TABLE users_groups_tests_quiz_answers;
ALTER TABLE tmp_users_groups_tests_quiz_answers RENAME TO users_groups_tests_quiz_answers;
COMMIT;
PRAGMA foreign_keys=ON;
