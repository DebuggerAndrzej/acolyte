package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"

	"github.com/DebuggerAndrzej/acolyte/backend/entities"
)

type Model struct {
	stopwatch         stopwatch.Model
	spinner           spinner.Model
	selectedComponent int8
	procs             []entities.Proc
}
type tickMsg time.Time

func initModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7FBBB3"))
	return Model{
		stopwatch:         stopwatch.NewWithInterval(time.Second),
		spinner:           s,
		selectedComponent: 0,
		procs: []entities.Proc{
			{Name: "LS LA", Command: "ls -la"},
			{Name: "GREP", Command: "grep -c 5 usefullgrep"},
			{Name: "PING", Command: "ping -c 30 google.pl"},
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, m.stopwatch.Init(), m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit

		case "tab":
			m.selectedComponent += 1
		case "shift+tab":
			m.selectedComponent -= 1
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var s string
	var cards []string
	mainTitle := mainTitleStyle.Render("Acolyte dashboard")
	for i, proc := range m.procs {
		cards = append(cards, generateProcCard(proc.Name, m, i))
	}
	s += lipgloss.JoinVertical(lipgloss.Top, mainTitle, lipgloss.JoinHorizontal(lipgloss.Top, cards...))
	return s
}
func generateProcCard(title string, m Model, index int) string {
	cardStyleRenderer := modelStyle.Render
	if index == int(m.selectedComponent) {
		cardStyleRenderer = focusedModelStyle.Render
	}
	return cardStyleRenderer(
		lipgloss.JoinVertical(
			lipgloss.Top,
			procNameStyle.Render(title),
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				fmt.Sprintf("Status: %s", m.spinner.View()),
				"     ",
				fmt.Sprintf("Time: %s", m.stopwatch.View()),
			),
		),
	)
}

func InitTui() {
	model := initModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
