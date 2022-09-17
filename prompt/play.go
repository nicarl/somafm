package prompt

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nicarl/somafm/audio"
	"github.com/nicarl/somafm/state"
)

type Control string

func PlayChannel(s tcell.Screen, playerState *state.PlayerState, quit func()) {
	done := make(chan bool)
	setVolume := make(chan float32)

	for {
		s.Clear()
		drawPlayer(s, *playerState)
		s.Show()
		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				if playerState.IsPlaying {
					done <- true
					playerState.PauseMusic()
				}
				return
			} else if ev.Key() == tcell.KeyCtrlC {
				quit()
			} else if ev.Key() == tcell.KeyDown {
				if playerState.SelectedControl == state.LOUDER {
					playerState.SelectControl(state.PLAY_BUTTON)
				} else {
					playerState.SelectControl(state.QUIETER)
				}
			} else if ev.Key() == tcell.KeyUp {
				if playerState.SelectedControl == state.QUIETER {
					playerState.SelectControl(state.PLAY_BUTTON)
				} else {
					playerState.SelectControl(state.LOUDER)
				}
			} else if ev.Key() == tcell.KeyEnter {
				if playerState.IsPlaying {
					switch playerState.SelectedControl {
					case state.PLAY_BUTTON:
						done <- true
						playerState.PauseMusic()
					case state.LOUDER:
						setVolume <- 0.5
					case state.QUIETER:
						setVolume <- -0.5
					}
				} else if playerState.SelectedControl == state.PLAY_BUTTON {
					go audio.PlayMusic(playerState.Channels[playerState.SelectedCh].StreamURL, done, setVolume)
					playerState.PlayMusic()
				}
			}
		}
	}
}
