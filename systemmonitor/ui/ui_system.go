package ui

import (
	layoutmanager "github.com/AndrewRazvanM/System-Monitor/systemmonitor/layout_manager"
)

const (
    HLine = '─'
    VLine = '│'
    TLCorner = '┌'
    TRCorner = '┐'
    BLCorner = '└'
    BRCorner rune = '┘'
)

func (w *Widget) SetAttributes(title string, style Color) {
	w.Style = style
	titleStr := "┐" + title + "┌"
	w.Title = []rune(titleStr)
}

func (w Widget) Draw() []DrawCommand {
	if !w.IsVisible {
		return nil
	}
	renderCmds := make([]DrawCommand, 0, 12)
	x0 := w.X
	y0 := w.Y
	x1 := w.X + w.W - 1
	y1 := w.Y + w.H - 1

	style := w.Style

	// fills top and bottom
	renderCmds = append(renderCmds, DrawCommand{
				Type:  CommandFill,
				X:     x0,
				Y:     y0,
				W:     x1 - x0,
				H:     1,
				Char:  HLine,
				Style: style,
			})
	renderCmds = append(renderCmds, DrawCommand{
				Type:  CommandFill,
				X:     x0,
				Y:     y1,
				W:     x1 - x0,
				H:     1,
				Char:  HLine,
				Style: style,
			})

	// fills the sides
	renderCmds = append(renderCmds, DrawCommand{
		Type:  CommandFill,
		X:     x0,
		Y:     y0,
		W:     1,
		H:     y1 - y0,
		Char:  VLine,
		Style: style,
	})

	renderCmds = append(renderCmds, DrawCommand{
		Type:  CommandFill,
		X:     x1,
		Y:     y0,
		W:     1,
		H:     y1 - y0,
		Char:  VLine,
		Style: style,
	})

	// corners
	renderCmds = append(renderCmds, DrawCommand{
		Type:  CommandFill,
		X:     x0,
		Y:     y0,
		W:     1,
		H:     1,
		Char:  TLCorner,
		Style: style,
	})
	renderCmds = append(renderCmds, DrawCommand{
		Type:  CommandFill,
		X:     x1,
		Y:     y0,
		W:     1,
		H:     1,
		Char:  TRCorner,
		Style: style,
	})
	renderCmds = append(renderCmds, DrawCommand{
		Type:  CommandFill,
		X:     x0,
		Y:     y1,
		W:     1,
		H:     1,
		Char:  BLCorner,
		Style: style,
	})
	renderCmds = append(renderCmds, DrawCommand{
		Type:  CommandFill,
		X:     x1,
		Y:     y1,
		W:     1,
		H:     1,
		Char:  BRCorner,
		Style: style,
	})

	//draw the title
	renderCmds = append(renderCmds, DrawCommand{
		Type: CommandText,
		X: x0 + 2,
		Y: y0,
		Text: w.Title,
		Style: style,
		Attr: Bold,
	})

	for _, b := range w.Buttons {
		renderCmds = b.Draw(renderCmds)
	}
	return renderCmds
}

func (w *Widget) UpdatePos(g layoutmanager.Geometry, isVisible bool) {
	w.Geometry = g
	w.IsVisible = isVisible
}

func (b Button) Draw(renderCmds []DrawCommand) []DrawCommand {
	if !b.Visible {
		return nil
	}
	x0 := b.X
	y0 := b.Y
	text := []rune("| " + b.Label + " |")

	renderCmds = append(renderCmds, DrawCommand{
        Type: CommandText,
        X: x0,
        Y: y0,
		H: 0,
        Text: text,
		W: len(text),
        Style: Blue,
        Attr: Bold,
    })

	return renderCmds
}
//Button possible behaviours

//Toggles the Widget on and off.
func ToggleWidget (widget *Widget) {
	if widget.IsVisible == true {
		widget.IsVisible = false
		return
	}

	widget.IsVisible = true
}