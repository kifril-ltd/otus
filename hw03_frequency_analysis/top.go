package hw03frequencyanalysis

import (
	"strings"
	"unicode"
)

func Top10(in string) []string {
	analyser := NewFrequencyAnalyser()
	analyser.AddNormalizeFunc([]NormalizeFunc{
		func(s string) (string, error) {
			return strings.TrimFunc(s, func(r rune) bool {
				return !unicode.IsLetter(r)
			}), nil
		},
		func(s string) (string, error) {
			if len(s) == 0 {
				return "", nil
			}

			runes := []rune(s)
			runes[0] = unicode.ToLower(runes[0])
			runes[len(runes)-1] = unicode.ToLower(runes[len(runes)-1])

			return string(runes), nil
		},
	}...)

	analyser.Analyse(in)

	return analyser.Top(10)
}
