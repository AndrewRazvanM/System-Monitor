package renderingsystem

import (
	"fmt"

	ui "github.com/AndrewRazvanM/System-Monitor/systemmonitor/ui"
	"github.com/gdamore/tcell/v2"
)

// render copies all modified cells from the current snapshot to the underlying
// tcell screen.
//
// Each cell in the current snapshot is compared against the corresponding cell
// in the previous snapshot. Only cells that have changed are written to the
// screen, minimizing terminal updates.
//
// After a cell is rendered, the previous snapshot is updated to match the
// current snapshot. This method calls Screen.Show();
func (sb *ScreenBuffer) Render() {
	styles := StyleList
	attributes := StyleAttr
	h := sb.Height
	w := sb.Width
	screen := sb.Screen
	for y := range h {
		for x := range w {
			index := y * w + x
			//render only if content is different
			if sb.Previous.Cells[index] == sb.Current.Cells[index] {
				continue
			}
			cell := sb.Current.Cells[index]
			screen.SetContent(x, y, cell.Rune, nil, styles[cell.Style].Attributes(attributes[cell.Attribute]))

			sb.Previous.Cells[index] = cell
		}	
	}
	sb.Screen.Show()	
}

//Processes commands issued by the composers. Needs to be called before Render()
func (sb *ScreenBuffer) ProcessCmds(cmds []ui.DrawCommand) {
    for _, c := range cmds {
        switch c.Type {

        case ui.CommandText:
            sb.writeRunes(
                c.X,
                c.Y,
                c.Text,
                c.Style,
				c.Attr,
            )

        case ui.CommandFill:
            sb.fillRect(
                c.X,
                c.Y,
                c.W,
                c.H,
                c.Char,
                c.Style,
				c.Attr,
            )

        case ui.CommandRune:
            sb.setCell(
                c.X,
                c.Y,
                c.Char,
                c.Style,
				c.Attr,
            )
        }
    }
}
// UpdateCell writes a single cell at absolute (x, y).
func (sb *ScreenBuffer) setCell(x, y int, r rune, style ui.Color, attr ui.CellAttr) {
	if x < 0 || y < 0 || x >= sb.Width || y >= sb.Height {
		return
	}
	index := y*sb.Width + x
	sb.Current.Cells[index].Rune = r
	sb.Current.Cells[index].Style = style
	sb.Current.Cells[index].Attribute = attr
}

func (sb *ScreenBuffer) writeRunes(x, y int, runes []rune, style ui.Color, attr ui.CellAttr) (written int) {
	if x < 0 || y < 0 || y >= sb.Height || x >= sb.Width {
		return 0
	}

	max := min(len(runes), sb.Width - x)

	index := y*sb.Width + x
	cells := sb.Current.Cells

	for i := range max {
		cells[index+i].Rune = runes[i]
		cells[index+i].Style = style
		cells[index+i].Attribute = attr
	}

	return max
}

func (sb *ScreenBuffer) fillRect(x, y, w, h int, r rune, style ui.Color, attr ui.CellAttr) (written int) {
	if w <= 0 || h <= 0 {
		return 0
	}

	if x < 0 || y < 0 {
		return 0
	}

	if x >= sb.Width || y >= sb.Height {
		return 0
	}

	if x+w > sb.Width {
		w = sb.Width - x
	}

	if y+h > sb.Height {
		h = sb.Height - y
	}

	cells := sb.Current.Cells
	screenW := sb.Width

	for row := 0; row < h; row++ {
		base := (y+row)*screenW + x

		for col := 0; col < w; col++ {
			cells[base+col].Rune = r
			cells[base+col].Style = style
			cells[base+col].Attribute = attr
			written++
		}
	}

	return written
}

func (sb *ScreenBuffer) Clear(r rune, style ui.Color) {
	cells := sb.Current.Cells

	for i := range cells {
		cells[i].Rune = r
		cells[i].Style = style
	}
}

func (sb *ScreenBuffer) Resize(){
	sb.Screen.Sync()
	width, height := sb.Screen.Size()
	sb.Current.Cells = make([]Cell, width * height)
	sb.Previous.Cells = make([]Cell, width * height)
	sb.Height = height
	sb.Width = width
}

func (sb *ScreenBuffer) Init(isVisible bool, style ui.Color) error {
	sb.IsVisible = isVisible
	screen, sErr := tcell.NewScreen()
	if sErr != nil {
		return fmt.Errorf("Error creating window: %v", sErr)
	}

	initError := screen.Init()
	if initError != nil {
		return fmt.Errorf("Error initializing window: %v", sErr)
	}
	sb.Screen = screen
	width, height := screen.Size()
	//create current and previous buffers

	sb.Current.Cells = make([]Cell, width * height)
	sb.Previous.Cells = make([]Cell, width * height)
	
	sb.Height = height
	sb.Width = width
	//enable mouse
	screen.EnableMouse()

	return nil
}