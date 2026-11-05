UPDATE "users" SET "is_teacher" = FALSE WHERE "is_teacher" IS NULL;
ALTER TABLE "users"
    ALTER COLUMN "is_teacher" SET DEFAULT FALSE,
    ALTER COLUMN "is_teacher" SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS "quizzes_path_key" ON "quizzes" ("path");

ALTER TABLE "tests_quizzes" ALTER COLUMN "position" SET NOT NULL;

ALTER TABLE "users_groups_tests"
    ALTER COLUMN "finalized_at" DROP NOT NULL;
ALTER TABLE "users_groups_tests"
    ALTER COLUMN "finalized_at" TYPE TIMESTAMPTZ
    USING (CASE WHEN "finalized_at"::text = '' THEN NULL ELSE "finalized_at"::TIMESTAMPTZ END);


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
