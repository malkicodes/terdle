package words

import (
	_ "embed"
	"math/rand/v2"
	"strings"
)

type Word = [5]rune

//go:embed answers.txt
var answersString string

//go:embed valid.txt
var validString string

var Answers []string
var Valid []string

func GetRandomAnswer() string {
	return Answers[rand.IntN(len(Answers))]
}

func StrToWord(s string) Word {
	var word Word

	for i := range 5 {
		word[i] = rune(s[i])
	}

	return word
}

func init() {
	for v := range strings.SplitSeq(answersString, "\n") {
		if len(v) != 0 {
			Answers = append(Answers, v)
		}
	}

	for v := range strings.SplitSeq(validString, "\n") {
		if len(v) != 0 {
			Valid = append(Valid, v)
		}
	}
}
