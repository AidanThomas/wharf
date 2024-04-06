package tui

import (
	"github.com/AidanThomas/wharf/internal/providers/docker"
	"github.com/rivo/tview"
)

type Table struct {
	*tview.Table
}

func NewContainersTable(theme Theme) *Table {
	tbl := &Table{tview.NewTable()}
	tbl.SetBackgroundColor(theme.Bg)
	tbl.SetFixed(1, 0)
	tbl.SetBorder(true)
	tbl.SetSelectable(true, false)

	return tbl
}

func (tbl *Table) DrawTable(containers []docker.Container) {
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
