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

func InitApp(appState *state.AppState) error {
	app := tview.NewApplication()

	channelDetails := getChannelDetails(appState)
	channelList := getChannelList(appState, app, channelDetails)

	flex := tview.NewFlex().
		AddItem(channelList, 0, 1, false).
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
