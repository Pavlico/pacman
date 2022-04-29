package score

import (
	config "packman/internal/utils"
)

func MinusLive(gameStatus *config.GameStatus) {
	gameStatus.Lives--
}

func ScoreUp(gameStatus *config.GameStatus) {
	gameStatus.Score += 10
}

func ScoreCandy(gameStatus *config.GameStatus) {
	gameStatus.Score += 20
}
