package main

import (
	"flag"
	"log"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"malki.codes/terdle/models"
	"malki.codes/terdle/words"
)

type MainModel struct {
	width    int
	height   int
	contrast bool

	answer  string
	endCase int8

	guesses models.GuessDisplayModel
	input   models.GuessInputModel
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			guess := m.input.Clear()

			if guess == "" {
				break
			}

			m.guesses.Update(models.NewGuessDisplayModelUpdateMsg(m.answer, guess))

			if m.guesses.Win {
				m.endCase = 1
				return m, tea.Quit
			} else if m.guesses.Index() == 6 {
				m.endCase = -1
				return m, tea.Quit
			}
		default:
			m.input.Update(msg)
		}
	}

	return m, nil
}

var gameContainerStyle = lipgloss.NewStyle().
	Width(11).Height(8).
	Padding(1, 3).
	Border(lipgloss.RoundedBorder())

var keyboardContainerStyle = lipgloss.NewStyle().
	Align(lipgloss.Center, lipgloss.Center)

func (m MainModel) View() string {
	rendered := ""

	if m.guesses.Index() == 0 {
		rendered = m.input.View()
	} else {
		rendered = lipgloss.JoinVertical(0, m.guesses.View(), m.input.View())
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			gameContainerStyle.Render(rendered),
			displayKeyboard(&m.guesses),
		),
	)
}

var keyboard [][]rune = [][]rune{
	{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p'},
	{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l'},
	{'z', 'x', 'c', 'v', 'b', 'n', 'm'},
}

func displayKeyboard(m *models.GuessDisplayModel) string {
	status := make(map[rune]uint8, 0)

	for i := 0; i < m.Index(); i++ {
		for j, char := range m.Guesses[i] {
			status[char] = m.Correctness[i][j]
		}
	}

	lines := make([]string, 3)

	for i, row := range keyboard {
		line := ""

		for _, char := range row {
			renderedChar := string(char)

			correctness, prs := status[char]

			var styles map[uint8]lipgloss.Style

			if m.Contrast {
				styles = models.ContrastCorrectnessStyle
			} else {
				styles = models.CorrectnessStyle
			}

			if prs {
				renderedChar = styles[correctness].Render(renderedChar)
			}

			line += " " + renderedChar + " "
		}

		lines[i] = strings.TrimSpace(line)
	}

	return keyboardContainerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			lines[0],
			lines[1],
			lines[2],
		),
	)
}

func displayEndscreen(m MainModel) string {
	renderedGuesses := make([]string, 0, 6)

	var styles map[uint8]lipgloss.Style

	if m.contrast {
		styles = models.ContrastCorrectnessStyle
	} else {
		styles = models.CorrectnessStyle
	}

	for i, guess := range m.guesses.Correctness {
		if i == m.guesses.Index() {
			break
		}

		rendered := ""

		for _, j := range guess {
			rendered = lipgloss.JoinHorizontal(0, rendered, styles[j].Render("██"))
		}

		renderedGuesses = append(renderedGuesses, rendered)
	}

	firstLine := "Terdle (" + m.answer + ") "

	switch m.endCase {
	case 0:
		firstLine += "Q"
	case 1:
		firstLine += strconv.Itoa(m.guesses.Index())
	case -1:
		firstLine += "X"
	}

	firstLine += "/6"

	return "\n" + lipgloss.JoinVertical(0, firstLine, lipgloss.JoinVertical(0, renderedGuesses...)) + "\n"
}

func main() {
	contrast := flag.Bool("c", false, "High contrast mode")
	answer := flag.String("a", words.GetRandomAnswer(), "Answer")

	flag.Parse()

	if len(*answer) != 5 {
		log.Fatalln("answer must be 5 letters long")
	}

	m := MainModel{
		answer:   *answer,
		contrast: *contrast,

		guesses: models.NewGuessDisplayModel(*contrast),
		input:   models.NewGuessInputModel(*contrast),
	}

	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}
	println(displayEndscreen(m))
}
