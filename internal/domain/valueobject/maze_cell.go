package valueobject

type MazeCell struct {
	Right  bool
	Bottom bool
}

type MazeMap [6][6]MazeCell
