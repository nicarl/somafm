package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case channelsLoadedMsg:
		m.state = stateReady
		m.channels = msg.channels
		m.updateFilter()
		if len(m.filtered) > 0 {
			ch := m.channels[m.filtered[m.cursor]]
			return m, fetchSongsCmd(ch.ID)
		}
		return m, nil

	case refreshTickMsg:
		cmds := []tea.Cmd{refreshTickCmd()}
		if m.state == stateReady && len(m.filtered) > 0 && m.cursor < len(m.filtered) {
			ch := m.channels[m.filtered[m.cursor]]
			cmds = append(cmds, fetchSongsCmd(ch.ID))
		}
		return m, tea.Batch(cmds...)

	case songsRefreshedMsg:
		for i := range m.channels {
			if m.channels[i].ID == msg.channelID {
				m.channels[i].Songs = msg.songs
				break
			}
		}
		return m, nil

	case errMsg:
		m.state = stateError
		m.err = msg.err
		return m, nil

	case playStartedMsg:
		m.connecting = false
		m.err = nil
		return m, nil

	case playErrorMsg:
		m.connecting = false
		m.playingIdx = -1
		m.err = msg.err
		return m, nil

	case spinner.TickMsg:
		if m.state == stateLoading || m.connecting {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.filterMode {
		return m.handleFilterKey(msg)
	}

	switch msg.String() {
	case "q", "ctrl+c":
		m.player.Stop()
		return m, tea.Quit

	case "tab":
		if m.focus == panelChannels {
			m.focus = panelDetails
		} else {
			m.focus = panelChannels
		}

	case "up", "k":
		if m.focus == panelDetails {
			if m.detailScroll > 0 {
				m.detailScroll--
			}
		} else if m.cursor > 0 && len(m.filtered) > 0 {
			m.cursor--
			m.detailScroll = 0
			ch := m.channels[m.filtered[m.cursor]]
			return m, fetchSongsCmd(ch.ID)
		}

	case "down", "j":
		if m.focus == panelDetails {
			m.detailScroll++
		} else if len(m.filtered) > 0 && m.cursor < len(m.filtered)-1 {
			m.cursor++
			m.detailScroll = 0
			ch := m.channels[m.filtered[m.cursor]]
			return m, fetchSongsCmd(ch.ID)
		}

	case "enter":
		if len(m.filtered) > 0 && !m.connecting {
			idx := m.filtered[m.cursor]
			m.playingIdx = idx
			m.connecting = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.playCmd(idx))
		}

	case " ":
		m.player.TogglePause()

	case "+", "=":
		m.player.VolumeUp()

	case "-":
		m.player.VolumeDown()

	case "/":
		m.filterMode = true
		m.filterInput.Focus()
		return m, m.filterInput.Cursor.BlinkCmd()

	case "esc":
		if m.filterInput.Value() != "" {
			m.filterInput.SetValue("")
			m.updateFilter()
			m.cursor = 0
			m.detailScroll = 0
		}
	}

	return m, nil
}

func (m model) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "esc":
		m.filterMode = false
		m.filterInput.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)
	m.updateFilter()
	m.cursor = 0
	m.detailScroll = 0
	return m, cmd
}

func (m *model) updateFilter() {
	filter := strings.ToLower(m.filterInput.Value())
	m.filtered = nil
	for i, ch := range m.channels {
		if filter == "" ||
			strings.Contains(strings.ToLower(ch.Title), filter) ||
			strings.Contains(strings.ToLower(ch.Genre), filter) {
			m.filtered = append(m.filtered, i)
		}
	}
}

func (m model) playCmd(idx int) tea.Cmd {
	streamURL := m.channels[idx].StreamURL
	player := m.player
	return func() tea.Msg {
		err := player.Play(streamURL)
		if err != nil {
			return playErrorMsg{err}
		}
		return playStartedMsg{}
	}
}
