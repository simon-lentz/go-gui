/*
markdown tests out fyne's native
markdown support.
*/
package markdown

import (
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI // fyne's URI system allows for the platform agnostic OS ops
	SaveMenuItem  *fyne.MenuItem
}

var cfg config

func Markdown() {
	// create a fyne app
	a := app.New()

	// theme can be customized as shown in the theme.go file
	a.Settings().SetTheme(&myTheme{})
	// create a window for the app
	win := a.NewWindow("markdown")

	// get the user interface
	edit, preview := cfg.makeUI()
	cfg.createMenuItems(win)

	// set the content of the window
	win.SetContent(container.NewHSplit(edit, preview))

	// show window and run app
	win.Resize(fyne.Size{Width: 800, Height: 500})
	win.CenterOnScreen()
	win.ShowAndRun()
}

func (app *config) makeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")

	app.EditWidget = edit
	app.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}

func (app *config) createMenuItems(win fyne.Window) {
	// create menu items
	openMenuItem := fyne.NewMenuItem("Open...", app.openFunc(win))
	saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(win))

	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true
	saveAsMenuItem := fyne.NewMenuItem("Save As...", app.saveAsFunc(win))

	// create file menu, add items
	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)

	// create and set main menu (top bar)
	menu := fyne.NewMainMenu(fileMenu)
	win.SetMainMenu(menu)
}

func (app *config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				// the dialog is tied to the specific window parameter it is passed, prevents dupes
				dialog.ShowError(err, win)
				return
			}

			if write == nil {
				// user cancels
				return
			}

			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "Only .md type files are supported!", win)
			}
			// save file
			write.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = write.URI()

			defer write.Close()

			// change filename header to reflect name entered on save
			win.SetTitle(win.Title() + " - " + write.URI().Name())
			app.SaveMenuItem.Disabled = false

		}, win)
		saveDialog.SetFileName("untitled.md") // set temp filename
		saveDialog.SetFilter(filter)          // prevent creation of anything other than .md
		saveDialog.Show()                     // after the write closer is associated with the window and func behavior it can be called

	}
}

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (app *config) openFunc(win fyne.Window) func() {
	return func() {
		// New vs Show allows for addition of a filter so that only .md files can be opened
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			if read == nil {
				//user cancelled
				return
			}

			defer read.Close()

			data, err := io.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			app.EditWidget.SetText(string(data))

			app.CurrentFile = read.URI()
			win.SetTitle(win.Title() + " - " + read.URI().Name())
			app.SaveMenuItem.Disabled = false

		}, win)

		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

func (app *config) saveFunc(win fyne.Window) func() {
	return func() {
		if app.CurrentFile != nil {
			write, err := storage.Writer(app.CurrentFile)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			write.Write([]byte(app.EditWidget.Text))
			defer write.Close()
		}
	}
}
