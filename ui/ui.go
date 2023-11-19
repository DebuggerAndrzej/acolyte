package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"

	"github.com/DebuggerAndrzej/acolyte/backend"
	"github.com/DebuggerAndrzej/acolyte/backend/entities"
)

type UiProc struct {
	proc entities.Proc
	isRunning bool
	timeRun string
	output string
}

type Model struct {
	stopwatch         stopwatch.Model
	spinner           spinner.Model
	selectedComponent int8
	procs             []UiProc
	commandOutput viewport.Model
	runningProcs []string
	view string
}

func initModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7FBBB3"))
	return Model{
		stopwatch:         stopwatch.NewWithInterval(time.Second),
		spinner:           s,
		selectedComponent: 0,
		procs: []UiProc{
			{proc: entities.Proc{Name: "LS LA", Command: "ls -la"}, isRunning: false, timeRun: "Hasn't been run yet", output: ""},
			{proc: entities.Proc{Name: "GREP", Command: "grep -c 5 usefullGrep"}, isRunning: false, timeRun: "Hasn't been run yet", output: ""},
			{proc: entities.Proc{Name: "PING", Command: "ping -c 20 google.pl"}, isRunning: false, timeRun: "Hasn't been run yet", output: ""},
		},
		commandOutput: viewport.New(200, 100),
		view: "dashboard",
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
		case "s":
			return m, tea.Sequence(m.signalProcRunning, m.startProc)
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
	for i, uiProc := range m.procs {
		cards = append(cards, generateProcCard(uiProc, m, i))
	}
	s += lipgloss.JoinVertical(lipgloss.Top, mainTitle, lipgloss.JoinHorizontal(lipgloss.Top, cards...))
	return s
}
func generateProcCard(uiProc UiProc, m Model, index int) string {
	cardStyleRenderer := modelStyle.Render
	status := "  "
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
				fmt.Sprintf("Time: %s", "0s"),
			),
		),
	)
}

func (m *Model) signalProcRunning() tea.Msg {
	m.procs[m.selectedComponent].isRunning = true
	return nil

}

func (m *Model) startProc() tea.Msg {
	m.procs[m.selectedComponent].output  =backend.DummyRunCommand(m.procs[m.selectedComponent].proc.Command)
	m.procs[m.selectedComponent].isRunning = false
	return nil
}


func InitTui() {
	model := initModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
