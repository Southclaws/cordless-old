package chatlog

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// ChatLog renders a list of Discord servers and their channels
type ChatLog struct {
	Inner *tview.Flex
	Log   *tview.TextView
	Input *tview.InputField
}

// Create creates a ChatLog and initialises it ready for rendering
func Create() (chanlist *ChatLog) {
	log := tview.NewTextView()
	log.SetBorder(true).SetTitle("Chat Log")

	input := tview.NewInputField()
	input.SetLabel("Message:")
	input.SetDoneFunc(func(key tcell.Key) {
		log.Write([]byte(input.GetText()))
		input.SetText("")
	})
	input.SetBorder(true).SetTitle("Input")

	log.Write([]byte("line\nline2"))

	container := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(log, 0, 1, false).
		AddItem(input, 3, 0, true)

	chanlist = &ChatLog{
		Inner: container,
		Log:   log,
		Input: input,
	}

	return
}
