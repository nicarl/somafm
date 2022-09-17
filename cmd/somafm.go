package main

import (
	"log"

	"github.com/nicarl/somafm/prompt"
	"github.com/nicarl/somafm/radioChannels"
	"github.com/nicarl/somafm/state"

	"github.com/nicarl/somafm/audio"
)

func main() {
	if err := audio.InitSpeaker(); err != nil {
		log.Fatalf("%+v", err)
	}

	radioCh, err := radioChannels.GetChannels()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	s, quit, err := prompt.InitScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	playerState := state.InitState(radioCh)
	for {
		prompt.SelectChannel(s, playerState, quit)
		err := prompt.PlayChannel(s, playerState, quit)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
}
