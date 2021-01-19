package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"regexp"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/cmd/fyne_settings/settings"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"

	ps "github.com/mitchellh/go-ps"
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
	app := app.New()

	w := app.NewWindow("List ")

	//Top bar

	newItem := fyne.NewMenuItem("New", nil)
	otherItem := fyne.NewMenuItem("Other", nil)
	otherItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Project", func() { fmt.Println("Menu New->Other->Project") }),
		fyne.NewMenuItem("Mail", func() { fmt.Println("Menu New->Other->Mail") }),
	)
	newItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("File", func() { fmt.Println("Menu New->File") }),
		fyne.NewMenuItem("Directory", func() { fmt.Println("Menu New->Directory") }),
		otherItem,
	)
	settingsItem := fyne.NewMenuItem("Settings", func() {
		w := app.NewWindow("Fyne Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(480, 480))
		w.Show()
	})

	cutItem := fyne.NewMenuItem("Cut", func() {
		shortcutFocused(&fyne.ShortcutCut{
			Clipboard: w.Clipboard(),
		}, w)
	})
	copyItem := fyne.NewMenuItem("Copy", func() {
		shortcutFocused(&fyne.ShortcutCopy{
			Clipboard: w.Clipboard(),
		}, w)
	})
	pasteItem := fyne.NewMenuItem("Paste", func() {
		shortcutFocused(&fyne.ShortcutPaste{
			Clipboard: w.Clipboard(),
		}, w)
	})
	findItem := fyne.NewMenuItem("Find", func() { fmt.Println("Menu Find") })

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://developer.fyne.io")
			_ = app.OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://fyne.io/support/")
			_ = app.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sponsor", func() {
			u, _ := url.Parse("https://github.com/sponsors/fyne-io")
			_ = app.OpenURL(u)
		}))
	mainMenu := fyne.NewMainMenu(
		// a quit item will be appended to our first menu
		fyne.NewMenu("File", newItem, fyne.NewMenuItemSeparator(), settingsItem),
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
		helpMenu,
	)
	w.SetMainMenu(mainMenu)
	w.SetMaster()

	//label1 := widget.NewLabel("Voici la liste ")
	//label2 := widget.NewLabel("Label3")

	b1 := widget.NewButton("Script Process", func() {
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

		/*out, err := exec.Command("/bin/sh", "myProcess.sh").Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(out)*/
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
	})

	re := regexp.MustCompile("[0-9]+")
	fmt.Println(re.FindAllString("abc123def987asdf", -1))

	numpor := widget.NewEntry()
	numpor.SetPlaceHolder("Port")

	largeText := widget.NewMultiLineEntry()

	formPort := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Port", Widget: numpor},
		},
		OnCancel: func() {
			fmt.Println("Remove")
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Port retiré: " + numpor.Text,
				Content: largeText.Text,
			})
			fmt.Println(numpor.Text)
			delete_port("Port.txt", numpor.Text)

		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			fmt.Println(numpor.Text)
			if doublonPort("Port.txt", numpor.Text) == true {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Port déja existant: " + numpor.Text,
					Content: largeText.Text,
				})
			} else {
				AddPort("Port.txt", numpor.Text)
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Port ajoué: " + numpor.Text,
					Content: largeText.Text,
				})
			}
		},
	}

	w.SetContent(

		fyne.NewContainerWithLayout(

			layout.NewVBoxLayout(),

			//fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), label1 ,label2, layout.NewSpacer()),

			fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), b1, b2, layout.NewSpacer()),

			fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), formPort, layout.NewSpacer()),
		),
	)

	w.Resize(fyne.NewSize(900, 200))
	w.ShowAndRun()

}
