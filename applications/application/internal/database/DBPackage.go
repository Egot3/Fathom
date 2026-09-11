package database

import (
	"github.com/egot3/fathom/internal/database/repositories"
	"github.com/samber/do/v2"
)

var DBPackage = do.Package(
	do.Lazy(InitDB),
	repositories.RepositoryPackage,
)
