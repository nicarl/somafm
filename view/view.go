package view

import (
	"fmt"
	"log"

	"github.com/nicarl/somafm/state"
	"github.com/rivo/tview"

	tcell "github.com/gdamore/tcell/v2"
)

func getChannelList(
	appState *state.AppState,
	app *tview.Application,
	channelDetails *tview.TextView,
	player *tview.List,
) *tview.List {
	channelList := tview.NewList()
	channelList.SetBorder(true).SetTitle("Channels")
	channelList.ShowSecondaryText(false)
	channelList.SetBorderPadding(1, 1, 1, 1)

	for _, radioCh := range appState.Channels {
		channelList.AddItem(radioCh.Title, "", 0, func() {
			err := appState.PlayMusic()
			if err != nil {
				app.Stop()
				log.Fatalf("%+v", err)
			}
			player.SetItemText(1, "Pause", "")
		})
	}
	channelList.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		appState.SelectCh(i)
		channelDetails.Clear()
		fmt.Fprint(channelDetails, appState.GetSelectedCh().GetDetails())
	})
	return channelList
}

func getChannelDetails(appState *state.AppState) *tview.TextView {
	channelDetails := tview.NewTextView()
	channelDetails.SetBorder(true).SetTitle("Details")
	channelDetails.SetBorderPadding(1, 0, 0, 0)
	fmt.Fprint(channelDetails, appState.GetSelectedCh().GetDetails())
	return channelDetails
}

func getPlayer(
	appState *state.AppState,
	app *tview.Application,
) *tview.List {
	player := tview.NewList()
	player.SetBorder(true)
	player.ShowSecondaryText(false)
	player.SetBorderPadding(1, 1, 1, 1)

	player.AddItem("Volume +", "", 0, func() {
		appState.IncreaseVolume()
	})
	player.AddItem("Play", "", 0, func() {
		if appState.IsPlaying {
			appState.PauseMusic()
			player.SetItemText(1, "Play", "")
		} else {
			err := appState.PlayMusic()
			if err != nil {
				app.Stop()
				log.Fatalf("%+v", err)
			}
			player.SetItemText(1, "Pause", "")
		}
	})
	player.AddItem("Volume -", "", 0, func() {
		appState.DecreaseVolume()
	})

	player.SetCurrentItem(1)
	return player
}

func InitApp(appState *state.AppState) error {
	app := tview.NewApplication()

	channelDetails := getChannelDetails(appState)
	player := getPlayer(appState, app)
	channelList := getChannelList(appState, app, channelDetails, player)

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
	flexWithHeader := tview.NewFrame(flex).
		SetBorders(2, 2, 2, 2, 4, 4).
		AddText("SomaFM", true, tview.AlignCenter, tcell.ColorWhite)

	app.SetRoot(flexWithHeader, true).SetFocus(channelList)

	if err := app.Run(); err != nil {
		return err
	}
	return nil
}
