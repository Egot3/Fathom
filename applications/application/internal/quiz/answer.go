package quiz

type AnswerInput struct {
	Input string `json:"input"`
}

type AnswerRadio struct {
	ChoiceIdx int `json:"chosen"`
}

type AnswerCheck struct {
	ChoiceIdxs []int `json:"chosen"`
}

type AnswerAccordance struct {
	Accordance []int `json:"accorded"`
}

type AnswerOrder struct {
	ItemIdxs []int `json:"item_indexes"`
}

type QuizAnswers struct {
	Radio      AnswerRadio      `json:"radio"`
	Check      AnswerCheck      `json:"check"`
	Accordance AnswerAccordance `json:"accordance"`
	Order      AnswerOrder      `json:"order"`
	Input      AnswerInput      `json:"input"`
}
