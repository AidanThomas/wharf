package main

import (
	"context"
	"fmt"
	"strings"

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

var (
	index = -1
)

func main() {
	defaultTheme := colours{
		bg: tcell.ColorNone,
		fg: tcell.ColorNone,
	}

	app := tview.NewApplication()

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBorder(true)
	flex.SetBackgroundColor(defaultTheme.bg)

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
		ctrFlex := tview.NewFlex()
		ctrFlex.SetBackgroundColor(defaultTheme.bg)
		caret := tview.NewTextView().SetText("   ")
		caret.SetBackgroundColor(defaultTheme.bg)
		text := createContainerFlex(ctr, defaultTheme)
		ctrFlex.
			AddItem(caret, 3, 0, false).
			AddItem(text, 0, 1, false)
		ctrFlex.SetFocusFunc(func() {
			caret.SetText(" > ")
		})
		ctrFlex.SetBlurFunc(func() {
			caret.SetText("   ")
		})
		flex.AddItem(ctrFlex, 3, 1, false)
	}

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			if index < flex.GetItemCount()-1 {
				index++
			} else {
				index = 0
			}
			app.SetFocus(flex.GetItem(index))
		}

		if event.Rune() == 'k' {
			if index > 0 {
				index--
			} else {
				index = flex.GetItemCount() - 1
			}
			app.SetFocus(flex.GetItem(index))
		}
		return event
	})

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

// Image, Names, Status, Ports, ID

func createContainerFlex(ctr types.Container, theme colours) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	flex.SetBackgroundColor(theme.bg)

	var colour string
	switch ctr.State {
	case "running":
		colour = "green"
	case "exited":
		colour = "red"
	default:
		colour = "yellow"
	}

	names := tview.NewTextView().SetDynamicColors(true)
	names.SetBackgroundColor(theme.bg)
	names.SetText(fmt.Sprintf("[%s]%s", colour, strings.Join(ctr.Names, ", ")))
	image := tview.NewTextView().SetDynamicColors(true)
	image.SetBackgroundColor(theme.bg)
	image.SetText(fmt.Sprintf("[%s]%s", colour, ctr.Image))
	status := tview.NewTextView().SetDynamicColors(true)
	status.SetBackgroundColor(theme.bg)
	status.SetText(fmt.Sprintf("[%s]%s", colour, ctr.Status))
	ports := tview.NewTextView().SetDynamicColors(true)
	ports.SetBackgroundColor(theme.bg)
	ports.SetText(fmt.Sprintf("[%s]%s", colour, "ports"))
	id := tview.NewTextView().SetDynamicColors(true)
	id.SetBackgroundColor(theme.bg)
	id.SetText(fmt.Sprintf("[%s]%s", colour, ctr.ID))

	flex.
		AddItem(names, 0, 1, false).
		AddItem(image, 0, 1, false).
		AddItem(status, 0, 1, false).
		AddItem(ports, 0, 1, false).
		AddItem(id, 0, 1, false)

	return flex
}
