package exportutlis

import "github.com/google/uuid"

type Kind string

const (
	Test string = "TEST"
	Quiz string = "QUIZ"
)

type YamlTest struct {
	Kind        Kind       `yaml:"kind"`
	UUID        uuid.UUID  `yaml:"uuid"`
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Quizzes     []YamlQuiz `yaml:"quizzes"`
}

type YamlQuiz struct {
	Kind Kind      `yaml:"kind"`
	UUID uuid.UUID `yaml:"uuid,omitzero"`
}
