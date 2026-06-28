package quizparser

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/egot3/fathom/internal/quiz"
)

func ParseQuizByBytes(fileBytes []byte) (*quiz.Quiz, error) {

	var quiz quiz.Quiz
	fm, source, err := ParseFrontmatter(fileBytes)
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

		if InputRegex.MatchString(trimmedLine) {
			break
		}
		if strings.HasPrefix(trimmedLine, "- ") {
			break
		}

		quiz.Body += line + " "

	}

	log.Printf("quiz with: %v, and body: %v", reader.Text(), quiz.Body)
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
