package prompt

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/nicarl/somafm/audio"
	"github.com/nicarl/somafm/state"
	"github.com/rivo/tview"
)

func getChannelList(
	appState *state.AppState,
	channelDetails *tview.TextView,
	player *tview.List,
	done chan bool,
	setVolume <-chan float32,
	errs chan<- error,
) *tview.List {
	channelList := tview.NewList()
	channelList.SetBorder(true).SetTitle("Channels")

	for _, radioCh := range appState.Channels {
		channelList.AddItem(radioCh.Title, "", 0, func() {
			if appState.IsPlaying {
				appState.PauseMusic()
				done <- true
			}
			appState.PlayMusic()
			go audio.PlayMusic(appState.GetSelectedCh().StreamURL, done, setVolume, errs)
			player.SetItemText(1, "Pause", "")
		})
	}
	channelList.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		appState.SelectCh(i)
		channelDetails.Clear()
		fmt.Fprint(channelDetails, appState.GetSelectedCh().Description)
	})
	return channelList
}

func getChannelDetails(appState *state.AppState) *tview.TextView {
	channelDetails := tview.NewTextView()
	channelDetails.SetBorder(true).SetTitle("Details")
	fmt.Fprint(channelDetails, appState.GetSelectedCh().Description)
	return channelDetails
}

func getPlayer(
	appState *state.AppState,
	done chan bool,
	setVolume chan float32,
	errs chan<- error,
) *tview.List {
	player := tview.NewList()
	player.SetBorder(true)
	player.AddItem("Volume +", "", 0, func() {
		setVolume <- 0.5
	})
	player.AddItem("Play", "", 0, func() {
		if appState.IsPlaying {
			appState.PauseMusic()
			done <- true
			player.SetItemText(1, "Play", "")
		} else {
			appState.PlayMusic()
			go audio.PlayMusic(appState.GetSelectedCh().StreamURL, done, setVolume, errs)
			player.SetItemText(1, "Pause", "")
		}
	})
	player.AddItem("Volume -", "", 0, func() {
		setVolume <- -0.5
	})

	player.SetCurrentItem(1)
	return player
}

func InitApp(appState *state.AppState) {
	app := tview.NewApplication()

	done := make(chan bool)
	setVolume := make(chan float32, 10)
	errs := make(chan error, 1)

	channelDetails := getChannelDetails(appState)
	player := getPlayer(appState, done, setVolume, errs)
	channelList := getChannelList(appState, channelDetails, player, done, setVolume, errs)

	channelList.SetSelectedFunc(func(_ int, _ string, _ string, _ rune) {
		app.SetFocus(player)
	})
	player.SetDoneFunc(func() {
		app.SetFocus(channelList)
	})

	flex := tview.NewFlex().
		AddItem(channelList, 0, 1, false).
		AddItem(player, 0, 1, false).
		AddItem(channelDetails, 0, 1, false)
	flexWithHeader := tview.NewFrame(flex).SetBorders(2, 2, 2, 2, 4, 4).AddText("SomaFM", true, tview.AlignCenter, tcell.ColorDefault)

	if err := app.SetRoot(flexWithHeader, true).SetFocus(channelList).Run(); err != nil {
		panic(err)
	}
}
