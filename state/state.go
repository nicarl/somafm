package state

import "github.com/nicarl/somafm/radioChannels"

type AppState struct {
	Channels   []radioChannels.RadioChan
	SelectedCh int
	IsPlaying  bool
}

func (state *AppState) SelectCh(i int) {
	state.SelectedCh = i
}

func (state *AppState) GetSelectedCh() radioChannels.RadioChan {
	return state.Channels[state.SelectedCh]
}

func (state *AppState) PauseMusic() {
	state.IsPlaying = false
}

func (state *AppState) PlayMusic() {
	state.IsPlaying = true
}

func InitState(channels []radioChannels.RadioChan) *AppState {
	return &AppState{
		Channels:   channels,
		SelectedCh: 0,
		IsPlaying:  false,
	}
}
