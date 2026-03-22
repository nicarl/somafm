package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	playingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	detailStyle   = lipgloss.NewStyle().Padding(1, 2)
	statusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	selectedPlayingStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("82"))
	borderStyle          = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	focusedBorderStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("205"))
	errorStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

func (m model) View() string {
	if m.state == stateLoading {
		return fmt.Sprintf("\n  %s Loading channels...\n", m.spinner.View())
	}

	if m.state == stateError {
		return fmt.Sprintf("\n  %s\n\n  Press q to quit.\n", errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	title := titleStyle.Width(m.width).Align(lipgloss.Center).Render("SomaFM")

	listWidth := m.width / 3
	if listWidth < 25 {
		listWidth = 25
	}
	detailWidth := m.width - listWidth - 6
	if detailWidth < 10 {
		detailWidth = 10
	}
	contentHeight := m.height - 7
	if contentHeight < 5 {
		contentHeight = 5
	}

	list := m.channelListView(listWidth, contentHeight)
	details := m.detailsView(detailWidth, contentHeight)
	content := lipgloss.JoinHorizontal(lipgloss.Top, list, details)

	var filter string
	if m.filterMode {
		filter = "\n  " + m.filterInput.View()
	} else if m.filterInput.Value() != "" {
		filter = fmt.Sprintf("\n  Filter: %s (esc to clear)", m.filterInput.Value())
	}

	status := m.statusView()
	help := helpStyle.Render("  enter:play  space:pause  +/-:volume  /:filter  tab:switch panel  q:quit")

	return title + "\n" + content + filter + "\n" + status + "\n" + help
}

func (m model) channelListView(width, height int) string {
	var lines []string
	for i, idx := range m.filtered {
		ch := m.channels[idx]
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		playing := ""
		if idx == m.playingIdx {
			playing = " ♫"
		}

		line := cursor + ch.Title + playing

		switch {
		case i == m.cursor && idx == m.playingIdx:
			line = selectedPlayingStyle.Render(line)
		case i == m.cursor:
			line = selectedStyle.Render(line)
		case idx == m.playingIdx:
			line = playingStyle.Render(line)
		}

		lines = append(lines, line)
	}

	start := 0
	if m.cursor >= height {
		start = m.cursor - height + 1
	}
	if m.cursor < start {
		start = m.cursor
	}
	end := start + height
	if end > len(lines) {
		end = len(lines)
	}
	if start > len(lines) {
		start = len(lines)
	}

	visible := strings.Join(lines[start:end], "\n")
	border := borderStyle
	if m.focus == panelChannels {
		border = focusedBorderStyle
	}
	return border.Width(width).Height(height).Render(visible)
}

func (m model) detailsView(width, height int) string {
	var content string
	if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
		ch := m.channels[m.filtered[m.cursor]]
		content = ch.GetDetails()
	}

	// Wrap content to the available inner width before counting lines,
	// so that long lines that visually wrap are properly accounted for.
	innerWidth := width - 4 // detailStyle horizontal padding (2 left + 2 right)
	if innerWidth < 1 {
		innerWidth = 1
	}
	wrapped := lipgloss.NewStyle().Width(innerWidth).Render(content)
	lines := strings.Split(wrapped, "\n")
	visibleLines := height - 2 // account for detailStyle vertical padding
	if visibleLines < 1 {
		visibleLines = 1
	}

	maxScroll := len(lines) - visibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}
	scroll := m.detailScroll
	if scroll > maxScroll {
		scroll = maxScroll
	}

	end := scroll + visibleLines
	if end > len(lines) {
		end = len(lines)
	}

	visible := strings.Join(lines[scroll:end], "\n")
	border := borderStyle
	if m.focus == panelDetails {
		border = focusedBorderStyle
	}
	return border.Width(width).Height(height).Render(
		detailStyle.Render(visible),
	)
}

func (m model) statusView() string {
	if m.err != nil && m.state == stateReady {
		return errorStyle.Render(fmt.Sprintf("  Error: %v", m.err))
	}

	if m.connecting {
		return statusStyle.Render(fmt.Sprintf("  %s Connecting...", m.spinner.View()))
	}

	if m.playingIdx >= 0 && m.playingIdx < len(m.channels) && m.player.IsPlaying() {
		ch := m.channels[m.playingIdx]
		state := "Playing"
		if m.player.IsPaused() {
			state = "Paused"
		}

		vol := m.player.GetVolume()
		volBar := volumeBar(vol)

		return statusStyle.Render(fmt.Sprintf("  ♫ %s: %s  |  Vol: %s", state, ch.Title, volBar))
	}

	return statusStyle.Render("  Select a channel and press enter to play")
}

func volumeBar(vol float64) string {
	normalized := int(vol + 5)
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 7 {
		normalized = 7
	}
	filled := strings.Repeat("█", normalized)
	empty := strings.Repeat("░", 7-normalized)
	return filled + empty
}
