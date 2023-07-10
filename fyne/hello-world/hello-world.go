/*
Fyne is a platform agnostic GUI framework
built in and for Go.
*/
package hello_world

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	output *widget.Label
}

var testApp App

func HelloWorld() {
	a := app.New()
	w := a.NewWindow("Hello, World!")

	w.SetContent(widget.NewLabel("Hello, World!"))

	output, entry, btn := testApp.makeUI()

	// creates a text entry box and button in the app window
	w.SetContent(container.NewVBox(output, entry, btn))

	// Resize sets the default size of the window so that it's not tiny.
	w.Resize(fyne.Size{Width: 500, Height: 500})
	// end point for code used in application
	w.ShowAndRun()
	// past this point no code will be called
	// until the application (window) is quit
}

func (app *App) makeUI() (*widget.Label, *widget.Entry, *widget.Button) {
	output := widget.NewLabel("Hello, World!")
	entry := widget.NewEntry()
	btn := widget.NewButton("Enter", func() {
		app.output.SetText(entry.Text)
	})
	// this changes the color of the button, one of the only ways this can
	// be done without signficant effort
	btn.Importance = widget.HighImportance
	app.output = output

	return output, entry, btn
}
