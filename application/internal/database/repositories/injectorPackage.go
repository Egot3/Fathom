package repositories

import (
	"github.com/egot3/fathom/internal/database/repositories/answer"
	"github.com/egot3/fathom/internal/database/repositories/group"
	"github.com/egot3/fathom/internal/database/repositories/quiz"
	"github.com/egot3/fathom/internal/database/repositories/test"
	"github.com/egot3/fathom/internal/database/repositories/user"
	"github.com/samber/do/v2"
)

var RepositoryPackage = do.Package(
	do.Lazy(user.NewUserRepository),
	do.Lazy(group.NewGroupRepository),
	do.Lazy(quiz.NewQuizRepository),
	do.Lazy(test.NewTestRepository),
	do.Lazy(answer.NewAnswerRepository),
)
