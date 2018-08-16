package chatlog

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// ChatLog renders a list of Discord servers and their channels
type ChatLog struct {
	Inner    *tview.Flex
	Log      *tview.TextView
	Input    *tview.InputField
	OnUpdate func()
}

// Create creates a ChatLog and initialises it ready for rendering
func Create() (chatlog *ChatLog) {
	chatlog = new(ChatLog)

	chatlog.Log = chatlog.createChatLog()
	chatlog.Input = chatlog.createInput()

	chatlog.Inner = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chatlog.Log, 0, 1, false).
		AddItem(chatlog.Input, 3, 0, true)

	return
}

// Write sends a message to the chat log
func (cl *ChatLog) Write(user, message string) {
	_, err := cl.Log.Write([]byte(fmt.Sprintf("%s: %s\n", user, message)))
	if err != nil {
		panic(err)
	}
	cl.Log.ScrollToEnd()
	cl.OnUpdate()
}

// Sys sends a system/debug message
func (cl *ChatLog) Sys(message string) {
	_, err := cl.Log.Write([]byte(fmt.Sprintf("SYSTEM: %s\n", message)))
	if err != nil {
		panic(err)
	}
	cl.Log.ScrollToEnd()
	cl.OnUpdate()
}

func (cl *ChatLog) createChatLog() (element *tview.TextView) {
	element = tview.NewTextView()
	element.SetBorder(true)
	element.SetTitle("Chat Log")
	return
}

func (cl *ChatLog) createInput() (element *tview.InputField) {
	element = tview.NewInputField()
	element.SetPlaceholder("message")
	element.SetDoneFunc(func(key tcell.Key) {
		text := element.GetText()
		if text == "" {
			return
		}

		_, err := cl.Log.Write([]byte(text + "\n"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		element.SetText("")
	})
	element.SetBorder(true).SetTitle("Input")
	return
}
