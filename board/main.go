package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"time"
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
	if err := g.SetKeybinding("", gocui.MouseLeft, gocui.ModNone, click); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

var side = 2

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	x0, y0 := maxX/2-6*side, maxY/2-2*side
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			v, err := g.SetView(fmt.Sprintf("cell-%d-%d", i, j), x0+i*3*side, y0+j*side, x0+i*3*side+3*side, y0+j*side+side)
			if err != nil && err != gocui.ErrUnknownView {
				return err
			}
			v.Overwrite = true
		}
	}
	return nil
}

var step = 0

func click(g *gocui.Gui, v *gocui.View) error {
	if finish {
		return nil
	}
	if len(v.Buffer()) > 0 {
		return nil
	}

	signal := "o"
	if (step & 1) == 0 {
		signal = "x"
	}
	fmt.Fprintf(v, "  %s", signal)
	step++

	g.Update(func(gui *gocui.Gui) error {
		if judge(g) {
			return win(g, signal)
		}
		return nil
	})

	return nil
}

var finish = false

func win(g *gocui.Gui, winner string) error {
	finish = true
	time.Sleep(1 * time.Second)
	x0, y0, x1, y1, err := g.ViewPosition("cell-1-1")
	if err != nil {
		return err
	}

	dx, dy := (x1-x0)/2, (y1-y0)/2
	v, err := g.SetView("win", x0-dx, y0-dy, x1+dx, y1+dy)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	fmt.Fprintf(v, "%s win!", winner)
	g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, restart)
	g.SetCurrentView(v.Name())
	return nil
}

func restart(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView(v.Name())
	for _, v := range g.Views() {
		v.Clear()
	}
	step = 0
	g.DeleteKeybinding("", gocui.KeyEnter, gocui.ModNone)
	return nil
}

func judge(g *gocui.Gui) bool {
	return step > 1
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
