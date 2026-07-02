package quiz

import "github.com/google/uuid"

type Frontmatter struct {
	Kind       string `yaml:"kind" json:"kind"`
	Randomized bool   `yaml:"randomized" json:"randomized"`
	Score      int    `yaml:"score" json:"score"`
	AllOrNone  bool   `yaml:"all-or-none" json:"all_or_none"`
}

type Quiz struct {
	Meta    Frontmatter `json:"meta"`
	Title   string      `json:"title"`
	Body    string      `json:"body"`
	UUID    uuid.UUID
	Options QuizOptions `json:"options"`
	Answer  QuizAnswers `json:"answers"`
}
