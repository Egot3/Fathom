package handler

import "github.com/samber/do/v2"

var HandlerPackage = do.Package(
	do.Lazy(NewTestService),
)
