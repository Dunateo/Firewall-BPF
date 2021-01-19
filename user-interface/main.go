

package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	_ "fyne.io/fyne/widget"
	"log"
)

func main() {
	//creation de l'application
	myApp := app.New()
	myWindow := myApp.NewWindow("NetCop")

	//toolbar
	toolbar := createToolbar()

	//form
	entry := widget.NewEntry()
	form := createForm(myWindow, entry)
	form.Append("Port to block:", entry)

	//Port
	left := widget.NewTabContainer()
	left.SetTabLocation(widget.TabLocationLeading)
	ports := []string{"10", "20", "30", "443"}
	updatePortlist(ports,left)

	//center layout
	center := fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		form)

	//the border big one
	content := fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbar, nil, left, nil),
		toolbar, left, center)

	//set what will be in the window
	myWindow.SetContent(content)

	myWindow.Resize(fyne.NewSize(720, 576))
	myWindow.SetFixedSize(true)
	//run the app
	myWindow.ShowAndRun()
}

func updatePortlist(tabPort []string, item *widget.TabContainer)  {
	for _, port := range tabPort {

		button := widget.NewButton("Delete", func() {
			log.Println(item.Items[item.CurrentTabIndex()].Text)

		})
		encap := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), button)

		item.Append(widget.NewTabItem(port,encap ))
	}
}
//gestion formulaire
func createForm(myWindow fyne.Window, entry *widget.Entry)*widget.Form  {
	Form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			},
		OnSubmit: func() { // optional, handle form submission
			log.Println("Form submitted:", entry.Text)

		},
	}
	return Form
}

//the top toolbar
func createToolbar() *widget.Toolbar{
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("New document")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)
	return toolbar
}


