package spinner

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var QuitMsg = errors.New("")

type errMsg error

type model struct {
	spinner   spinner.Model
	staticMsg string
	msg       string
	quit      bool
}

func NewSpinner(staticMsg string) model {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	if staticMsg == "" {
		staticMsg = "Procesando"
	}

	return model{spinner: s, staticMsg: staticMsg}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quit = true

			return m, tea.Quit
		default:
			m.msg = "to exit press ctrl+c"

			return m, nil
		}
	case errMsg:
		m.msg = msg.Error()
		m.quit = true

		return m, tea.Quit
	default:
		s, cmd := m.spinner.Update(msg)
		m.spinner = s

		return m, cmd
	}
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%s %s...", m.spinner.View(), m.staticMsg))

	if m.msg != "" {
		s.WriteString(fmt.Sprintf("\n%s", m.msg))
	}

	return s.String()
}
