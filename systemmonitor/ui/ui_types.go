package ui


type CommandType uint8

const (
	CommandText CommandType = iota
	CommandFill
	CommandRune
)

type Color uint8

const (
	Green Color = iota
	Yellow
	Red
	Standard 
	Blue 
	StandardBold 
)

type CellAttr uint8

const (
	NoAttr CellAttr = iota
	Bold   
	Dim    
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
	Style Color
	//Title contains the widget title
	Title string
}

//used by the Composer interface to issue commands to the renderer
// each pipeline will has it's own Compose method
type DrawCommand struct {
    Type   CommandType // Text, Rect, Bar
    X, Y   int
    W, H   int
    Text   string // used when Type == CommandText
    Char   rune   // used when Type == CommandFill or CommandRune
    Style  Color
	Attr    CellAttr// used to change tcell style
}

//used to build the layout for each widget. Does not handle global layout
type Composer interface {
	Compose(area Geometry) []DrawCommand
}