package state

import (
	"github.com/nicarl/somafm/audio"
	"github.com/nicarl/somafm/radioChannels"
)

type AppState struct {
	Channels   []radioChannels.RadioChan
	SelectedCh int
	IsPlaying  bool
	done       chan bool
	setVolume  chan float32
	errs       chan error
}

func (appState *AppState) SelectCh(i int) {
	appState.SelectedCh = i
}

func (appState *AppState) GetSelectedCh() radioChannels.RadioChan {
	return appState.Channels[appState.SelectedCh]
}

func (appState *AppState) PauseMusic() {
	if appState.IsPlaying {
		appState.IsPlaying = false
		appState.done <- true
	}
}

func (appState *AppState) PlayMusic() {
	if appState.IsPlaying {
		appState.PauseMusic()
	}
	appState.IsPlaying = true
	go audio.PlayMusic(appState.GetSelectedCh().StreamURL, appState.done, appState.setVolume, appState.errs)
}

func (appState *AppState) IncreaseVolume() {
	appState.setVolume <- 0.5
}

func (appState *AppState) DecreaseVolume() {
	appState.setVolume <- -0.5
}

func InitState(channels []radioChannels.RadioChan) *AppState {
	return &AppState{
		Channels:   channels,
		SelectedCh: 0,
		IsPlaying:  false,
		done:       make(chan bool),
		setVolume:  make(chan float32, 10),
		errs:       make(chan error, 1),
	}
}
