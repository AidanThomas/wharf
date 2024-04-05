package main

import "github.com/AidanThomas/wharf/internal/tui"

func main() {
	app, err := tui.NewTui()
	if err != nil {
		panic(err)
	}
	defer app.Close()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
