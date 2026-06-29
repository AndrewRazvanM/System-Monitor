package inputhandler

import (
	layoutmanager "github.com/AndrewRazvanM/System-Monitor/systemmonitor/layout_manager"
	renderingsystem "github.com/AndrewRazvanM/System-Monitor/systemmonitor/rendering_system"
	ui "github.com/AndrewRazvanM/System-Monitor/systemmonitor/ui"
	"github.com/gdamore/tcell/v2"
)

func (h *InputHandler) HandleInput (x,y int, root layoutmanager.Split, buffer renderingsystem.ScreenBuffer){
    screen := buffer.Screen
    for {
        ev := screen.PollEvent()

        switch e := ev.(type) {
        case *tcell.EventMouse:
            x, y := e.Position()

            if e.Buttons () & tcell.Button1 != 0 {
                h.HandleClick(x, y)
            }
        
        // need to redo this part later. Doesn't redraw everything on the screen and causes artifacts
        case *tcell.EventResize:
            buffer.Resize()
            root.UpdatePos(layoutmanager.Geometry{
                X: 0,
                Y: 0,
                W: buffer.Width,
                H: buffer.Height,
            }, true)

        case *tcell.EventKey:
            if e.Key() == tcell.KeyEscape ||
               e.Key() == tcell.KeyCtrlC {
                screen.Fini()
                return
            }
        }
    }
}

func (h *InputHandler) HandleClick(x, y int) {
    widgets := h.WidgetButtons
    for _, w := range widgets {
        if !w.IsVisible {
            continue
        }

        // check widget bounds first
        if !hit(w.Geometry, x, y) {
            continue
        }

        // convert to widget-local space
        lx := x - w.X
        ly := y - w.Y

        for _, b := range w.Buttons {
            if !b.Visible {
                continue
            }

            if hitRelative(b, lx, ly) {
                b.Behavior()
                return
            }
        }
    }

    // global buttons (no widget offset)
    for _, b := range h.GlobalButtons {
        if hit(b.Geometry, x, y) {
            b.Behavior()
            return
        }
    }
}
// hit checks if the click is within a Geometry
func hit (pos layoutmanager.Geometry, mx, my int) bool {
    if mx >= pos.X && mx <= pos.X + pos.W && my >= pos.Y && my <= pos.H + pos.Y {
        return true
    }
    return false
}
// hitRelative checks if the clicks is within a relative Geometry. 
// The button is embedded into another UI element and it uses that element's co-ordinates as the absolute ones.
func hitRelative (b ui.Button, lx, ly int) bool {
    if lx >= b.X && lx <= b.X + b.W && ly >= b.Y && ly <= b.H + b.Y {
        return true
    }
    return false
}