package config

import "time"

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
	StepOn   string
	Status   string
}

type Board struct {
	Maze      []string
	Ghosts    []Ghost
	Packman   Packman
	DotsNum   int
	PillTimer *time.Timer
}

type GameStatus struct {
	On       bool
	Score    int
	Lives    int
	DotsLeft int
}

const GhostStatusNormal = "normal"
const GhostStatusBlue = "blue"
const StartingRow = 14
const StartingCol = 13
const StartingGRow = 11
const StartingGCol = 13
const PillDuration = 10
