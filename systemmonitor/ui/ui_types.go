package ui

import (
	layoutmanager "github.com/AndrewRazvanM/System-Monitor/systemmonitor/layout_manager"
)

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
)

type CellAttr uint8

const (
	NoAttr CellAttr = iota
	Bold   
	Dim    
)

type Button struct {
    Label string

	//stores the position of a button
	//if it is within a widget, it's position is relative to the widget.
    layoutmanager.Geometry
	
	Visible bool
	Behavior func()
}

// Widget stores the position and dimensions for 
// each dashboard used by the renderer
type Widget struct {
	//the Widget knows where it is, cause it knows where it isn't
	//only the layout manager updates this
	layoutmanager.Geometry
	//IsVisible says if the widget should show or not
	IsVisible bool
	//Style stores the style used by the widget
	Style Color
	//Title contains the widget title
	Title []rune
	//Tells the layout manager, on which node this widget belongs
	Node layoutmanager.Node
	//Lists buttons associated with the widget
	Buttons []Button
}

//used by the Composer interface to issue commands to the renderer
// each pipeline will has it's own Compose method
type DrawCommand struct {
    Type   CommandType // Text, Rect, Bar
    X, Y   int
    W, H   int
    Text   []rune // used when Type == CommandText
    Char   rune   // used when Type == CommandFill or CommandRune
    Style  Color
	Attr   CellAttr// used to change tcell style
}

//used to build the layout for each widget. Does not handle global layout
type Composer interface {
	Compose(area layoutmanager.Geometry, isVisible bool) []DrawCommand
}