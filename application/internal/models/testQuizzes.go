package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TestsQuizzies struct {
	bun.BaseModel `bun:"table:tests_quizzes"`

	TestUUID uuid.UUID `bun:"test_uuid,pk"`
	Test     *Test     `bun:"rel:belongs-to,join:test_uuid=uuid"`

	QuizPath string `bun:"quiz_path,pk"`
	Quiz     *Quiz  `bun:"rel:belongs-to,join:quiz_path=path"`
}
