package valueobject

type MazeCell struct {
	Right  bool
	Bottom bool
}

type Maze struct {
	Marker1 Point2D
	Marker2 Point2D
	Map     [6][6]MazeCell
}
