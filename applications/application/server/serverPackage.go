package server

import (
	"github.com/egot3/fathom/internal/handler"
	"github.com/samber/do/v2"
)

var ServerPackage = do.Package(
	do.Lazy(ChiServer),
	handler.HandlerPackage,
)
