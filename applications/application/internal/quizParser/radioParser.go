package quizparser

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/egot3/fathom/internal/quiz"
)

func RadioParser(reader *bufio.Scanner, quizP *quiz.Quiz) error {
	quizP.Options.Radio = quiz.OptionsRadioAndCheck{Choices: make([]quiz.Choice, 2)}
	quizP.Answer.Radio.ChoiceIdx = -1
	for id := 0; ; id++ {
		line := reader.Text()
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "- [x]") {
			if quizP.Answer.Radio.ChoiceIdx != -1 {
				return fmt.Errorf("radio can't have multiple answers")
			}

			opt := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "- [x] "))
			quizP.Answer.Radio.ChoiceIdx = id

			quizP.Options.Radio.Choices = append(quizP.Options.Radio.Choices, quiz.Choice{Id: id, Label: opt})
		} else if strings.HasPrefix(trimmedLine, "- [ ]") {
			opt := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "- [ ] "))

			quizP.Options.Radio.Choices = append(quizP.Options.Radio.Choices, quiz.Choice{Id: id, Label: opt})
		}
		if !reader.Scan() {
			break
		}
	}

	if len(quizP.Options.Radio.Choices) < 2 {
		return fmt.Errorf("radio can't have less than 2 options")
	}

	if quizP.Answer.Radio.ChoiceIdx == -1 {
		return fmt.Errorf("radio must have exactly 1 answer")
	}

	if quizP.Meta.Randomized {
		rand.Shuffle(len(quizP.Options.Radio.Choices), func(i, j int) {
			quizP.Options.Radio.Choices[i], quizP.Options.Radio.Choices[j] = quizP.Options.Radio.Choices[j], quizP.Options.Radio.Choices[i]
		})
	}

	return nil
}
