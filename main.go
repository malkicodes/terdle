package main

import (
	"flag"
	"log"
	"strconv"

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
			} else if m.guesses.Index() == -1 {
				m.endCase = -1
				return m, tea.Quit
			}
		default:
			m.input.Update(msg)
		}
	}

	return m, nil
}

var containerStyle = lipgloss.NewStyle().
	Width(11).Height(8).
	Padding(1, 3).
	Border(lipgloss.RoundedBorder())

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
		containerStyle.Render(rendered),
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
