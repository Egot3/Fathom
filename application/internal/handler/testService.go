package handler

import (
	"net/http"

	"github.com/egot3/fathom/internal/database/repositories/answer"
	"github.com/egot3/fathom/internal/database/repositories/group"
	"github.com/egot3/fathom/internal/database/repositories/quiz"
	"github.com/egot3/fathom/internal/database/repositories/test"
	"github.com/egot3/fathom/internal/database/repositories/user"
	quizparser "github.com/egot3/fathom/internal/quizParser"
	"github.com/samber/do/v2"
)

type chiTestService struct {
	userRepo   user.UserRepository
	groupRepo  group.GroupRepository
	quizRepo   quiz.QuizRepository
	testRepo   test.TestRepository
	answerRepo answer.AnswerRepository

	runner *quizparser.TestRunner
}

type TestService interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	PatchUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
}

func NewTestService(i do.Injector) (TestService, error) {
	uR := do.MustInvoke[user.UserRepository](i)
	gR := do.MustInvoke[group.GroupRepository](i)
	qR := do.MustInvoke[quiz.QuizRepository](i)
	tR := do.MustInvoke[test.TestRepository](i)
	aR := do.MustInvoke[answer.AnswerRepository](i)

	r := do.MustInvoke[*quizparser.TestRunner](i)

	return &chiTestService{userRepo: uR,
		quizRepo:   qR,
		groupRepo:  gR,
		testRepo:   tR,
		answerRepo: aR,
		runner:     r,
	}, nil
}
