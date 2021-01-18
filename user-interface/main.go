package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"./port_toolbox"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
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

func write(text string, file *os.File) {
	if _, err := file.WriteString(text); err != nil {
		panic(err)
	}
}

func read(filename string) string {
	data, err := ioutil.ReadFile(filename)
	check(err)
	return string(data)
}

func main() {
	app := app.New()

	w := app.NewWindow("List ")

	label1 := widget.NewLabel("Voici la liste ")

	b1 := widget.NewButton("Button1", func() {
		out, err := exec.Command("/bin/sh", "myProcess.sh").Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(out)
	})

	b2 := widget.NewButton("Button2",

		func() {
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
			file, err := os.OpenFile("Port.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)

			defer file.Close()
			check(err)

			write(numpor.Text, file)

			data := read(file.Name())
			fmt.Print(data)
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Port ajoué: " + numpor.Text,
				Content: largeText.Text,
			})
			file, err := os.OpenFile("Port.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
			defer file.Close()

			reader := bufio.NewReader(file)
			st, e := Readln(reader)
			for e == nil {
				fmt.Println(st)
				st, e = Readln(reader)

				if strings.Compare(largeText.Text, st) == 0 {
					fyne.CurrentApp().SendNotification(&fyne.Notification{
						Title:   "Pas de port saisie " + numpor.Text,
						Content: numpor.Text,
					})
					fmt.Println("passe par la ")
				} else {
					fmt.Println("passe par ici ")
					write(numpor.Text+"\n", file)
				}

			}

			check(err)

			data := read(file.Name())
			fmt.Print(data)
		},
	}

	label2 := widget.NewLabel("Label3")

	list := widget.NewVBox(

		widget.NewLabel("Item 1"),

		widget.NewLabel("Item 2"),
	)

	bar := widget.NewToolbar(

		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {}),
	)

	w.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewVBoxLayout(),
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), label1, layout.NewSpacer(), makeTable(
				[]string{"User", "PID", "App", ""},
				[][]string{{"1", "2", "3"}, {"4", "5", "6"}},
			), layout.NewSpacer()),
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), b1, b2, layout.NewSpacer()),

			fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), formPort, layout.NewSpacer()),

			fyne.NewContainerWithLayout(layout.NewBorderLayout(bar, label2, nil, nil), bar, list), layout.NewSpacer(),

			/*fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), label2,
			widget.NewCheck("Optional", func(value bool) {
				log.Println("Check set to", value)
			}),
			layout.NewSpacer(),
			widget.NewRadio([]string{"Option 1", "Option 2"}, func(value string) {
				log.Println("Radio set to", value)
			}),
			layout.NewSpacer(),
			widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {
				log.Println("Select set to", value)
			}),
			layout.NewSpacer()),*/
		),
	)

	w.Resize(fyne.NewSize(900, 200))
	w.ShowAndRun()

}

func makeTable(headings []string, rows [][]string) *widget.Box {

	columns := rowsToColumns(headings, rows)

	objects := make([]fyne.CanvasObject, len(columns))
	for k, col := range columns {
		box := widget.NewVBox(widget.NewLabelWithStyle(headings[k], fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		for _, val := range col {
			box.Append(widget.NewLabel(val))
		}
		objects[k] = box
	}
	return widget.NewHBox(objects...)
}

func rowsToColumns(headings []string, rows [][]string) [][]string {
	columns := make([][]string, len(headings))
	for _, row := range rows {
		for colK := range row {
			columns[colK] = append(columns[colK], row[colK])
		}
	}
	return columns
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
