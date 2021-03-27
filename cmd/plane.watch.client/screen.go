package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"plane.watch/lib/tracker"
	"sort"
	"sync"
	"time"
)

type (
	display struct {
		app    *tview.Application
		top    *tview.Table
		bottom *tview.TextView

		planes sync.Map
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
		defer c.Stop()
		for {
			select {
			case <- c.C:
				d.OnEvent(newPacer())
			}
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
	d.top.
		SetCell(0, 0, hdr("ICAO", 9)).
		SetCell(0, 1, hdr("Ident", 9)).
		SetCell(0, 2, hdr("squawk", 6)).
		SetCell(0, 3, hdr("Lat/Lon", 19)).
		SetCell(0, 4, hdr("altitude", 12)).
		SetCell(0, 5, hdr("Speed", 12)).
		SetCell(0, 6, hdr("heading", 12)).
		SetCell(0, 7, hdr("# Msgs", 6)).
		SetCell(0, 8, hdr("Age (s)", 7)).
		SetCell(0, 9, hdr("Extra", 10)).
		SetBorder(true)

	d.bottom = tview.NewTextView().SetDynamicColors(true)
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

func (d *display) drawTable() {
	table := make(map[uint32]*tracker.Plane)
	icaos := make([]uint32,0)
	d.planes.Range(func(key, value interface{}) bool {
		table[key.(uint32)] = value.(*tracker.Plane)
		icaos  = append(icaos, key.(uint32))
		return true
	})

	sort.Slice(icaos, func(i, j int) bool {
		return icaos[i] < icaos[j]
	})

	row := 1
	var latLon string
	for _, icao := range icaos {
		d.top.SetCellSimple(row, 0, table[icao].IcaoIdentifierStr())
		d.top.SetCellSimple(row, 1, table[icao].FlightNumber())
		d.top.SetCellSimple(row, 2, table[icao].SquawkIdentityStr())
		if table[icao].HasLocation() {
			latLon = fmt.Sprintf("%0.4f/%0.4f", table[icao].Lat(), table[icao].Lon())
		} else {
			latLon = "?"
		}
		d.top.SetCellSimple(row, 3, latLon)
		d.top.SetCellSimple(row, 4, fmt.Sprint(table[icao].Altitude()))
		d.top.SetCellSimple(row, 5, fmt.Sprintf("%0.2f",table[icao].Velocity()))
		d.top.SetCellSimple(row, 6, table[icao].HeadingStr())
		d.top.SetCellSimple(row, 7, fmt.Sprint(table[icao].MsgCount()))

		since := time.Now().Sub(table[icao].LastSeen()).Seconds()
		d.top.SetCellSimple(row, 8, fmt.Sprintf("%0.0f",since))
		d.top.SetCellSimple(row, 9, table[icao].Special())

		row++
	}
}

func (d *display) OnEvent(e tracker.Event) {
	d.app.QueueUpdate(func() {
		_, _, _, height := d.bottom.GetRect()
		d.bottom.SetWordWrap(false).SetMaxLines(height)
	})
	switch e.(type) {
	case *tracker.LogEvent:
		w := tview.ANSIWriter(d.bottom)
		_, _ = fmt.Fprintln(w, e)
		d.bottom.ScrollToEnd()

	case *tracker.PlaneLocationEvent:
		ple := e.(*tracker.PlaneLocationEvent)
		if ple.Removed() {
			d.planes.Delete(ple.Plane().IcaoIdentifier())
		} else {
			d.planes.Store(ple.Plane().IcaoIdentifier(), ple.Plane())
		}
		//d.drawTable()

	case *tracker.FrameEvent:
		// show the received frame
	case *pacerEvent:
		// cleanup our planes list
		d.drawTable()
	}
}
