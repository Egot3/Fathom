package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Answer struct {
	bun.BaseModel `bun:"table:users_groups_tests_quiz_answers,alias:ugtqa"` //hell of an alias

	TestUUID uuid.UUID `bun:"test_uuid,pk"`
	Test     *Test     `bun:"rel:belongs-to,join:test_uuid=uuid"`

	GroupUUID uuid.UUID `bun:"group_uuid,pk"`
	Group     *Group    `bun:"rel:belongs-to,join:group_uuid=uuid"`

	UserUUID uuid.UUID `bun:"user_uuid,pk"`
	User     *User     `bun:"rel:belongs-to,join:user_uuid=uuid"`

	QuizPath string `bun:"quiz_path,pk"`
	Quiz     *Quiz  `bun:"rel:belongs-to,join:quiz_path=path"`

	AnsweredAt time.Time `bun:"answered_at,pk"`

	Score uint64 `bun:"score"`
}

var _ bun.BeforeAppendModelHook = (*Answer)(nil)

func (a *Answer) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		a.AnsweredAt = time.Now()
	}
	return nil
}
