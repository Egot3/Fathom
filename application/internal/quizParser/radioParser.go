package quizparser

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"strings"
)

func RadioParser(reader *bufio.Scanner, quiz *Quiz) error {
	quiz.Options.Radio = &OptionsRadioAndCheck{Choices: make([]Choice, 2)}
	for id := 0; reader.Scan(); {
		line := reader.Text()
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "- [x]") {
			if quiz.Answer.Radio != nil {
				return fmt.Errorf("radio can't have multiple answers")
			}

			opt := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "- [x] "))
			quiz.Answer.Radio = &AnswerRadio{ChoiceIdx: id}

			quiz.Options.Radio.Choices = append(quiz.Options.Radio.Choices, Choice{Id: id, Label: opt})
			id++
			continue
		}
		if strings.HasPrefix(trimmedLine, "- [ ]") {
			opt := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "- [ ] "))

			quiz.Options.Radio.Choices = append(quiz.Options.Radio.Choices, Choice{Id: id, Label: opt})
			id++
			continue
		}
	}

	if len(quiz.Options.Radio.Choices) < 2 {
		return fmt.Errorf("radio can't have less than 2 options")
	}

	if quiz.Answer.Radio == nil {
		return fmt.Errorf("radio must have exactly 1 answer")
	}

	if quiz.Meta.Randomized {
		rand.Shuffle(len(quiz.Options.Radio.Choices), func(i, j int) {
			quiz.Options.Radio.Choices[i], quiz.Options.Radio.Choices[j] = quiz.Options.Radio.Choices[j], quiz.Options.Radio.Choices[i]
		})
	}

	return nil
}
