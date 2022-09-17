package state

import "github.com/nicarl/somafm/radioChannels"

type PlayerState struct {
	Channels        []radioChannels.RadioChan
	SelectedCh      int
	SelectedControl Control
	IsPlaying       bool
}

type Control string

const (
	LOUDER      Control = "LOUDER"
	QUIETER     Control = "QUIETER"
	PLAY_BUTTON Control = "PLAY_BUTTON"
)

func (state *PlayerState) SelectNextCh() {
	if state.SelectedCh < len(state.Channels)-1 {
		state.SelectedCh++
	}
}

func (state *PlayerState) SelectPrevCh() {
	if state.SelectedCh > 0 {
		state.SelectedCh--
	}
}

func (state *PlayerState) SelectControl(control Control) {
	state.SelectedControl = control
}

func (state *PlayerState) PauseMusic() {
	state.IsPlaying = false
}

func (state *PlayerState) PlayMusic() {
	state.IsPlaying = true
}

func InitState(channels []radioChannels.RadioChan) *PlayerState {
	return &PlayerState{
		Channels:        channels,
		SelectedCh:      0,
		SelectedControl: PLAY_BUTTON,
		IsPlaying:       false,
	}
}
