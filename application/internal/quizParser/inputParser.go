package quizparser

import (
	"bufio"
	"fmt"
	"regexp"

	"github.com/egot3/fathom/internal/quiz"
)

var InputRegex = regexp.MustCompile(`\[([^\[\]]+)\]`)

func InputParser(reader *bufio.Scanner, quizP *quiz.Quiz) error {
	for reader.Scan() {
		line := reader.Text()
		matches := InputRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			quizP.Answer.Input = &quiz.AnswerInput{Input: matches[1]}
			break
		}
	}
	if quizP.Answer.Input == nil {
		return fmt.Errorf("Can't have no answer in input")
	}

	return nil
}
