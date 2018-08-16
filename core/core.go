package core

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"

	"github.com/Southclaws/cordless/interface/chanlist"
	"github.com/Southclaws/cordless/interface/chatlog"
)

// App stores application state and handles all events
type App struct {
	Discord  *discordgo.Session
	Renderer *tview.Application

	config   config
	chanlist *chanlist.ChanList
	chatlog  *chatlog.ChatLog
}

// Initialise sets up the UI ready for running
func Initialise() (app *App, err error) {
	app = new(App)

	app.config, err = loadConfig()
	if err != nil {
		err = errors.Wrap(err, "failed to load config")
		return
	}

	app.Renderer = tview.NewApplication()

	app.chanlist = chanlist.Create(app.config.CurrentGuild, app.config.CurrentChannel)
	app.chatlog = chatlog.Create()

	flex := tview.NewFlex().
		AddItem(app.chanlist.Inner, 0, 1, false).
		AddItem(app.chatlog.Inner, 0, 7, false)

	app.Renderer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == 0 && event.Key() == tcell.KeyRune && !app.chatlog.Input.HasFocus() {
			// if input doesn't have focus and the event is not a modifier
			// write the captured key to the input field - this prevents the
			// awkward feeling when the first key you hit will focus the input
			// field but not actually insert the key so you have to hit it again
			app.chatlog.Input.SetText(app.chatlog.Input.GetText() + string(event.Rune()))
			app.Renderer.SetFocus(app.chatlog.Input)
		} else if event.Key() == tcell.KeyCtrlS {
			app.chatlog.Sys("Getting guild list...")
			app.Renderer.SetFocus(app.chanlist.Inner)
			app.chanlist.PromptForGuild()
			app.chatlog.Sys("Awaiting selection...")
		} else if event.Key() == tcell.KeyCtrlA {
			if app.config.CurrentGuild == "" {
				app.chatlog.Sys("No guild selected, use ^S to select a guild.")
			} else {
				app.chatlog.Sys("Getting channel list...")
				app.Renderer.SetFocus(app.chanlist.Inner)
				app.chanlist.PromptForChan()
				app.chatlog.Sys("Awaiting selection...")
			}
		} else if (event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown) &&
			(app.chatlog.Input.HasFocus() || app.chatlog.Log.HasFocus()) {
			app.Renderer.SetFocus(app.chatlog.Log)
		}

		return event
	})

	fmt.Println("Connecting to Discord...")
	app.Discord, err = discordgo.New(app.config.Token)
	if err != nil {
		fmt.Println(err)
		return
	}

	app.chanlist.OnSelectGuild = app.onSelectGuild
	app.chanlist.OnSelectChan = app.onSelectChan
	app.chatlog.OnUpdate = func() {
		app.Renderer.Draw()
	}

	app.Discord.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		app.chatlog.Sys(fmt.Sprint(
			"Connected! Current guild: ",
			app.chanlist.ActiveGuild,
			" current channel: ",
			app.chanlist.ActiveChan,
		))
		app.chanlist.SetChannels(m.Guilds)
	})
	app.Discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if app.chanlist.ActiveChan == m.ChannelID {
			app.chatlog.Write(m.Author.Username, m.Content)
		} else {
			app.chanlist.SetUnread(m.ChannelID)
		}
	})

	err = app.Discord.Open()
	if err != nil {
		return
	}

	app.Renderer.SetRoot(flex, true)
	return app, err
}

// Run runs the application and blocks until an error occurs
func (app *App) Run() (err error) {
	return app.Renderer.Run()
}

func (app *App) onSelectGuild(id string) {
	app.chatlog.Sys(fmt.Sprint("Selected guild ", id))
	app.chanlist.PromptForChan()

	app.config.CurrentGuild = id
	err := updateConfig(app.config)
	if err != nil {
		panic(err)
	}
}

func (app *App) onSelectChan(id string) {
	app.chatlog.Sys(fmt.Sprint("Selected channel ", id))

	app.config.CurrentChannel = id
	err := updateConfig(app.config)
	if err != nil {
		panic(err)
	}
}
