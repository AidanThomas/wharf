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
	table := tview.NewTable()
	table.SetBackgroundColor(defaultTheme.bg)

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	headings := []string{"Image", "Names", "Status", "Ports", "ID"}
	for i, h := range headings {
		table.SetCell(0, i+1, tview.NewTableCell(h))
	}

	for i, ctr := range containers {
		color, data := parseContainer(ctr)
		for j, d := range data {
			table.SetCell(i+1, j+1, tview.NewTableCell(d).SetTextColor(color))
		}
	}

	table.SetFixed(1, 0)
	table.SetBorder(true)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			table.SetSelectable(true, false)
			table.Select(1, 0)
		}
		return event
	})

	flex := tview.NewFlex()
	flex.AddItem(table, 0, 1, true)

	if err := app.SetRoot(flex, true).EnableMouse(true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

// Image, Names, Status, Ports, ID
func parseContainer(ctr types.Container) (tcell.Color, []string) {
	var color tcell.Color
	switch ctr.State {
	case "running":
		color = tcell.ColorGreen
	case "exited":
		color = tcell.ColorRed
	default:
		color = tcell.ColorYellow
	}

	var ports []string
	for _, p := range ctr.Ports {
		var port string
		if p.IP != "" {
			port = fmt.Sprintf("%s:%d->%d/%s", p.IP, p.PublicPort, p.PrivatePort, p.Type)
		} else {
			port = fmt.Sprintf("%d/%s", p.PrivatePort, p.Type)
		}
		ports = append(ports, port)
	}

	image := ctr.Image
	names := strings.Join(ctr.Names, ", ")
	status := ctr.Status
	id := ctr.ID

	return color, []string{image, names, status, strings.Join(ports, ", "), id}
}
