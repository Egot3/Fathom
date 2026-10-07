package exportutils

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Kind string

const (
	Test Kind = "TEST"
	Quiz Kind = "QUIZ"
	Bank Kind = "BANK"
)

type Manifest struct {
	Kind    Kind       `yaml:"kind"`
	Quizzes []YamlQuiz `yaml:"quizzes"`

	// for tests' manifest
	UUID uuid.UUID `yaml:"uuid,omitempty"`
	Name string    `yaml:"name,omitempty"`
}

func (m Manifest) Validate() error {
	switch m.Kind {
	case Bank:
		if m.Name != "" || m.UUID != uuid.Nil {
			return errors.New("bank manifest must not carry test fields")
		}
	case Test:
		if m.Name == "" || m.UUID == uuid.Nil {
			return errors.New("test manifest requires uuid and name")
		}
	default:
		return fmt.Errorf("unknown manifest kind %q", m.Kind)
	}
	return nil
}

type YamlQuiz struct {
	UUID     uuid.UUID `yaml:"uuid,omitempty"`
	Path     string    `yaml:"path"`
	Checksum string    `yaml:"checksum"`
}
