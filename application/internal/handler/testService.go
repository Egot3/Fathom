package handler

import (
	"net/http"

	"github.com/egot3/fathom/internal/database/repositories/answer"
	"github.com/egot3/fathom/internal/database/repositories/group"
	"github.com/egot3/fathom/internal/database/repositories/quiz"
	"github.com/egot3/fathom/internal/database/repositories/test"
	"github.com/egot3/fathom/internal/database/repositories/user"
	testrunner "github.com/egot3/fathom/internal/testRunner"
	"github.com/samber/do/v2"
)

type chiService struct {
	userRepo   user.UserRepository
	groupRepo  group.GroupRepository
	quizRepo   quiz.QuizRepository
	testRepo   test.TestRepository
	answerRepo answer.AnswerRepository

	runner testrunner.TestRunner
}

type UserService interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	PatchUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
}

type GroupService interface {
	PostGroup(w http.ResponseWriter, r *http.Request)
	GetGroup(w http.ResponseWriter, r *http.Request)
	DeleteGroup(w http.ResponseWriter, r *http.Request)
	PatchGroup(w http.ResponseWriter, r *http.Request)
	AppendUsers(w http.ResponseWriter, r *http.Request)
	RemoveUsers(w http.ResponseWriter, r *http.Request)
	ListGroups(w http.ResponseWriter, r *http.Request)
}

type Service interface {
	UserService
	GroupService
}

func NewTestService(i do.Injector) (Service, error) {
	uR := do.MustInvoke[user.UserRepository](i)
	gR := do.MustInvoke[group.GroupRepository](i)
	qR := do.MustInvoke[quiz.QuizRepository](i)
	tR := do.MustInvoke[test.TestRepository](i)
	aR := do.MustInvoke[answer.AnswerRepository](i)

	r := do.MustInvoke[testrunner.TestRunner](i)

	return &chiService{
		userRepo:   uR,
		quizRepo:   qR,
		groupRepo:  gR,
		testRepo:   tR,
		answerRepo: aR,
		runner:     r,
	}, nil
}
