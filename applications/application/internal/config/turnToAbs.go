package config

import (
	"path/filepath"
)

var pathToQuizzes string

func TurnToAbs(name string) (string, error) {
	return filepath.Abs(pathToQuizzes + name + ".md")
}
