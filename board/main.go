package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)
	g.Mouse = true

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

var step = 0

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	side := 2
	x0, y0 := maxX/2-6*side, maxY/2-2*side
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			v, err := g.SetView(fmt.Sprintf("cell-%d-%d", i, j), x0+i*3*side, y0+j*side, x0+i*3*side+3*side, y0+j*side+side)
			if err != nil && err != gocui.ErrUnknownView {
				return err
			}
			v.Overwrite = true
			if err := g.SetKeybinding(v.Name(), gocui.MouseLeft, gocui.ModNone, click); err != nil {
				log.Panicln(err)
			}
		}
	}
	return nil
}

func click(g *gocui.Gui, v *gocui.View) error {
	if len(v.Buffer()) > 0 {
		return nil
	}

	signal := "o"
	if (step & 1) == 0 {
		signal = "x"
	}
	fmt.Fprintf(v, "  %s", signal)
	step++

	if judge(g, v) {
		win(g, v, signal)
		return nil
	}

	return nil
}

func win(g *gocui.Gui, v *gocui.View, winner string) {
	v.Clear()
	v.Write([]byte(winner + " win"))
	g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, restart)
}

func restart(g *gocui.Gui, v *gocui.View) error {
	for _, v := range g.Views() {
		v.Clear()
	}
	step = 0
	g.DeleteKeybinding("", gocui.KeyEnter, gocui.ModNone)
	return nil
}

func judge(g *gocui.Gui, v *gocui.View) bool {
	return step > 7
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
