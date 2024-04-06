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

var (
	mainFlex     *tview.Flex
	resultsTable *tview.Table
)

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
	mainFlex = tview.NewFlex()
	resultsTable = tview.NewTable()
	resultsTable.SetBackgroundColor(a.theme.Bg)
	resultsTable.SetFixed(1, 0)
	resultsTable.SetBorder(true)
	resultsTable.SetSelectable(true, false)
	resultsTable.Select(1, 0)
	resultsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		rune := event.Rune()
		switch key {
		case tcell.KeyEnter:
			row, _ := resultsTable.GetSelection()
			id := resultsTable.GetCell(row, 4).Text
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
				drawTable(resultsTable, containers)
				a.Draw()
			}()
		}
		switch rune {
		case '/':
			a.createSearchUI()
		}
		return event
	})
	containers, err := a.cli.GetAll()
	if err != nil {
		return err
	}

	drawTable(resultsTable, containers)

	mainFlex.AddItem(resultsTable, 0, 1, true)
	a.SetRoot(mainFlex, true).EnableMouse(true).SetFocus(mainFlex)

	return nil
}

func (a *Application) createSearchUI() {
	mainFlex.Clear()

	searchBox := tview.NewTextArea()
	searchBox.SetBorder(true)
	searchBox.SetBackgroundColor(a.theme.Bg)
	searchBox.SetPlaceholder("Search by name")
	searchBox.SetPlaceholderStyle(tcell.StyleDefault)
	searchBox.SetTextStyle(tcell.StyleDefault)

	searchBox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			searchTerm := searchBox.GetText()
			containers, err := a.cli.SearchByName(searchTerm)
			if err != nil {
				panic(err)
			}
			drawTable(resultsTable, containers)
			return nil
		}
		return event
	})

	mainFlex.SetDirection(tview.FlexRow)

	mainFlex.
		AddItem(searchBox, 3, 1, true).
		AddItem(resultsTable, 0, 1, false)

	a.SetFocus(searchBox)
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
