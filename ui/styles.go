package ui

import (
	lipgloss "github.com/charmbracelet/lipgloss"
)

var (
	procNameStyle = lipgloss.NewStyle().
			Width(45).
			Height(1).
			Bold(true).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.ThickBorder()).
			BorderBottom(true)
	modelStyle = lipgloss.NewStyle().
			Width(45).
			Height(4).
			BorderStyle(lipgloss.RoundedBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(45).
				Height(4).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("69"))
	mainTitleStyle = lipgloss.NewStyle().
			Height(1).
			Bold(true).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.ThickBorder()).
			BorderBottom(true)
	titleStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(0, 1)
)
