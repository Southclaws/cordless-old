package chanlist

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rivo/tview"
)

// ChanList renders a list of Discord servers and their channels
type ChanList struct {
	Inner         *tview.List
	Guilds        map[string]Guild
	OnSelectGuild func(id string)
	OnSelectChan  func(id string)

	ActiveGuild string
	ActiveChan  string
}

// Guild contains a list of channels
type Guild struct {
	ID       string
	Name     string
	Channels map[string]Channel
}

// Channel wraps a discord channel and adds an Unread field
type Channel struct {
	ID     string
	Name   string
	Unread bool
}

// Create creates a ChanList and initialises it ready for rendering
func Create(ActiveGuild, ActiveChan string) (chanlist *ChanList) {
	chanlist = &ChanList{
		Inner:  tview.NewList(),
		Guilds: make(map[string]Guild),

		ActiveGuild: ActiveGuild,
		ActiveChan:  ActiveChan,
	}
	chanlist.Inner.ShowSecondaryText(false)
	chanlist.Inner.SetBorder(true)
	return
}

// SetChannels takes a hierarchy of servers and channels and renders it
func (e *ChanList) SetChannels(guilds []*discordgo.Guild) {
	for _, g := range guilds {
		e.Guilds[g.ID] = Guild{
			ID:       g.ID,
			Name:     g.Name,
			Channels: make(map[string]Channel),
		}
		for _, c := range g.Channels {
			e.Guilds[g.ID].Channels[c.ID] = Channel{
				ID:     c.ID,
				Name:   c.Name,
				Unread: false, // Todo: retain unread status
			}
		}
	}
	if e.ActiveGuild != "" {
		e.PromptForChan()
	}
}

// PromptForGuild prompts the user to select a guild
func (e *ChanList) PromptForGuild() {
	e.Inner.Clear()
	for _, guild := range e.Guilds {
		e.Inner.AddItem("•"+guild.Name, "", 0, func() {
			e.OnSelectGuild(guild.ID)
			e.ActiveGuild = guild.ID
		})
	}
	e.Inner.SetCurrentItem(0)
}

// PromptForChan prompts the user to select a channel from the specified guild
func (e *ChanList) PromptForChan() {
	if e.ActiveGuild == "" {
		return
	}

	e.Inner.Clear()
	var guild *Guild
	for _, g := range e.Guilds {
		if g.ID == e.ActiveGuild {
			guild = &g
		}
	}
	if guild == nil {
		panic(fmt.Sprint("could not find guild in tree with the ID", e.ActiveGuild))
	}

	activeChanIdx := -1
	idx := 0
	for _, channel := range guild.Channels {
		if channel.ID == e.ActiveChan {
			activeChanIdx = idx
		}
		chanName := fmt.Sprintf("#%s - %s", channel.Name, channel.ID)
		if channel.Unread {
			chanName += "*"
		}
		e.Inner.AddItem(chanName, "", 0, func() {
			e.OnSelectChan(channel.ID)
			e.ActiveChan = channel.ID
		})
		idx++
	}

	if activeChanIdx != -1 {
		e.Inner.SetCurrentItem(activeChanIdx)
	}
}

// SetUnread marks a channel as unread and updates the UI to reflect this
func (e *ChanList) SetUnread(channelID string) {
	for _, g := range e.Guilds {
		for _, c := range g.Channels {
			c.Unread = true
		}
	}
}
