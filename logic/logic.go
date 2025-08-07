package logic

import (
	"slices"

	"malki.codes/terdle/words"
)

func ValidateGuess(answer, guess words.Word) ([5]uint8, bool) {
	var correctness [5]uint8

	// Find position of all letters in guesses
	var answerPositions = make(map[rune][]int)

	for i, letter := range answer {
		_, prs := answerPositions[letter]

		if !prs {
			answerPositions[letter] = make([]int, 0, 5)
		}

		answerPositions[letter] = append(answerPositions[letter], i)
	}

	// Pass for correct letters
	for i, letter := range guess {
		if slices.Contains(answerPositions[letter], i) {
			correctness[i] = 2
			answerPositions[letter] = answerPositions[letter][1:]
		}
	}

	// Pass for yellow letters
	for i, letter := range guess {
		if slices.Contains(answer[0:], letter) && len(answerPositions[letter]) != 0 && correctness[i] != 2 {
			correctness[i] = 1

			// At this stage, all that matters is the length of answerPositions[letter]
			answerPositions[letter] = answerPositions[letter][1:]
		}
	}

	// Win validation
	win := true

	for _, v := range correctness {
		if v != 2 {
			win = false
			break
		}
	}

	return correctness, win
}
