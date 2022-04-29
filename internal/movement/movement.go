package movement

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"packman/internal/gameboard"
	"packman/internal/score"
	config "packman/internal/utils"
	"time"
)

func ReadInput() (string, error) {
	buffer := make([]byte, 100)
	cnt, err := os.Stdin.Read(buffer)
	if err != nil {
		return "", errors.New("Cant read key")
	}
	movement := convertInput(buffer, cnt)
	return movement, nil
}

func convertInput(buffer []byte, cnt int) string {
	if buffer[0] == 0x1b {
		if cnt == 1 {
			return "ESC"
		}
		if cnt >= 3 && buffer[1] == '[' {
			switch buffer[2] {
			case 'A':
				return "UP"
			case 'B':
				return "DOWN"
			case 'C':
				return "RIGHT"
			case 'D':
				return "LEFT"
			}
		}
	}
	return ""
}
func ProcessMoveAction(game *config.Board, input string, characterType string, gameStatus *config.GameStatus) (*config.Board, *config.GameStatus) {
	if characterType == "packman" {

		pPos := &game.Packman.Position
		pPos.Row, pPos.Col = getNewPos(game, pPos.PrevRow, pPos.PrevCol, input)
		if CheckCollision(game, pPos.Row, pPos.Col, "#") {
			return game, gameStatus
		}
		if CheckCollision(game, pPos.Row, pPos.Col, "-") {
			return game, gameStatus
		}
		if CheckCollision(game, pPos.Row, pPos.Col, "G") {
			MoveToStartingPos(pPos)
			score.MinusLive(gameStatus)
			game.Maze[pPos.PrevRow] = gameboard.ClearSpace(game.Maze[pPos.PrevRow], pPos.PrevCol)
		}
		if CheckCollision(game, pPos.Row, pPos.Col, "X") {
			score.ScoreCandy(gameStatus)
			game.Maze[pPos.PrevRow] = gameboard.ClearSpace(game.Maze[pPos.PrevRow], pPos.PrevCol)
		}
		if CheckCollision(game, pPos.Row, pPos.Col, ".") {
			gameStatus.DotsLeft--
			score.ScoreUp(gameStatus)
			game.Maze[pPos.PrevRow] = gameboard.ClearSpace(game.Maze[pPos.PrevRow], pPos.PrevCol)
		}
		if CheckCollision(game, pPos.Row, pPos.Col, " ") {
			game.Maze[pPos.PrevRow] = gameboard.ClearSpace(game.Maze[pPos.PrevRow], pPos.PrevCol)
		}
		pPos.PrevRow, pPos.PrevCol = pPos.Row, pPos.Col
		game.Maze[pPos.Row] = changeCharacter(game.Maze[pPos.Row], pPos.Col, "P")
	}
	if characterType == "ghost" {
		gPos := &game.Ghosts[1].Position
		gPos.Row, gPos.Col = getNewPos(game, gPos.PrevRow, gPos.PrevCol, input)
		if CheckCollision(game, gPos.Row, gPos.Col, "#") {
			return game, gameStatus
		}
		if CheckCollision(game, gPos.Row, gPos.Col, "P") {
			MoveToStartingPos(&game.Packman.Position)
			score.MinusLive(gameStatus)
		}
		game.Maze[gPos.PrevRow] = changeCharacter(game.Maze[gPos.PrevRow], gPos.PrevCol, "G")
		gPos.PrevRow, gPos.PrevCol = gPos.Row, gPos.Col
		game.Maze[gPos.Row] = changeCharacter(game.Maze[gPos.Row], gPos.Col, "G")

	}

	return game, gameStatus
}

func CheckCollision(game *config.Board, currentRow, currentCol int, collisionCharacter string) bool {
	rowStr := game.Maze[currentRow]
	character := string(rowStr[currentCol])
	if character == collisionCharacter {
		return true
	}
	return false
}

func RandomDirection() string {
	dir := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "RIGHT",
		3: "LEFT",
	}
	return move[dir]
}

func GetCharacterDirection() (string, error) {
	input, err := ReadInput()
	if err != nil {
		log.Println("error reading input:", err)
		return "", err
	}

	return input, nil
}

func MoveGhosts(game *config.Board, gameStatus *config.GameStatus) {
	for {
		select {
		case <-time.After(200 * time.Millisecond):
			ProcessMoveAction(game, RandomDirection(), "ghost", gameStatus)
		}
	}

}

func MoveToStartingPos(pPos *config.CharacterPosition) {
	pPos.Row, pPos.Col = config.StartingRow, config.StartingCol
}

func getNewPos(game *config.Board, oldRow, oldCol int, dir string) (newRow, newCol int) {
	newRow, newCol = oldRow, oldCol

	switch dir {
	case "UP":
		newRow = newRow - 1
		if newRow < 0 {
			newRow = len(game.Maze) - 1
		}
	case "DOWN":
		newRow = newRow + 1
		if newRow == len(game.Maze)-1 {
			newRow = 0
		}
	case "RIGHT":
		newCol = newCol + 1
		if newCol == len(game.Maze[0]) {
			newCol = 0
		}
	case "LEFT":
		newCol = newCol - 1
		if newCol < 0 {
			newCol = len(game.Maze[0]) - 1
		}
	}

	if game.Maze[newRow][newCol] == '#' {
		newRow = oldRow
		newCol = oldCol
	}

	return
}

func changeCharacter(rowStr string, col int, characterType string) string {
	character := []byte(characterType)
	row := []byte(rowStr)
	row[col] = character[0]
	return string(row)
}
