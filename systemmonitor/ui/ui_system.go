package ui

const (
    HLine = '─'
    VLine = '│'
    TLCorner = '┌'
    TRCorner = '┐'
    BLCorner = '└'
    BRCorner = '┘'
)

func (w *Widget) Initalize(x, y, W, H int, title string, style uint8) {
	w.X, w.Y, w.W, w.H = x, y, W, H
	w.IsVisible = true
	w.Style = style
	w.Title = title
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
				Data:  HLine,
				Style: style,
			})
	renderCmds = append(renderCmds, DrawCommand{
				Type:  CommandFill,
				X:     x0,
				Y:     y1,
				W:     x1 - x0,
				H:     1,
				Data:  HLine,
				Style: style,
			})

	// fills the sides
	renderCmds = append(renderCmds, DrawCommand{
		Type:  CommandFill,
		X:     x0,
		Y:     y0,
		W:     1,
		H:     y1 - y0,
		Data:  VLine,
		Style: style,
	})

	renderCmds = append(renderCmds, DrawCommand{
		Type:  CommandFill,
		X:     x1,
		Y:     y0,
		W:     1,
		H:     y1 - y0,
		Data:  VLine,
		Style: style,
	})

	// corners
	renderCmds = append(renderCmds, DrawCommand{
		Type:  CommandFill,
		X:     x0,
		Y:     y0,
		W:     1,
		H:     1,
		Data:  TLCorner,
		Style: style,
	})
	renderCmds = append(renderCmds, DrawCommand{
		Type:  CommandFill,
		X:     x1,
		Y:     y0,
		W:     1,
		H:     1,
		Data:  TRCorner,
		Style: style,
	})
	renderCmds = append(renderCmds, DrawCommand{
		Type:  CommandFill,
		X:     x0,
		Y:     y1,
		W:     1,
		H:     1,
		Data:  BLCorner,
		Style: style,
	})
	renderCmds = append(renderCmds, DrawCommand{
		Type:  CommandFill,
		X:     x1,
		Y:     y1,
		W:     1,
		H:     1,
		Data:  BRCorner,
		Style: style,
	})

	return renderCmds
}