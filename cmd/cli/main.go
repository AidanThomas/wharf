package main

import (
	"time"

	"github.com/AidanThomas/wharf/internal/providers/docker"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type colours struct {
	bg tcell.Color
	fg tcell.Color
}

func main() {
	defaultTheme := colours{
		bg: tcell.ColorNone,
		fg: tcell.ColorNone,
	}

	app := tview.NewApplication()

	dClient, err := docker.NewClient()
	if err != nil {
		panic(err)
	}
	defer dClient.Close()

	containers, err := dClient.GetAll()
	if err != nil {
		panic(err)
	}

	table := tview.NewTable()
	table.SetBackgroundColor(defaultTheme.bg)
	table.SetFixed(1, 0)
	table.SetBorder(true)
	table.SetSelectable(true, false)
	table.Select(1, 0)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			row, _ := table.GetSelection()
			id := table.GetCell(row, 4).Text
			ctr, err := dClient.GetById(id)
			if err != nil {
				panic(err)
			}
			switch ctr.State {
			case docker.Running:
				dClient.StopContainer(id)
			case docker.Exited:
				dClient.StartContainer(id)
			}
			go func() {
				time.Sleep(500 * time.Millisecond)
				containers, err := dClient.GetAll()
				if err != nil {
					panic(err)
				}
				drawTable(table, containers)
				app.Draw()
			}()
		}
		return event
	})

	drawTable(table, containers)

	flex := tview.NewFlex()
	flex.AddItem(table, 0, 1, true)

	if err := app.SetRoot(flex, true).EnableMouse(true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

// Image, Names, Status, Ports, ID
func getColour(ctr docker.Container) tcell.Color {
	switch ctr.State {
	case docker.Running:
		return tcell.ColorGreen
	case docker.Exited:
		return tcell.ColorRed
	default:
		return tcell.ColorYellow
	}
}

func drawTable(tbl *tview.Table, containers []docker.Container) {
	headings := []string{"Image", "Names", "Status", "Ports", "ID"}
	for i, h := range headings {
		tbl.SetCell(0, i, tview.NewTableCell(h).SetExpansion(1).SetSelectable(false))
	}
	tbl.Clear()

	for i, ctr := range containers {
		color := getColour(ctr)
		data := []string{ctr.Image, ctr.Names, ctr.Status, ctr.Ports, ctr.ID}
		for j, d := range data {
			tbl.SetCell(i+1, j, tview.NewTableCell(d).SetTextColor(color).SetExpansion(1))
		}
	}
}
