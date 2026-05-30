package quizparser

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

type Frontmatter struct {
	Kind       string `yaml:"kind"`
	Randomized bool   `yaml:"randomized"`
	Score      int    `yaml:"score"`
	AllOrNone  bool   `yaml:"all-or-none"`
}

type Quiz struct {
	Meta    Frontmatter
	Title   string
	Options QuizOptions
	Answer  QuizAnswers
}

func ParseQuizByPath(path string) (*Quiz, error) {
	sourceFull, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var quiz Quiz
	fm, source, err := parseFrontmatter(sourceFull)
	if err != nil {
		return nil, err
	}

	quiz.Meta = fm

	reader := bufio.NewScanner(bytes.NewReader(source))

	for reader.Scan() {
		line := reader.Text()
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "# ") {
			quiz.Title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			break
		}
	}

	// each quiz has 1 typeof question
	switch quiz.Meta.Kind {
	case "input":
		if err := InputParser(reader, &quiz); err != nil {
			return nil, err
		}
	case "check":
		if err := CheckParser(reader, &quiz); err != nil {
			return nil, err
		}
	case "radio":
		RadioParser(reader, &quiz)
	case "order":
		OrderParser(reader, &quiz)
	case "accordance":
		AccordanceParser(reader, &quiz)
	}

	return &quiz, nil
}

// P.S. on parts where I am putting comments, my brain melts
