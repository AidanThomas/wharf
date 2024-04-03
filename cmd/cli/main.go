package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
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

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
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
			filters := filters.NewArgs()
			filters.Add("id", id)
			filtered, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true, Filters: filters})
			if err != nil {
				panic(err)
			}
			switch filtered[0].State {
			case "running":
				apiClient.ContainerStop(context.Background(), id, container.StopOptions{})
			case "exited":
				apiClient.ContainerStart(context.Background(), id, container.StartOptions{})
			}
			containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
			if err != nil {
				panic(err)
			}
			table.Clear()
			drawTable(table, containers)
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

		if len(ports) > 3 {
			ports = ports[:3]
			ports = append(ports, "...")
		}
	}

	image := ctr.Image
	names := strings.Join(ctr.Names, ", ")
	status := ctr.Status
	id := ctr.ID

	return color, []string{image, names, status, strings.Join(ports, ", "), id}
}

func drawTable(tbl *tview.Table, containers []types.Container) {

	headings := []string{"Image", "Names", "Status", "Ports", "ID"}
	for i, h := range headings {
		tbl.SetCell(0, i, tview.NewTableCell(h).SetExpansion(1).SetSelectable(false))
	}

	for i, ctr := range containers {
		color, data := parseContainer(ctr)
		for j, d := range data {
			tbl.SetCell(i+1, j, tview.NewTableCell(d).SetTextColor(color).SetExpansion(1))
		}
	}
}
