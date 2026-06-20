package renderingsystem

type ScreenBuffer struct {
	Current Snapshot
	Previous Snapshot
	Width uint16
	Height uint16
}

// Cell stores the actual character (element) that will be rendered.
// Also contains the position it should be rendered at and it's style
type Cell struct {
    Rune rune
    Style uint8
}

// Snapshot stores the elements that will be rendered.
// Each Snapshot is a 2D Grid mapped on a slice (1D).
type Snapshot struct {
	// Slice of 2D positions (x,y). Basically mapping 2D to 1D:
	// Cells[Index] := Y * W + X 
    Cells []Cell
}