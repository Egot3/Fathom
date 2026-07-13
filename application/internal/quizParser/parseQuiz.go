package quizparser

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/egot3/fathom/internal/quiz"
)

func ParseQuizByBytes(fileBytes []byte) (*quiz.Quiz, error) {

	var q quiz.Quiz
	fm, source, err := ParseFrontmatter(fileBytes)
	if err != nil {
		return nil, err
	}

	q.Meta = fm

	reader := bufio.NewScanner(bytes.NewReader(source))

	for reader.Scan() {
		line := reader.Text()
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "# ") {
			q.Title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			break
		}
	}

	for reader.Scan() {
		line := reader.Text()
		trimmedLine := strings.TrimSpace(line)

		if InputRegex.MatchString(trimmedLine) {
			break
		}
		if strings.HasPrefix(trimmedLine, "- ") {
			break
		}

		q.Body += line + " "

	}
	q.Body = strings.TrimSpace(q.Body)

	// each quiz has 1 typeof question
	switch q.Meta.Kind {
	case quiz.Input:
		if err := InputParser(reader, &q); err != nil {
			return nil, err
		}
	case quiz.Check:
		if err := CheckParser(reader, &q); err != nil {
			return nil, err
		}
	case quiz.Radio:
		if err := RadioParser(reader, &q); err != nil {
			return nil, err
		}
	case quiz.Order:
		if err := OrderParser(reader, &q); err != nil {
			return nil, err
		}
	case quiz.Accordance:
		if err := AccordanceParser(reader, &q); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unsupported kind")
	}

	return &q, nil
}

// P.S. on parts where I am putting comments, my brain melts
