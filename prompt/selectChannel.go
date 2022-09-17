package prompt

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nicarl/somafm/state"
)

func SelectChannel(s tcell.Screen, playerState *state.PlayerState, quit func()) {
	for {
		s.Clear()
		drawRadioChannels(s, *playerState)
		s.Show()
		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			} else if ev.Key() == tcell.KeyDown {
				playerState.SelectNextCh()
			} else if ev.Key() == tcell.KeyUp {
				playerState.SelectPrevCh()
			} else if ev.Key() == tcell.KeyEnter {
				return
			}
		}
	}
}
