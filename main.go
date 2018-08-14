package main

import (
	"fmt"

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

	app.SetFocus(chatlog.Input)
	app.SetRoot(flex, true)
	err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
}
