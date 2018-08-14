package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/Southclaws/cordless/interface/chanlist"
	"github.com/Southclaws/cordless/interface/chatlog"
)

func main() {
	app := tview.NewApplication()

	chanlist := chanlist.Create()
	chanlist.SetChannels(map[string][]string{
		"SA:MP": {"#programming", "#general", "#verified"},
	})
	chatlog := chatlog.Create()

	flex := tview.NewFlex().
		AddItem(chanlist.Inner, 0, 1, false).
		AddItem(chatlog.Inner, 0, 7, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == 0 && event.Key() == tcell.KeyRune && !chatlog.Input.HasFocus() {
			// if input doesn't have focus and the event is not a modifier
			// write the captured key to the input field - this prevents the
			// awkward feeling when the first key you hit will focus the input
			// field but not actually insert the key so you have to hit it again
			chatlog.Input.SetText(chatlog.Input.GetText() + string(event.Rune()))
			app.SetFocus(chatlog.Input)
		} else if event.Key() == tcell.KeyCtrlA {
			app.SetFocus(chanlist.Inner)
		}

		return event
	})

	app.SetRoot(flex, true)
	err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
}
