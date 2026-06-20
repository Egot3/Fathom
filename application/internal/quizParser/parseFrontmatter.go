package quizparser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/egot3/fathom/internal/quiz"
	"go.yaml.in/yaml/v4"
)

func ParseFrontmatter(source []byte) (quiz.Frontmatter, []byte, error) {
	var fm quiz.Frontmatter

	lines := bytes.Split(source, []byte("\n"))
	if len(lines) < 2 || strings.TrimSpace(string(lines[0])) != "---" {
		return fm, source, fmt.Errorf("missing frontmatter delimeter(---)")
	}

	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(string(lines[i])) == "---" {
			end = i
			break
		}
	}
	if end == -1 {
		return fm, source, fmt.Errorf("unclosed frontmatter")
	}

	if err := yaml.Unmarshal(bytes.Join(lines[1:end], []byte("\n")), &fm); err != nil {
		return fm, nil, err
	}
	if fm.Kind == "" {
		return fm, nil, fmt.Errorf("mising entry in frontmatter: kind")
	}
	if fm.Score == 0 {
		return fm, nil, fmt.Errorf("missing score/set to zero in frontmatter")
	}

	return fm, bytes.Join(lines[end+1:], []byte("\n")), nil
}
