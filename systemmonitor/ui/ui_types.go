package ui


type CommandType uint8

const (
	CommandText CommandType = iota
	CommandFill
	CommandRune
)

const (
	Green uint8 = 0
	Yellow uint8 = 1
	Red uint8 = 2
	Standard uint8 = 3 // maps to tcell.StyleDefault
	Blue uint8 = 4
	Vline rune =  '│'
)

// Geometry stores the pos on screen (x, y) and size (W, H)
type Geometry struct {
	//window positions 
	X, Y int
	// Width and Height of Widget
	W, H int
}

// Widget stores the position and dimensions for 
// each dashboard used by the renderer
type Widget struct {
	//the Widget knows where it is, cause it knows where it isn't
	//only the layout manager updates this
	Geometry
	//IsVisible says if the widget should show or not
	IsVisible bool
	//Dirty says if the widget needs to be redrawn
	Dirty bool
	//Style stores the style used by the widget
	Style uint8
	//Title contains the widget title
	Title string
}

//used by the Composer interface to issue commands to the renderer
// each pipeline will has it's own Compose method
type DrawCommand struct {
    Type   CommandType // Text, Rect, Bar
    X, Y   int
    W, H   int
    Data   any
    Style  uint8
}
//used to build the layout for each widget. Does not handle global layout
type Composer interface {
	Compose()
}