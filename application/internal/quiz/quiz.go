package quiz

type Frontmatter struct {
	Kind       string `yaml:"kind"`
	Randomized bool   `yaml:"randomized"`
	Score      int    `yaml:"score"`
	AllOrNone  bool   `yaml:"all-or-none"`
}

type Quiz struct {
	Meta    Frontmatter
	Title   string
	Body    string
	Options QuizOptions
	Answer  QuizAnswers
}
