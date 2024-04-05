package main

import (
	"github.com/AidanThomas/wharf/internal/providers/docker"
	"github.com/AidanThomas/wharf/internal/tui"
)

func main() {
	dClient, err := docker.NewClient()
	if err != nil {
		panic(err)
	}
	defer dClient.Close()
	app, err := tui.NewTui(dClient)
	if err != nil {
		panic(err)
	}
	if err := app.Run(); err != nil {
		panic(err)
	}
}
