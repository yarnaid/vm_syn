package main

import (
	"fmt"
	"time"

	"github.com/yarnaid/vm_syn/vminstance"

	ui "github.com/gizak/termui"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

func startGui(vm vminstance.VM) {
	app := tview.NewApplication()
	timerText := tview.NewTextView()
	timerText.SetTitle("timer").SetBorder(true)

	loggerText := tview.NewTextView()
	loggerText.SetTitle("logger").SetBorder(true)

	start := time.Now()

	infoList := tview.NewList()
	infoList.SetTitle("info").SetBorder(true)
	infoList.AddItem("status", "???", 'a', nil)
	infoList.AddItem("addr", "???", 'a', nil)

	regList := tview.NewList()
	regList.SetBorder(true).SetTitle("regs")
	regList.ShowSecondaryText(false)
	for i := 0; i < 8; i++ {
		regList.AddItem("", "", '0'+rune(i), nil)
	}

	stackList := tview.NewList()
	stackList.SetBorder(true).SetTitle("stack")
	stackList.ShowSecondaryText(false)
	for i := 0; i < 10; i++ {
		stackList.AddItem("", "", '0'+rune(i), nil)
	}

	output := tview.NewTextView()
	output.SetTitle("output").SetBorder(true)
	vm.SetTerminal(output)

	flex := tview.NewFlex().
		AddItem(regList, 0, 1, false).
		AddItem(stackList, 0, 1, false).
		AddItem(output, 0, 2, false).
		AddItem(infoList, 0, 3, false).
		AddItem(timerText, 0, 1, false).
		AddItem(loggerText, 0, 5, false)
	ticker := time.NewTicker(time.Second * time.Duration(1))
	go func() {
		for {
			select {
			case <-ticker.C:
				app.QueueUpdateDraw(
					func() {
						timerText.Write([]byte(time.Since(start).String() + "\n"))
						reg := vm.CPU().Registry()
						for i := 0; i < 8; i++ {
							regList.SetItemText(i, reg.ToStringList()[i], "")
						}
						stackHead := vm.CPU().Stack().Head()
						for i := 0; i < len(stackHead); i++ {
							stackList.SetItemText(i, string(stackHead[i]), "")
						}
						m, _ := infoList.GetItemText(0)
						infoList.SetItemText(0, m, vm.Status().String())
						m, _ = infoList.GetItemText(1)
						infoList.SetItemText(1, m, fmt.Sprint(vm.Addr()))
					})
			}
		}
	}()
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func startGui2(vm vminstance.VM) {
	if err := ui.Init(); err != nil {
		logrus.Error(errors.Wrap(err, "gui init error"))
	}
	defer ui.Close()

	output := ui.NewParagraph("")
	output.BorderLabel = "output"
	// output.
	// vm.SetTerminal(output.Buffer().)

	regsData := ui.NewList()
	regsData.ItemFgColor = ui.ColorYellow
	regsData.BorderLabel = "regs"
	regsData.Height = 10

	infoList := ui.NewList()
	infoList.BorderLabel = "info"
	infoList.ItemFgColor = ui.ColorYellow
	infoList.Height = 4

	stackList := ui.NewList()
	stackList.BorderLabel = "stack(0)"
	stackList.ItemFgColor = ui.ColorBlue
	stackList.Height = 12

	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(4, 0, regsData), ui.NewCol(4, 0, stackList)),
		ui.NewRow(ui.NewCol(4, 0, infoList)),
	)

	ui.Body.Align()
	ui.Render(ui.Body)

	tickerCount := 1
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Millisecond * time.Duration(100)).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				ui.Body.Width = payload.Width
				ui.Body.Align()
				ui.Clear()
				ui.Render(ui.Body)
			}
		case <-ticker:
			regsData.Items = append(vm.Registry().ToStringList())
			regsData.Height = len(regsData.Items) + 2

			infoList.Items = []string{
				formatInfo("status", vm.Status()),
				formatInfo("stack depth", vm.CPU().Stack().Len()),
			}

			stackList.Items = formatStack(vm.CPU().Stack().Head())
			stackList.BorderLabel = fmt.Sprintf("stack(%v)", vm.CPU().Stack().Len())
			ui.Render(ui.Body)
			tickerCount++
		}
	}
}

func formatInfo(key, value interface{}) string {
	return fmt.Sprintf("[%v]: %v", key, value)
}

func formatStack(stack []vminstance.StackData) []string {
	stackHeadLen := len(stack)
	res := make([]string, stackHeadLen)
	for i := 0; i < stackHeadLen; i++ {
		res[i] = fmt.Sprint(stack[i])
	}
	return res
}
