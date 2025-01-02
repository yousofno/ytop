package main

import (
	"github.com/mitchellh/go-ps"
	"github.com/rivo/tview"
	"os"
	"strconv"
)

var pages = tview.NewPages()
var app = tview.NewApplication()

func contains(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

func ProcessListPage() *tview.Flex {
	var proccesDetail = tview.NewTextView()
	proccesDetail.SetBorder(true)
	proccesDetail.SetTitle("Process Detail")
	uniqueProcc := []string{}
	var flex = tview.NewFlex()
	var proccesList = tview.NewList()
	proccesList.SetBorder(true)
	proccesList.SetTitle("Process list")
	proccesNames, err := ps.Processes()
	if err == nil {
		for _, procces := range proccesNames {
			if !contains(uniqueProcc, procces.Executable()) {
				uniqueProcc = append(uniqueProcc, procces.Executable())
				proccesList.AddItem(procces.Executable(), "", '-', func() {
					var detail = "Pid : " + strconv.Itoa(procces.Pid()) + "\n" + "Parent proccess id : " + strconv.Itoa(procces.PPid())
					proccesDetail.SetText(detail)
				})
			}
		}
	}
	var helpList = tview.NewList()
	helpList.SetBorder(true)
	helpList.SetTitle("Help")
	flex.AddItem(proccesList, 0, 3, true)
	flex.AddItem(proccesDetail, 0, 6, false)
	return flex
}

func menuPage() *tview.List {
	var list = tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Y-TOP")
	list.AddItem("Process List", "Show All Procces lists", 'p', func() {
		pages.SwitchToPage("procc")
	})
	list.AddItem("Quit", "Quit", 'q', func() {
		os.Exit(0)
	})
	return list
}

func init() {
	var menu = menuPage()
	var procc = ProcessListPage()
	pages.AddPage("menu", menu, true, true)
	pages.AddPage("procc", procc, true, false)
	app.SetRoot(pages, true)
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
func main() {
}
