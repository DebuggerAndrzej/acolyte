package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"

	"github.com/DebuggerAndrzej/acolyte/backend"
	"github.com/DebuggerAndrzej/acolyte/backend/entities"
)

type UiProc struct {
	proc      entities.Proc
	isRunning bool
	output    string
	stopwatch stopwatch.Model
}

type Model struct {
	spinner           spinner.Model
	selectedComponent int
	procs             []UiProc
	commandOutput     viewport.Model
	runningProcs      []string
	view              string
	keys              keyMap
	help              help.Model
	cardsInRow        int
}

func initModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7FBBB3"))
	return Model{
		spinner:           s,
		selectedComponent: 0,
		procs: []UiProc{
			{
				proc:      entities.Proc{Name: "LS LA", Command: "ls -la"},
				isRunning: false,
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "GREP", Command: "grep -c 5 usefullGrep"},
				isRunning: false,
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "PING", Command: "ping -c 150 google.pl"},
				isRunning: false,
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "LS LA 2", Command: "ls -la"},
				isRunning: false,
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "GREP 2", Command: "grep -c 5 usefullGrep"},
				isRunning: false,
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "PING 2", Command: "ping -c 40 google.pl"},
				isRunning: false,
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "LS LA 3", Command: "ls -la"},
				isRunning: false,
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "GREP 3", Command: "grep -c 5 usefullGrep"},
				isRunning: false,
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "PING 3", Command: "ping -c 70 google.pl"},
				isRunning: false,
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
		},
		commandOutput: viewport.New(200, 47),
		view:          "dashboard",
		keys:          keys,
		help:          help.New(),
		cardsInRow:    0,
	}
}

func (m Model) headerView() string {
	title := titleStyle.Render(m.procs[m.selectedComponent].proc.Name)
	line := strings.Repeat("─", max(0, m.commandOutput.Width-lipgloss.Width(title))/2)
	return lipgloss.JoinHorizontal(lipgloss.Center, line, title, line)
}
func (m Model) footerView() string {
	info := titleStyle.Render(fmt.Sprintf("%3.f%%", m.commandOutput.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.commandOutput.Width-lipgloss.Width(info))/2)
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info, line)
}

func (m *Model) signalProcRunning() tea.Msg {
	uiProc := &m.procs[m.selectedComponent]
	uiProc.isRunning = true
	return uiProc.stopwatch.Init()
}

func (m *Model) startProc() tea.Msg {
	m.procs[m.selectedComponent].output = backend.DummyRunCommand(m.procs[m.selectedComponent].proc.Command)
	m.procs[m.selectedComponent].isRunning = false
	return nil
}

func (m *Model) generateProcCard(uiProc UiProc, index int) string {
	cardStyleRenderer := modelStyle.Render
	status := ""
	if index == int(m.selectedComponent) {
		cardStyleRenderer = focusedModelStyle.Render
	}
	if uiProc.isRunning == true {
		status = m.spinner.View()
	}
	return cardStyleRenderer(
		lipgloss.JoinVertical(
			lipgloss.Top,
			procNameStyle.Render(uiProc.proc.Name),
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				fmt.Sprintf("Status: %s", status),
				"     ",
				fmt.Sprintf("Time: %s", uiProc.stopwatch.View()),
			),
		),
	)
}

func (m *Model) validateSelectedComponent() {
	amountOfProcs := len(m.procs)
	totalPossibleProcs := int(math.Ceil(float64(amountOfProcs)/float64(m.cardsInRow))) * m.cardsInRow
	switch {
	case m.selectedComponent == amountOfProcs:
		m.selectedComponent = 0
	case m.selectedComponent == -1:
		m.selectedComponent = amountOfProcs - 1
	case m.selectedComponent > amountOfProcs:
		m.selectedComponent -= totalPossibleProcs
	case m.selectedComponent < 0:
		m.selectedComponent += totalPossibleProcs
	}

}
