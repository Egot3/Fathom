package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TestsQuizzies struct {
	bun.BaseModel `bun:"table:tests_quizzes"`

	TestUUID uuid.UUID `bun:"test_uuid,pk"`
	Test     *Test     `bun:"rel:belongs-to,join:test_uuid=uuid"`

	QuizUUID uuid.UUID `bun:"quiz_uuid,pk"`
	Quiz     *Quiz     `bun:"rel:belongs-to,join:quiz_uuid=uuid"`

	Position int `bun:"position,notnull,unique"`
}
