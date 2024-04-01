package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

	ptrFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	ctrFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	ctrFlex.SetBackgroundColor(defaultTheme.bg)
	ctrFlex.SetTitle("Containers").SetTitleAlign(tview.AlignLeft)

	mainFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainFlex.SetBorder(true)
	mainFlex.SetBackgroundColor(defaultTheme.bg)
	mainFlex.
		AddItem(ptrFlex, 3, 0, true).
		AddItem(ctrFlex, 0, 1, false)

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		ptrFlex.AddItem(createPointerEntry(defaultTheme), 1, 0, false)
		ctrFlex.AddItem(createContainerEntry(ctr, defaultTheme), 1, 0, false)
	}

	app.SetRoot(mainFlex, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func createContainerEntry(ctr types.Container, theme colours) *tview.TextView {
	var ctrString string
	switch ctr.State {
	case "running":
		ctrString = fmt.Sprintf("[green]%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.State)
	case "exited":
		ctrString = fmt.Sprintf("[red]%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.State)
	default:
		ctrString = fmt.Sprintf("[yellow]%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.State)
	}

	view := tview.NewTextView()
	view.SetTextAlign(tview.AlignLeft)
	view.SetBackgroundColor(theme.bg)
	view.SetDynamicColors(true)
	view.SetBorder(false)
	view.SetText(ctrString)
	return view
}

func createPointerEntry(theme colours) *tview.TextView {
	view := tview.NewTextView()
	view.SetTextAlign(tview.AlignLeft)
	view.SetBackgroundColor(theme.bg)
	view.SetDynamicColors(true)
	view.SetBorder(false)
	view.SetText(" > ")
	return view
}
