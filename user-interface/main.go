package main

import (
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var topWindow fyne.Window

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func main() {
	//creation de l'application
	myApp := app.New()
	myWindow := myApp.NewWindow("NetCop")

	//toolbar
	toolbar := createToolbar()

	ports := fileToSlice("Port.txt")




	//Port
	left := widget.NewTabContainer()
	left.SetTabLocation(widget.TabLocationLeading)

	updatePortlist(ports, left)

	//form
	entry := widget.NewEntry()
	form := createForm(left, entry)
	form.Append("Port to block:", entry)

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

func updatePortlist(tabPort []string, item *widget.TabContainer) {
	for _, port := range tabPort {

		button := widget.NewButton("Delete", func() {

			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title: "Port retiré: " + item.Items[item.CurrentTabIndex()].Text})
			delete_port("Port.txt", item.Items[item.CurrentTabIndex()].Text)
			log.Println(item.Items[item.CurrentTabIndex()].Text)

			//update

			//tabPort = updateTab(tabPort,item.Items[item.CurrentTabIndex()].Text)
			item.Remove(item.Items[item.CurrentTabIndex()])
			item.Refresh()

		})
		encap := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), button)

		item.Append(widget.NewTabItem(port, encap))
	}
}

//update for adding
func addUIPort(item *widget.TabContainer,port string)  {
	button := widget.NewButton("Delete", func() {

		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title: "Port retiré: " + item.Items[item.CurrentTabIndex()].Text})
		delete_port("Port.txt", item.Items[item.CurrentTabIndex()].Text)
		log.Println(item.Items[item.CurrentTabIndex()].Text)

		//update on remove
		item.Remove(item.Items[item.CurrentTabIndex()])
		item.Refresh()

	})
	encap := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), button)

	item.Append(widget.NewTabItem(port, encap))
}
//gestion formulaire
func createForm(item *widget.TabContainer, entry *widget.Entry) *widget.Form {
	Form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
		},
		OnSubmit: func() { // optional, handle form submission
			fmt.Println("Form submitted")
			fmt.Println(entry.Text)
			if doublonPort("Port.txt", entry.Text) == true {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title: "Port déja existant: " + entry.Text,
				})
			} else if len(entry.Text) == 0 {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title: "There is no port in the field !",

				})
			}else {
				AddPort("Port.txt", entry.Text)
				addUIPort(item,entry.Text)
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title: "Port ajoué: " + entry.Text,

				})
			}
		},
	}
	return Form
}

//the top toolbar
func createToolbar() *widget.Toolbar {
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
func updateTab(tabs []string, port string) []string {
	var result  []string
	for _,content := range tabs{
		if strings.Compare(content, port) != 0{
			result = append(result, content)
		}
	}
	return result
}

/*b1 := widget.NewButton("Script Process", func() {
	cmd := exec.Command("/bin/sh", "../myProcess.sh")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
})

b2 := widget.NewButton("List process", func() {
	processList, err := ps.Processes()
	if err != nil {
		log.Println("ps.Processes() Failed, are you using windows?")
		return
	}

	infoStat, _ := host.Info()
	fmt.Printf("Total processes: %d\n", infoStat.Procs)

	miscStat, _ := load.Misc()
	fmt.Printf("Running processes: %d\n", miscStat.ProcsRunning)

	for x := range processList {
		var process ps.Process
		process = processList[x]
		log.Printf("%d\t%s\t%d\n", process.Pid(), process.Executable(), process.PPid())
		//log.Printf("%d\t%s\n", process.Pid(), process.Executable())

	}
})*/
