package quizparser

import (
	"bufio"
	"bytes"
	"fmt"
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
	Body    string
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

	for reader.Scan() {
		line := reader.Text()
		trimmedLine := strings.TrimSpace(line)

		if regex.MatchString(trimmedLine) {
			break
		}
		if strings.HasPrefix(trimmedLine, "- ") {
			break
		}

		quiz.Body += line + " "
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
		if err := RadioParser(reader, &quiz); err != nil {
			return nil, err
		}
	case "order":
		if err := OrderParser(reader, &quiz); err != nil {
			return nil, err
		}
	case "accordance":
		if err := AccordanceParser(reader, &quiz); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unsupported kind")
	}

	return &quiz, nil
}

// P.S. on parts where I am putting comments, my brain melts
