package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Left       key.Binding
	Right      key.Binding
	Up         key.Binding
	Down       key.Binding
	Start      key.Binding
	ShowOutput key.Binding
	Quit       key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Start,
		k.ShowOutput,
		k.Quit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{}, {}}
}

var keys = keyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "shift+tab"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "tab"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
	),
	Start: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Start proc"),
	),
	ShowOutput: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Show proc output"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
