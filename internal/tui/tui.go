package tui

import (
	"time"

	"github.com/AidanThomas/wharf/internal/providers/docker"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Application struct {
	*tview.Application
	theme Theme
	cli   *docker.Client
}

func NewTui() (*Application, error) {
	defaultTheme := Theme{
		Bg: tcell.ColorNone,
		Fg: tcell.ColorNone,
	}

	dCli, err := docker.NewClient()
	if err != nil {
		return &Application{}, err
	}

	app := &Application{
		Application: tview.NewApplication(),
		theme:       defaultTheme,
		cli:         dCli,
	}

	if err := app.createUI(); err != nil {
		return &Application{}, err
	}

	return app, nil
}

func (a *Application) Close() {
	a.cli.Close()
}

func (a *Application) createUI() error {
	flx := tview.NewFlex()
	tbl := tview.NewTable()
	tbl.SetBackgroundColor(a.theme.Bg)
	tbl.SetFixed(1, 0)
	tbl.SetBorder(true)
	tbl.SetSelectable(true, false)
	tbl.Select(1, 0)
	tbl.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			row, _ := tbl.GetSelection()
			id := tbl.GetCell(row, 4).Text
			ctr, err := a.cli.GetById(id)
			if err != nil {
				panic(err)
			}
			switch ctr.State {
			case docker.Running:
				a.cli.StopContainer(id)
			case docker.Exited:
				a.cli.StartContainer(id)
			}
			go func() {
				time.Sleep(500 * time.Millisecond)
				containers, err := a.cli.GetAll()
				if err != nil {
					panic(err)
				}
				drawTable(tbl, containers)
				a.Draw()
			}()
		}
		return event
	})
	containers, err := a.cli.GetAll()
	if err != nil {
		return err
	}

	drawTable(tbl, containers)

	flx.AddItem(tbl, 0, 1, true)
	a.SetRoot(flx, true).EnableMouse(true).SetFocus(flx)

	return nil
}

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
	tbl.Clear()

	headings := []string{"Image", "Names", "Status", "Ports", "ID"}
	for i, h := range headings {
		tbl.SetCell(0, i, tview.NewTableCell(h).SetExpansion(1).SetSelectable(false))
	}

	for i, ctr := range containers {
		color := getColour(ctr)
		data := []string{ctr.Image, ctr.Names, ctr.Status, ctr.Ports, ctr.ID}
		for j, d := range data {
			tbl.SetCell(i+1, j, tview.NewTableCell(d).SetTextColor(color).SetExpansion(1))
		}
	}
}
