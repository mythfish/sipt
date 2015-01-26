package sipt

import (
	"github.com/andlabs/ui"
)

type ModelDialog struct {
	win  ui.Window
	btn1 ui.Button
	btn2 ui.Button
}

func NewModelDialog(title, notice string) *ModelDialog {
	btn1 := ui.NewButton("Yes")
	//btn2 := ui.NewButton("No")
	n := ui.NewLabel(notice)

	stack1 := ui.NewHorizontalStack(
		btn1,
		//btn2,
	)

	stack2 := ui.NewVerticalStack(
		n,
		stack1,
	)
	win := ui.NewWindow(title, 200, 100, stack2)

	btn1.OnClicked(func() {
		win.Close()
		return
	})
	return &ModelDialog{
		win:  win,
		btn1: btn1,
		//btn2: btn2,
	}
}

func (md *ModelDialog) Show() {
	md.win.Show()
}
