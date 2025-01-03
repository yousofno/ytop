package main

import (
	"github.com/mitchellh/go-ps"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/mem"
	"os"
	"strconv"
	"time"
)

var pages = tview.NewPages()
var app = tview.NewApplication()
var quitChannel = make(chan bool)

func ProcessListPage() *tview.Flex {
	var flex = tview.NewFlex()
	var treeView = tview.NewTreeView()
	var root = tview.NewTreeNode(".")
	root.SetExpanded(true)
	treeView.SetRoot(root)
	treeView.SetBorder(true)
	treeView.SetTitle("Procces list")

	var onlineUsage = tview.NewTextView()
	onlineUsage.SetBorder(true)
	onlineUsage.SetTitle("Online usage")

	var detail = tview.NewTextView()
	detail.SetBorder(true)
	detail.SetTitle("Detail")

	go func() {
		for {
			select {
			case val := <-quitChannel:
				if val {
					pages.RemovePage("menu")
					var menu = menuPage()
					pages.AddPage("menu", menu, true, true)
					pages.SwitchToPage("menu")
					app.Draw()
					return

				}

			default:
				treeView.SetRoot(nil)
				v, err := mem.VirtualMemory()
				if err != nil {
					panic(err)
				}
				blackBlock := "█"
				graphic_memory := ""
				var numberOfGreenBlocksVirtual = v.UsedPercent / 10
				var numberOfBlackBlocksVirtual = (10 - numberOfGreenBlocksVirtual) + 1
				graphic_memory += "["
				for index := 1; index <= int(numberOfGreenBlocksVirtual); index++ {
					graphic_memory += " " + blackBlock
				}
				for index := 1; index <= int(numberOfBlackBlocksVirtual); index++ {
					graphic_memory += "  "
				}
				graphic_memory += "]"

				s, err := mem.SwapMemory()
				if err != nil {
					panic(err)
				}

				swapp_memory := ""
				var numberOfGreenBlocksSwap = s.UsedPercent / 10
				var numberOfBlackBlocksSwap = (10 - numberOfGreenBlocksSwap) + 1
				swapp_memory += "["
				for index := 1; index <= int(numberOfGreenBlocksSwap); index++ {
					swapp_memory += " " + blackBlock
				}
				for index := 1; index <= int(numberOfBlackBlocksSwap); index++ {
					swapp_memory += "  "
				}
				swapp_memory += "]"

				onlineUsage.SetText("Memory : " + graphic_memory + " " + strconv.Itoa(int(v.UsedPercent)) + "%" +
					"\n" + "Swap   : " + swapp_memory + " " + strconv.Itoa(int(s.UsedPercent)) + "%")
				res, err := ps.Processes()
				if err != nil {
					panic(err)
				}
				for _, procc := range res {
					var node = tview.NewTreeNode(procc.Executable())
					if treeView.GetRoot() == nil {
						treeView.SetRoot(node)
						continue
					}
					treeView.GetRoot().AddChild(node)
				}
				app.Draw()
				time.Sleep(1 * time.Second)
			}
		}
	}()

	flex.AddItem(treeView, 0, 1, true)
	flex.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).AddItem(onlineUsage, 0, 3, false).AddItem(detail, 0, 4, false), 0, 4, false)

	return flex
}

func menuPage() *tview.List {
	var list = tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Y-TOP")
	list.AddItem("Process List", "Show All Procces lists", 'p', func() {
		var procc = ProcessListPage()
		pages.RemovePage("procc")
		pages.AddPage("procc", procc, true, true)
		quitChannel <- false
		pages.SwitchToPage("procc")

	})
	list.AddItem("Quit", "Quit", 'q', func() {
		os.Exit(0)
	})
	return list
}

func main() {
	var menu = menuPage()
	pages.AddPage("menu", menu, true, true)
	app.SetRoot(pages, true)
	err := app.Run()
	if err != nil {
		panic(err)
	}

}
