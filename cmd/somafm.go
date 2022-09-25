package main

import (
	"log"

	"github.com/nicarl/somafm/audio"
	"github.com/nicarl/somafm/radioChannels"
	"github.com/nicarl/somafm/state"
	"github.com/nicarl/somafm/view"
)

func main() {
	if err := audio.InitSpeaker(); err != nil {
		log.Fatalf("%+v", err)
	}

	radioCh, err := radioChannels.GetChannels()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	playerState := state.InitState(radioCh)

	view.InitApp(playerState)
}
