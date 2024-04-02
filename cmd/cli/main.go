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
		text := tview.NewTextView().SetText(getContainerText(ctr))
		text.SetBackgroundColor(defaultTheme.bg)
		text.SetDynamicColors(true)
		ctrFlex.
			AddItem(caret, 3, 0, false).
			AddItem(text, 0, 1, false)
		ctrFlex.SetFocusFunc(func() {
			caret.SetText(" > ")
		})
		ctrFlex.SetBlurFunc(func() {
			caret.SetText("   ")
		})
		flex.AddItem(ctrFlex, 1, 1, false)
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

func getContainerText(ctr types.Container) string {
	out := fmt.Sprintf("%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.State)
	switch ctr.State {
	case "running":
		out = fmt.Sprintf("[green]%s", out)
	case "exited":
		out = fmt.Sprintf("[red]%s", out)
	default:
		out = fmt.Sprintf("[yellow]%s", out)
	}

	return out
}
