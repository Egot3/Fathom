package quiz

import (
	"github.com/google/uuid"
)

type Kind string

const (
	Input      Kind = "INPUT"
	Radio      Kind = "RADIO"
	Check      Kind = "CHECK"
	Accordance Kind = "ACCORDANCE"
	Order      Kind = "ORDER"
)

type Frontmatter struct {
	Kind       Kind `yaml:"kind" json:"kind"`
	Randomized bool `yaml:"randomized" json:"randomized"`
	Score      int  `yaml:"score" json:"score"`
	AllOrNone  bool `yaml:"all-or-none" json:"all_or_none"`
}

type Quiz struct {
	Meta    Frontmatter `json:"meta"`
	Title   string      `json:"title"`
	Body    string      `json:"body"`
	UUID    uuid.UUID
	Options QuizOptions `json:"options"`
	Answer  QuizAnswers `json:"answers"`
}

func (q Quiz) EvaluateScore(answer QuizAnswers) float32 {
	switch q.Meta.Kind {
	case Input:
		return evaluateInput(q.Answer.Input, answer.Input, q.Meta.Score)
	case Radio:
		return evaluateRadio(q.Answer.Radio, answer.Radio, q.Meta.Score)
	case Order:
		if q.Meta.AllOrNone {
			return evaluateAllOrNoneOrder(q.Answer.Order, answer.Order, q.Meta.Score)
		}
		return evaluateOrder(q.Answer.Order, answer.Order, q.Meta.Score)
	case Check:
		if q.Meta.AllOrNone {
			return evaluateAllOrNoneCheck(q.Answer.Check, answer.Check, q.Meta.Score)
		}
		return evaluateCheck(q.Answer.Check, answer.Check, q.Meta.Score)
	case Accordance:
		if q.Meta.AllOrNone {
			return evaluateAllOrNoneAccordance(q.Answer.Accordance, answer.Accordance, q.Meta.Score)
		}
		return evaluateAccordance(q.Answer.Accordance, answer.Accordance, q.Meta.Score)
	}

	return 0
}
