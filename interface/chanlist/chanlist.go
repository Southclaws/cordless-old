package chanlist

import (
	"github.com/rivo/tview"
)

// ChanList renders a list of Discord servers and their channels
type ChanList struct {
	Inner *tview.List
}

// Create creates a ChanList and initialises it ready for rendering
func Create() (chanlist *ChanList) {
	chanlist = &ChanList{
		Inner: tview.NewList(),
	}
	chanlist.Inner.SetBorder(true)
	chanlist.Inner.AddItem("Loading Channels...", "", 0, nil)
	return
}

// SetChannels takes a hierarchy of servers and channels and renders it
func (e ChanList) SetChannels(data map[string][]string) {
	e.Inner.Clear()
	for server, channels := range data {
		e.Inner.AddItem(server, "", 0, nil)
		for _, channel := range channels {
			e.Inner.AddItem(channel, "", 0, nil)
		}
	}
}
