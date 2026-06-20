package database

import (
	"github.com/samber/do/v2"
)

var DBPackage = do.Package(
	do.Lazy(InitDB),
)
