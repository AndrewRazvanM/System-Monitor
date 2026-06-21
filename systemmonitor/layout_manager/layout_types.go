package layoutmanager

type Direction uint8

const(
	Row Direction = iota
	Column
)

// Geometry stores the pos on screen (x, y) and size (W, H)
type Geometry struct {
	//window positions 
	X, Y int
	// Width and Height of Widget
	W, H int
}

type Split struct {
    Direction // 0 == Horizontal - Row or 1 == Vertical - Column
    Children  []Entry
}

type Entry struct {
    Node   Node
    Weight float32 // 1.0 = share remaining space equally
}

type Node interface {
    UpdatePos(geometry Geometry, isVisible bool)
}