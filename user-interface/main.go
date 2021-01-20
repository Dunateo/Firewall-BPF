package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
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

	//button cmd
	cmd := exec.Command("/bin/sh", "../myProcess.sh")
	buttonCmd := createButtonWithCmd(left, cmd)

	//button process
	/*processList, err := go-ps.Processes()
	check(err)
	buttonProc := createButtonProcess(left, processList)*/

	//center layout
	center := fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		form)

	//right layout
	right := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		buttonCmd)

	//the border big one
	content := fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbar, nil, left, nil),
		toolbar, left, center, right)

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
func addUIPort(item *widget.TabContainer, port string) {
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

//gestion de commande
func createButtonWithCmd(item *widget.TabContainer, command *exec.Cmd) *widget.Button {
	Button := widget.NewButton("Script Process", func() {
		var out bytes.Buffer
		var stderr bytes.Buffer
		command.Stdout = &out
		command.Stderr = &stderr
		err := command.Run()
		check(err)
		fmt.Println("Result: " + out.String())
	})
	return Button
}

/*func createButtonProcess(item *widget.TabContainer, processList *ps.Process) *widget.Button {
	Button := widget.NewButton("Script Process", func() {



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
	})
	return Button
}*/

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
			} else {
				AddPort("Port.txt", entry.Text)
				addUIPort(item, entry.Text)
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
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ColorChromaticIcon(), func() {

		}),
	)
	return toolbar
}

//update port tab
func updateTab(tabs []string, port string) []string {
	var result []string
	for _, content := range tabs {
		if strings.Compare(content, port) != 0 {
			result = append(result, content)
		}
	}
	return result
}
