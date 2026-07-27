package quiz

type OptionsRadioAndCheck struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
}

type OptionsAccordance struct {
	Static  []string `json:"static"`
	Dynamic []string `json:"dynamic"`
}

type OptionsOrder struct {
	Items []string `json:"items"`
}

type QuizOptions struct {
	Radio      OptionsRadioAndCheck `json:"radio,omitzero"`
	Check      OptionsRadioAndCheck `json:"check,omitzero"`
	Accordance OptionsAccordance    `json:"accordance,omitzero"`
	Order      OptionsOrder         `json:"orders,omitzero"`
}
