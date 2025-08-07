package models

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"malki.codes/terdle/logic"
	"malki.codes/terdle/words"
)

type GuessDisplayModel struct {
	contrast bool

	guesses     [6]words.Word
	Correctness [6][5]uint8

	Win bool
}

var CORRECTNESS_0 = lipgloss.NewStyle().Inline(true).
	Faint(true).
	Width(1).Height(1)
var CORRECTNESS_1 = lipgloss.NewStyle().Inline(true).
	Bold(true).
	Foreground(lipgloss.ANSIColor(3)).
	Width(1).Height(1)
var CORRECTNESS_2 = lipgloss.NewStyle().Inline(true).
	Bold(true).
	Foreground(lipgloss.ANSIColor(2)).
	Width(1).Height(1)
var CorrectnessStyle = map[uint8]lipgloss.Style{
	0: CORRECTNESS_0,
	1: CORRECTNESS_1,
	2: CORRECTNESS_2,
}

var CONTRAST_CORRECTNESS_0 = lipgloss.NewStyle().Inline(true).
	Faint(true).
	Width(1).Height(1)
var CONTRAST_CORRECTNESS_1 = lipgloss.NewStyle().Inline(true).
	Bold(true).
	Foreground(lipgloss.ANSIColor(4)).
	Width(1).Height(1)
var CONTRAST_CORRECTNESS_2 = lipgloss.NewStyle().Inline(true).
	Bold(true).
	Foreground(lipgloss.ANSIColor(1)).
	Width(1).Height(1)
var ContrastCorrectnessStyle = map[uint8]lipgloss.Style{
	0: CONTRAST_CORRECTNESS_0,
	1: CONTRAST_CORRECTNESS_1,
	2: CONTRAST_CORRECTNESS_2,
}

func (m GuessDisplayModel) Index() int {
	for i, guess := range m.guesses {
		if guess[0] == 0 {
			return i
		}
	}
	return -1
}

func (m GuessDisplayModel) Init() tea.Cmd {
	return nil
}

func (m *GuessDisplayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case GuessDisplayModelUpdateMsg:
		m.Correctness[m.Index()] = msg.correctness
		m.guesses[m.Index()] = msg.guess
		m.Win = msg.win
	}

	return m, nil
}

func (m GuessDisplayModel) View() string {
	renderedGuesses := make([]string, 0, 6)

	var styles map[uint8]lipgloss.Style

	if m.contrast {
		styles = ContrastCorrectnessStyle
	} else {
		styles = CorrectnessStyle
	}

	for i, guess := range m.guesses {
		guessCorrectness := m.Correctness[i]
		rendered := ""

		if guess[0] == 0 {
			break
		}

		for j, letter := range guess {
			rendered = lipgloss.JoinHorizontal(0, rendered, styles[guessCorrectness[j]].Render(string(letter)))
		}

		renderedGuesses = append(renderedGuesses, rendered)
	}

	return lipgloss.JoinVertical(0, renderedGuesses...)
}

func NewGuessDisplayModel(contrast bool) GuessDisplayModel {
	return GuessDisplayModel{
		contrast: contrast,
	}
}

type GuessDisplayModelUpdateMsg struct {
	guess       words.Word
	correctness [5]uint8

	win bool
}

func NewGuessDisplayModelUpdateMsg(answer, guess string) GuessDisplayModelUpdateMsg {
	parsedGuess := words.StrToWord(guess)
	correctness, win := logic.ValidateGuess(words.StrToWord(answer), parsedGuess)

	return GuessDisplayModelUpdateMsg{
		guess:       parsedGuess,
		correctness: correctness,

		win: win,
	}
}

var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1))
var contrastErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(3))
var cursor = lipgloss.NewStyle().Underline(true).Render(" ")

type GuessInputModel struct {
	current  string
	hasError bool
	contrast bool
}

func (m GuessInputModel) Init() tea.Cmd {
	return nil
}

func (m *GuessInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		letter := msg.String()

		if (letter == "backspace" || letter == "delete") && len(m.current) > 0 {
			m.current = m.current[:len(m.current)-1]
			m.hasError = false
		} else if len(strings.TrimSpace(letter)) == 1 && len(m.current) < 5 {
			m.current += letter
			m.hasError = false
		}
	}

	return m, nil
}

func (m GuessInputModel) View() string {
	rendered := m.current

	if len(m.current) < 5 {
		rendered = lipgloss.JoinHorizontal(0, rendered, cursor)
	}

	if m.hasError {
		if m.contrast {
			rendered = contrastErrorStyle.Render(rendered)
		} else {
			rendered = errorStyle.Render(rendered)
		}
	}

	return rendered
}

func (m *GuessInputModel) Clear() string {
	word := m.current
	m.hasError = true

	if len(word) != 5 {
		return ""
	}

	if !slices.Contains(words.Valid, word) {
		return ""
	}

	m.current = ""
	m.hasError = false
	return word
}

func NewGuessInputModel(contrast bool) GuessInputModel {
	return GuessInputModel{
		current:  "",
		contrast: contrast,
	}
}
