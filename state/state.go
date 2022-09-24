package state

import "github.com/nicarl/somafm/radioChannels"

type PlayerState struct {
	Channels   []radioChannels.RadioChan
	SelectedCh int
	IsPlaying  bool
}

func (state *PlayerState) SelectCh(i int) {
	state.SelectedCh = i
}

func (state *PlayerState) GetSelectedCh() radioChannels.RadioChan {
	return state.Channels[state.SelectedCh]
}

func (state *PlayerState) PauseMusic() {
	state.IsPlaying = false
}

func (state *PlayerState) PlayMusic() {
	state.IsPlaying = true
}

func InitState(channels []radioChannels.RadioChan) *PlayerState {
	return &PlayerState{
		Channels:   channels,
		SelectedCh: 0,
		IsPlaying:  false,
	}
}
