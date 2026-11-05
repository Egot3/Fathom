CREATE INDEX IF NOT EXISTS idx_groups_users_user_group
    ON "groups_users" ("user_uuid", "group_uuid");

CREATE INDEX IF NOT EXISTS idx_tests_quizzes_quiz
    ON "tests_quizzes" ("quiz_uuid");

CREATE INDEX IF NOT EXISTS idx_ugt_user_finalized
    ON "users_groups_tests" ("user_uuid", "finalized_at" DESC);

CREATE INDEX IF NOT EXISTS idx_ugt_group_finalized
    ON "users_groups_tests" ("group_uuid", "finalized_at" DESC);

CREATE INDEX IF NOT EXISTS idx_answers_quiz
    ON "users_groups_tests_quiz_answers" ("quiz_uuid");

CREATE INDEX IF NOT EXISTS idx_users_created_active
    ON "users" ("created_at") WHERE "deleted_at" IS NULL;
