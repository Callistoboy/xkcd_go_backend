package words

import (
	"maps"
	"slices"
	"strings"
	"unicode"

	snowball "github.com/kljensen/snowball"
	"github.com/kljensen/snowball/english"
)

func cleanText(c rune) bool {
	return !(unicode.IsLetter(c) || unicode.IsNumber(c))
}

func Norm(phrase string) []string {
	phrases := strings.FieldsFunc(phrase, cleanText)
	stemmed := make(map[string]bool, len(phrases))
	for _, word := range phrases {

		// empty check
		if word == "" {
			continue
		}

		word = strings.ToLower(word)

		if english.IsStopWord(word) {
			continue
		}

		// stemming
		stem, err := snowball.Stem(word, "english", true)
		if err != nil {
			continue
		}

		// empty after stemming check
		if stem == "" {
			continue
		}

		stemmed[stem] = true
	}
	return slices.Collect(maps.Keys(stemmed))
}
