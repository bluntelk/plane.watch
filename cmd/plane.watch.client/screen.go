package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"plane.watch/lib/tracker"
	"time"
)

type (
	display struct {
		app    *tview.Application
		top    *tview.Table
		bottom *tview.TextView
	}
	pacerEvent struct {

	}
)

func (p pacerEvent) Type() string {
	return "pacer"
}

func (p pacerEvent) String() string {
	return "beep"
}

func newPacer() *pacerEvent {
	return &pacerEvent{}
}

func (d *display) Run() error {
	go func() {
		c := time.NewTicker(time.Second)
		select {
		case <- c.C:
			d.OnEvent(newPacer())
		}
	}()
	if err := d.app.Run(); nil != err {
		return err
	}
	return nil
}

func newAppDisplay() (*display, error) {
	d := display{}
	d.app = tview.NewApplication()
	hdr := func(title string, width int) *tview.TableCell {
		return &tview.TableCell{
			Text:          title,
			MaxWidth:      width,
			Color:         tcell.ColorYellowGreen,
			NotSelectable: true,
			Expansion: 1,
		}
	}
	d.top = tview.NewTable()
	d.top.SetTitle("Current Planes")
	d.top.SetBorders(true).
		SetCell(0, 0, hdr("ICAO", 9)).
		SetCell(0, 1, hdr("Ident", 9)).
		SetCell(0, 2, hdr("Squawk", 6)).
		SetCell(0, 3, hdr("Altitude", 12)).
		SetCell(0, 4, hdr("Speed", 12)).
		SetCell(0, 5, hdr("Heading", 12)).
		SetCell(0, 6, hdr("# Msgs", 6)).
		SetCell(0, 7, hdr("Age (s)", 7))

	d.bottom = tview.NewTextView()
	d.bottom.SetBorder(true).SetTitle("Logs")

	_, _, _, height := d.bottom.GetRect()
	d.bottom.SetWordWrap(false).SetMaxLines(height + 50)
	d.bottom.SetChangedFunc(func() {
		d.app.Draw()
	})
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(d.top, 0, 1, true).
		AddItem(d.bottom, 0, 1, false)
	d.app.SetRoot(flex, true).EnableMouse(false)

	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyEscape {
			// we need to exit!
			d.app.Stop()
		}
		return event
	})

	return &d, nil
}

func (d *display) OnEvent(e tracker.Event) {
	d.app.QueueUpdate(func() {
		_, _, _, height := d.bottom.GetRect()
		d.bottom.SetWordWrap(false).SetMaxLines(height)
	})
	switch e.(type) {
	case *tracker.LogEvent:
		_, _ = fmt.Fprintln(d.bottom, e)
		d.bottom.ScrollToEnd()
	case *tracker.PlaneLocationEvent:
	case *tracker.FrameEvent:
		_, _ = fmt.Fprintln(d.bottom, e)
		d.bottom.ScrollToEnd()
	}
}
