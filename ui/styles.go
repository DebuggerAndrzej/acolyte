package ui

import (
	lipgloss "github.com/charmbracelet/lipgloss"
)

var (
	procNameStyle = lipgloss.NewStyle().
			Width(30).
			Height(1).
			Bold(true).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.ThickBorder()).
			BorderBottom(true)
	modelStyle = lipgloss.NewStyle().
			Width(30).
			Height(3).
			BorderStyle(lipgloss.RoundedBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(30).
				Height(3).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	mainTitleStyle = lipgloss.NewStyle().
			Width(200).
			Height(1).
			Bold(true).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.ThickBorder()).
			BorderBottom(true)
)
