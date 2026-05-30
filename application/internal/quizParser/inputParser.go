package quizparser

import (
	"bufio"
	"fmt"
	"regexp"
)

var regex = regexp.MustCompile(`\[([^\[\]]+)\]`)

func InputParser(reader *bufio.Scanner, quiz *Quiz) error {
	for reader.Scan() {
		line := reader.Text()
		matches := regex.FindStringSubmatch(line)
		if len(matches) > 1 {
			quiz.Answer.Input = &AnswerInput{Input: matches[1]}
			break
		}
	}
	if quiz.Answer.Input == nil {
		return fmt.Errorf("Can't have no answer in input")
	}

	return nil
}
