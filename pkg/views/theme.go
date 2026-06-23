package views

import "github.com/charmbracelet/lipgloss"

var (
	PrimaryColor = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7D7AFF"}
	AccentColor  = lipgloss.AdaptiveColor{Light: "#008F72", Dark: "#5EE6C4"}
	TextColor    = lipgloss.AdaptiveColor{Light: "#202124", Dark: "#F2F2F2"}
	MutedColor   = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#8B8B99"}
	BorderColor  = lipgloss.AdaptiveColor{Light: "#D5D7E0", Dark: "#3A3A48"}
	ErrorColor   = lipgloss.AdaptiveColor{Light: "#C62828", Dark: "#FF6B6B"}
	SuccessColor = lipgloss.AdaptiveColor{Light: "#137A5B", Dark: "#5EE6A8"}

	AppTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor)

	PageTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextColor)

	LabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(MutedColor)

	BodyStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	MutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor)

	ActiveItemStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentColor)

	InactiveItemStyle = lipgloss.NewStyle().
				Foreground(TextColor)

	InputPromptStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(PrimaryColor)

	InputTextStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	InputPlaceholderStyle = lipgloss.NewStyle().
				Foreground(MutedColor)

	CursorStyle = lipgloss.NewStyle().
			Foreground(AccentColor)

	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextColor).
			Background(PrimaryColor)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(MutedColor)
)

func CardStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(1, 2).
		Width(max(1, width-6))
}

func PanelStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 1).
		Width(max(1, width-4))
}

func ComposerStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(0, 1).
		Width(max(1, width-4))
}

func KeyStyle(key string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		Render(key)
}
