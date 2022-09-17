package main

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/nicarl/somafm/prompt"
	"github.com/nicarl/somafm/radioChannels"
	"github.com/nicarl/somafm/state"

	"github.com/gdamore/tcell/v2"
)

func main() {
	sampleRate := beep.SampleRate(44100)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	radioCh, err := radioChannels.GetChannels()
	if err != nil {
		log.Fatal(err)
	}

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	selectedCh := 0
	playerState := state.PlayerState{
		Channels:        radioCh,
		SelectedCh:      selectedCh,
		SelectedControl: state.PLAY_BUTTON,
		IsPlaying:       false,
	}
	s.Clear()
	quit := func() {
		s.Fini()
		os.Exit(0)
	}
	for {
		prompt.SelectChannel(s, &playerState, quit)
		prompt.PlayChannel(s, &playerState, quit)
	}
}
