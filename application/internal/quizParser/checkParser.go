package quizparser

import (
	"bufio"
	"math/rand/v2"
	"strings"
)

func CheckParser(reader *bufio.Scanner, quiz *Quiz) error {
	for id := 0; reader.Scan(); {
		line := reader.Text()
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "- [x]") {
			opt := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "- [x] "))
			quiz.Answer.Check.Choices = append(quiz.Answer.Check.Choices, Choice{Id: id, Label: opt})

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
	if quiz.Meta.Randomized {
		rand.Shuffle(len(quiz.Options.Check.Choices), func(i, j int) {
			quiz.Options.Check.Choices[i], quiz.Options.Check.Choices[j] = quiz.Options.Check.Choices[j], quiz.Options.Check.Choices[i]
		})
	}

	return nil
}
