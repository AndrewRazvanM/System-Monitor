package ui

const (
    HLine = '─'
    VLine = '│'
    TLCorner = '┌'
    TRCorner = '┐'
    BLCorner = '└'
    BRCorner rune = '┘'
)

func (w *Widget) Initalize(x, y, W, H int, title string, style Color) {
	w.X, w.Y, w.W, w.H = x, y, W, H
	w.IsVisible = true
	w.Style = style
	w.Title = "┐" + title + "┌"
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

	return renderCmds
}