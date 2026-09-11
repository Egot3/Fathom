package testrunner

import "github.com/samber/do/v2"

var ManagerPackage = do.Package(
	do.Lazy(NewManager),
)
