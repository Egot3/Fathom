package findingutils

import "strings"

func FuzzyMatch(pattern, target string) bool {
	pattern = strings.ToLower(pattern)
	target = strings.ToLower(target)

	patternRunes := []rune(pattern)
	targetRunes := []rune(target)

	if len(patternRunes) == 0 {
		return true
	}
	if len(patternRunes) > len(targetRunes) {
		return false
	}

	patternIdx := 0
	for _, targetRune := range targetRunes {
		if targetRune == patternRunes[patternIdx] {
			patternIdx++
			if patternIdx == len(patternRunes) {
				return true
			}
		}
	}

	return false
}
