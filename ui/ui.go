package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type Model struct {
	stopwatch         stopwatch.Model
	spinner           spinner.Model
	selectedComponent int8
}
type tickMsg time.Time

func initModel() Model {
	return Model{stopwatch: stopwatch.NewWithInterval(time.Second), spinner: spinner.New(), selectedComponent: 0}
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
			if m.selectedComponent == 0 {
				m.selectedComponent += 1
			} else {
				m.selectedComponent -= 1
			}
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
	if m.selectedComponent == 0 {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(fmt.Sprintf("%4s", m.stopwatch.View())), modelStyle.Render(m.spinner.View()))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%4s", m.stopwatch.View())), focusedModelStyle.Render(m.spinner.View()))
	}
	return s
}

func InitTui() {
	model := initModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
