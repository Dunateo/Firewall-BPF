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
	"github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//const portFile = "/usr/lib/block_rules/blockedPort"
const procesFile = "/usr/lib/block_rules/myProcess.sh"

const portFile = "Port.txt"
//const procesFile = "myProcess.sh"

func main() {
	//creation de l'application
	myApp := app.New()
	myWindow := myApp.NewWindow("NetCop")

	//toolbar
	toolbar := createToolbar()

	//ports file
	ports := fileToSlice(portFile)

	//prog file
	launchedProg := fileToSlice("/home/app.txt")

	//prog
	prog := widget.NewAccordionContainer()
	createTab(prog, launchedProg)

	flex := widget.NewScrollContainer(prog)

	//Port
	left := widget.NewTabContainer()
	left.SetTabLocation(widget.TabLocationLeading)

	updatePortlist(ports, left)

	//form
	entry := widget.NewEntry()
	//entry.Resize(40,3)
	form := createForm(left, entry)
	form.Append("Port to block:", entry)

	//center
	center := fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		form)

	//Hbox layout
	formProg := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(2),
		center, flex)

	//the border big one
	content := fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbar, nil, left, nil),
		toolbar, left, formProg)

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
			delete_port(portFile, item.Items[item.CurrentTabIndex()].Text)
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
		delete_port(portFile, item.Items[item.CurrentTabIndex()].Text)
		log.Println(item.Items[item.CurrentTabIndex()].Text)

		//update on remove
		item.Remove(item.Items[item.CurrentTabIndex()])
		item.Refresh()

	})
	encap := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), button)

	item.Append(widget.NewTabItem(port, encap))
}

//gestion prog display
func createTab(item *widget.AccordionContainer, tabProg []string) {
	//fmt.Println(tabProg)
	for _, prog := range tabProg {

		//fmt.Println(prog)

		item.Append(widget.NewAccordionItem(prog, widget.NewLabel("la")))
	}

}

//gestion de commande
func execCmd(command *exec.Cmd) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	err := command.Run()
	check(err)
	fmt.Println("Result: " + out.String())

}

//gestion formulaire
func createForm(item *widget.TabContainer, entry *widget.Entry) *widget.Form {
	Form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
		},
		OnSubmit: func() { // optional, handle form submission
			fmt.Println("Form submitted")
			fmt.Println(entry.Text)
			if doublonPort(portFile, entry.Text) == true {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title: "Port déja existant: " + entry.Text,
				})
			} else if len(entry.Text) == 0 {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title: "There is no port in the field !",
				})
			} else {
				AddPort(portFile, entry.Text)
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
		widget.NewToolbarAction(theme.FileApplicationIcon(), func() {
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
			}
		}),
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			cmd := exec.Command("/bin/sh", procesFile)
			execCmd(cmd)

		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
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
