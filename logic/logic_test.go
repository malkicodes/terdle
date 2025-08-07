package logic_test

import (
	"strconv"
	"testing"

	"malki.codes/terdle/logic"
	"malki.codes/terdle/words"
)

var correctnessTestCases map[string]map[string]string = map[string]map[string]string{
	"aaaaa": {
		"bbbbb": "00000",
		"bbabb": "00200",
		"bbbaa": "00022",
		"aaaaa": "22222",
	},
	"greet": {
		"eexxx": "11000",
		"eexxe": "11000",
		"eeexx": "10200",
		"xeeex": "00220",
	},
	// just some hellowordl.net games i played lol
	"alive": {
		"salet": "01110",
		"learn": "11100",
		"alien": "22210",
		"alike": "22202",
		"alive": "22222",
	},
	"blank": {
		"salet": "01100",
		"alike": "12010",
		"clank": "02222",
		"plank": "02222",
		"blank": "22222",
	},
	"tenet": {
		"salet": "00022",
		"poker": "00020",
		"bingo": "00200",
		"tenet": "22222",
	},
	"torch": {
		"salet": "00001",
		"boner": "02001",
		"corky": "12200",
	},
	"knock": {
		"onion": "12000",
	},
	// Wordle 1,510 4/6
	"coral": {
		"salet": "01100",
		"plain": "01100",
		"comal": "22022",
		"coral": "22222",
	},
}

func TestValidateGuessCorrectness(t *testing.T) {
	for answerRaw, guesses := range correctnessTestCases {
		answer := words.StrToWord(answerRaw)

		for guess, expected := range guesses {
			correctness, _ := logic.ValidateGuess(answer, words.StrToWord(guess))

			for i, v := range correctness {
				value := int(v)
				expectedVal, err := strconv.Atoi(string(expected[i]))
				if err != nil {
					t.Fatal(err)
				}

				if value != expectedVal {
					actual := ""

					for _, v := range correctness {
						actual += strconv.Itoa(int(v))
					}

					t.Errorf("ValidateGuess(%s, %s) = %s, expected %s", answerRaw, guess, actual, expected)
					break
				}
			}
		}
	}
}

func TestValidateGuessWin(t *testing.T) {
	_, win := logic.ValidateGuess(words.StrToWord("abcde"), words.StrToWord("abcde"))

	if !win {
		t.Errorf("ValidateGuess(abcde, abcde) = %t, expected %t", win, true)
	}

	_, win = logic.ValidateGuess(words.StrToWord("abcde"), words.StrToWord("bcdea"))

	if win {
		t.Errorf("ValidateGuess(abcde, bcdea) = %t, expected %t", win, false)
	}

	_, win = logic.ValidateGuess(words.StrToWord("abcde"), words.StrToWord("bcdea"))

	if win {
		t.Errorf("ValidateGuess(abcde, bcdea) = %t, expected %t", win, false)
	}
}
