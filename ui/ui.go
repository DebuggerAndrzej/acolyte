package ui

import (
	"fmt"
	"log"
	"os"
	"strings"
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
	proc      entities.Proc
	isRunning bool
	timeRun   string
	output    string
	stopwatch stopwatch.Model
}

type Model struct {
	spinner           spinner.Model
	selectedComponent int8
	procs             []UiProc
	commandOutput     viewport.Model
	runningProcs      []string
	view              string
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
				timeRun:   "Hasn't been run yet",
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "GREP", Command: "grep -c 5 usefullGrep"},
				isRunning: false,
				timeRun:   "Hasn't been run yet",
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "PING", Command: "ping -c 150 google.pl"},
				isRunning: false,
				timeRun:   "Hasn't been run yet",
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "LS LA 2", Command: "ls -la"},
				isRunning: false,
				timeRun:   "Hasn't been run yet",
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "GREP 2", Command: "grep -c 5 usefullGrep"},
				isRunning: false,
				timeRun:   "Hasn't been run yet",
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "PING 2", Command: "ping -c 150 google.pl"},
				isRunning: false,
				timeRun:   "Hasn't been run yet",
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "LS LA 3", Command: "ls -la"},
				isRunning: false,
				timeRun:   "Hasn't been run yet",
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "GREP 3", Command: "grep -c 5 usefullGrep"},
				isRunning: false,
				timeRun:   "Hasn't been run yet",
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
			{
				proc:      entities.Proc{Name: "PING 3", Command: "ping -c 150 google.pl"},
				isRunning: false,
				timeRun:   "Hasn't been run yet",
				output:    "",
				stopwatch: stopwatch.NewWithInterval(time.Second),
			},
		},
		commandOutput: viewport.New(200, 47),
		view:          "dashboard",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.view == "dashboard" {
				return m, tea.Quit
			} else {
				m.view = "dashboard"
			}
		case "tab":
			m.selectedComponent += 1
		case "shift+tab":
			m.selectedComponent -= 1
		case "s":
			return m, tea.Sequence(m.signalProcRunning, m.startProc)
		case "enter":
			m.view = "command output"
			m.commandOutput.SetContent(m.procs[m.selectedComponent].output)
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		m.commandOutput.Width = msg.Width
		m.commandOutput.Height = msg.Height - 6
	}
	for i := 0; i < len(m.procs); i++ {
		uiProc := &m.procs[i]
		if uiProc.isRunning == true {
			uiProc.stopwatch, cmd = uiProc.stopwatch.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	if m.view == "command output" {
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
			cards = append(cards, generateProcCard(uiProc, m, i))
		}
		var cardRowsFormatted []string
		cardsInRow := max(0, int(m.commandOutput.Width/lipgloss.Width(cards[0])))
		for i := 0; i < len(cards); i += cardsInRow {
			end := i + cardsInRow

			if end > len(cards) {
				end = len(cards)
			}

			cardRowsFormatted = append(cardRowsFormatted, lipgloss.JoinHorizontal(lipgloss.Top, cards[i:end]...))
		}
		cardsFormatted := lipgloss.JoinVertical(lipgloss.Top, cardRowsFormatted...)
		s += lipgloss.JoinVertical(lipgloss.Top, mainTitle, cardsFormatted)
		return s
	}

	if m.view == "command output" {
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.commandOutput.View(), m.footerView())
	}
	return "How did we get here?"

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

func generateProcCard(uiProc UiProc, m Model, index int) string {
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

func InitTui() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	model := initModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
