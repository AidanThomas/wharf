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
	resultsTable *Table
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

	mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		rune := event.Rune()
		switch rune {
		case '/':
			a.drawSearch()
			return nil
		}
		switch key {
		case tcell.KeyEsc:
			if err := a.drawDefault(); err != nil {
				panic(err)
			}
		}
		return event
	})

	resultsTable = NewContainersTable(a.theme)
	resultsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
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
				resultsTable.DrawTable(containers)
				a.Draw()
			}()
		}
		return event
	})
	containers, err := a.cli.GetAll()
	if err != nil {
		return err
	}

	resultsTable.DrawTable(containers)

	mainFlex.AddItem(resultsTable, 0, 1, true)
	a.SetRoot(mainFlex, true).EnableMouse(true).SetFocus(mainFlex)

	return nil
}

func (a *Application) drawDefault() error {
	mainFlex.Clear()

	resultsTable.Clear()
	containers, err := a.cli.GetAll()
	if err != nil {
		return err
	}

	mainFlex.AddItem(resultsTable, 0, 1, true)
	resultsTable.DrawTable(containers)
	a.SetFocus(resultsTable)

	return nil
}

func (a *Application) drawSearch() {
	mainFlex.Clear()

	searchBox := tview.NewTextArea()
	searchBox.SetBorder(true)
	searchBox.SetBackgroundColor(a.theme.Bg)
	searchBox.SetPlaceholder("Search by name")
	searchBox.SetPlaceholderStyle(tcell.StyleDefault)
	searchBox.SetTextStyle(tcell.StyleDefault)

	searchBox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			searchTerm := searchBox.GetText()
			containers, err := a.cli.SearchByName(searchTerm)
			if err != nil {
				panic(err)
			}
			resultsTable.DrawTable(containers)
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
