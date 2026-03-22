package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/nicarl/somafm/audio"
	"github.com/nicarl/somafm/radiochannels"
)

type appState int

const (
	stateLoading appState = iota
	stateReady
	stateError
)

type focusedPanel int

const (
	panelChannels focusedPanel = iota
	panelDetails
)

type channelsLoadedMsg struct {
	channels []radiochannels.RadioChan
}

type errMsg struct{ err error }
type playStartedMsg struct{}
type playErrorMsg struct{ err error }
type refreshTickMsg time.Time
type songsRefreshedMsg struct {
	channelID string
	songs     []radiochannels.Song
}

type model struct {
	state       appState
	channels    []radiochannels.RadioChan
	filtered    []int
	cursor      int
	playingIdx  int
	connecting  bool
	player      *audio.Player
	spinner     spinner.Model
	filterInput textinput.Model
	filterMode  bool
	detailScroll int
	focus        focusedPanel
	width        int
	height       int
	err          error
}

func newModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	fi := textinput.New()
	fi.Placeholder = "Filter channels..."
	fi.CharLimit = 50

	return model{
		state:       stateLoading,
		playingIdx:  -1,
		player:      audio.NewPlayer(),
		spinner:     s,
		filterInput: fi,
	}
}

func refreshTickCmd() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return refreshTickMsg(t)
	})
}

func fetchSongsCmd(channelID string) tea.Cmd {
	return func() tea.Msg {
		songs, err := radiochannels.FetchSongs(channelID)
		if err != nil {
			return nil
		}
		return songsRefreshedMsg{channelID: channelID, songs: songs}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchChannelsCmd, refreshTickCmd())
}

func fetchChannelsCmd() tea.Msg {
	channels, err := radiochannels.GetChannels()
	if err != nil {
		return errMsg{err}
	}
	return channelsLoadedMsg{channels}
}

func Run() error {
	m := newModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	m.player.Stop()
	return err
}
