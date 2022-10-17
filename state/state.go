package state

import (
	"github.com/nicarl/somafm/audio"
	"github.com/nicarl/somafm/radioChannels"
)

type AppState struct {
	Channels   []radioChannels.RadioChan
	SelectedCh int
	PlayingCh  int
	IsPlaying  bool
	done       chan bool
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

func (appState *AppState) PlayMusic() error {
	if appState.IsPlaying {
		if appState.SelectedCh == appState.PlayingCh {
			return nil
		}
		appState.PauseMusic()
	}
	errs := make(chan error, 1)
	appState.IsPlaying = true
	appState.PlayingCh = appState.SelectedCh
	go audio.PlayMusic(appState.GetSelectedCh().StreamURL, appState.done, errs)

	return <-errs
}

func InitState(channels []radioChannels.RadioChan) *AppState {
	return &AppState{
		Channels:   channels,
		SelectedCh: 0,
		IsPlaying:  false,
		done:       make(chan bool),
	}
}
