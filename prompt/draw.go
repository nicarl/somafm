package prompt

import (
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/nicarl/somafm/state"
)

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func drawRadioChannels(s tcell.Screen, playerState state.PlayerState) {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorDefault)
	selectedStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorDarkCyan)

	width, height := s.Size()

	start := func() int {
		if playerState.SelectedCh >= height {
			return int(math.Min(float64(len(playerState.Channels)), float64(playerState.SelectedCh-height+1)))
		}
		return 0
	}()
	end := int(math.Min(float64(len(playerState.Channels)), float64(height+start)))

	for i := 0; i < end-start; i++ {
		if (i + start) == playerState.SelectedCh {
			drawText(s, 1, i, width, i, selectedStyle, playerState.Channels[i+start].Title)
		} else {
			drawText(s, 1, i, width, i, defStyle, playerState.Channels[i+start].Title)
		}
	}
}

func drawPlayer(s tcell.Screen, playerState state.PlayerState) {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorDefault)
	selectedStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorDarkCyan)

	_, height := s.Size()

	var playButtonText string
	if playerState.IsPlaying {
		playButtonText = "Pause ⏸"
	} else {
		playButtonText = "Play  ▶"
	}
	if playerState.SelectedControl == state.LOUDER {
		drawText(s, 8, height/2-2, 7, height/2-2, selectedStyle, "⏶")
		drawText(s, 8, height/2+2, 7, height/2+2, defStyle, "⏷")
		drawText(s, 5, height/2, 12, height/2, defStyle, playButtonText)

	} else if playerState.SelectedControl == state.PLAY_BUTTON {
		drawText(s, 8, height/2-2, 7, height/2-2, defStyle, "⏶")
		drawText(s, 8, height/2+2, 7, height/2+2, defStyle, "⏷")
		drawText(s, 5, height/2, 12, height/2, selectedStyle, playButtonText)
	} else {
		drawText(s, 8, height/2-2, 7, height/2-2, defStyle, "⏶")
		drawText(s, 8, height/2+2, 7, height/2+2, selectedStyle, "⏷")
		drawText(s, 5, height/2, 12, height/2, defStyle, playButtonText)
	}
}
