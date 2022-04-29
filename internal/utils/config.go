package config

type CharacterPosition struct {
	Row     int
	Col     int
	PrevRow int
	PrevCol int
}
type Packman struct {
	Position CharacterPosition
}

type Ghost struct {
	Position CharacterPosition
	Status   string
}

type Board struct {
	Maze    []string
	Ghosts  []Ghost
	Packman Packman
	DotsNum int
}

type GameStatus struct {
	On       bool
	Score    int
	Lives    int
	DotsLeft int
}

const GhostStatusNormal = "normal"
const StartingRow = 14
const StartingCol = 13
