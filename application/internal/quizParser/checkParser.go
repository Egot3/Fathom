package quizparser

import (
	"bufio"
	"fmt"
	"log"
	"math/rand/v2"
	"strings"
)

func CheckParser(reader *bufio.Scanner, quiz *Quiz) error {
	quiz.Options.Check = &OptionsRadioAndCheck{Choices: make([]Choice, 0)}
	quiz.Answer.Check = &AnswerCheck{ChoiceIdxs: make([]int, 0)}
	for id := 0; reader.Scan(); {
		line := reader.Text()
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "- [x]") {
			opt := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "- [x] "))
			quiz.Answer.Check.ChoiceIdxs = append(quiz.Answer.Check.ChoiceIdxs, id)

			quiz.Options.Check.Choices = append(quiz.Options.Check.Choices, Choice{Id: id, Label: opt})
			id++
			continue
		}
		if strings.HasPrefix(trimmedLine, "- [ ]") {
			opt := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "- [ ] "))

			quiz.Options.Check.Choices = append(quiz.Options.Check.Choices, Choice{Id: id, Label: opt})
			id++
			continue
		}
	}

	log.Printf("%v with %v", len(quiz.Options.Check.Choices), len(quiz.Answer.Check.ChoiceIdxs))

	if len(quiz.Options.Check.Choices) < 2 {
		return fmt.Errorf("can't have less than 2 options for check")
	}
	if len(quiz.Answer.Check.ChoiceIdxs) == 0 {
		return fmt.Errorf("can't have less than 1 answer for check")
	}

	if quiz.Meta.Randomized {
		rand.Shuffle(len(quiz.Options.Check.Choices), func(i, j int) {
			quiz.Options.Check.Choices[i], quiz.Options.Check.Choices[j] = quiz.Options.Check.Choices[j], quiz.Options.Check.Choices[i]
		})
	}

	return nil
}
