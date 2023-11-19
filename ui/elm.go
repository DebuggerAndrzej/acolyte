package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.view {
	case "dashboard":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.Right):
				m.selectedComponent += 1
			case key.Matches(msg, m.keys.Left):
				m.selectedComponent -= 1
			case key.Matches(msg, m.keys.Down):
				m.selectedComponent += m.cardsInRow
			case key.Matches(msg, m.keys.Up):
				m.selectedComponent -= m.cardsInRow
			case key.Matches(msg, m.keys.Start):
				return m, tea.Sequence(m.signalProcRunning, m.startProc)
			case key.Matches(msg, m.keys.ShowOutput):
				m.view = "command output"
				m.commandOutput.SetContent(m.procs[m.selectedComponent].output)
			}
		case spinner.TickMsg:
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		case tea.WindowSizeMsg:
			m.commandOutput.Width = msg.Width
			m.commandOutput.Height = msg.Height - 6
			m.cardsInRow = int(msg.Width / lipgloss.Width(modelStyle.Render()))
		}
		for i := 0; i < len(m.procs); i++ {
			uiProc := &m.procs[i]
			if uiProc.isRunning == true {
				uiProc.stopwatch, cmd = uiProc.stopwatch.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
		m.validateSelectedComponent()
	case "command output":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.Quit):
				m.view = "dashboard"
			}
		case spinner.TickMsg:
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
		m.commandOutput, cmd = m.commandOutput.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.view == "dashboard" {
		var s string
		var cards []string
		mainTitle := mainTitleStyle.Width(m.commandOutput.Width).Render("Acolyte dashboard")
		for i, uiProc := range m.procs {
			cards = append(cards, m.generateProcCard(uiProc, i))
		}
		var cardRowsFormatted []string
		m.cardsInRow = max(0, int(m.commandOutput.Width/lipgloss.Width(cards[0])))
		for i := 0; i < len(cards); i += m.cardsInRow {
			end := i + m.cardsInRow

			if end > len(cards) {
				end = len(cards)
			}

			cardRowsFormatted = append(cardRowsFormatted, lipgloss.JoinHorizontal(lipgloss.Top, cards[i:end]...))
		}
		cardsFormatted := lipgloss.JoinVertical(lipgloss.Top, cardRowsFormatted...)
		s += lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.NewStyle().
				Height(m.commandOutput.Height+5).
				Render(lipgloss.JoinVertical(lipgloss.Top, mainTitle, cardsFormatted)),
			m.help.View(m.keys),
		)

		return s
	}

	if m.view == "command output" {
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.commandOutput.View(), m.footerView())
	}
	return "How did we get here?"

}
