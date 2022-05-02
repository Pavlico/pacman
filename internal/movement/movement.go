package movement

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"packman/internal/gameboard"
	"packman/internal/score"
	config "packman/internal/utils"
	"sync"
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
func ProcessMoveAction(game *config.Board, input string, gameStatus *config.GameStatus) (*config.Board, *config.GameStatus) {
	pPos := &game.Packman.Position
	pPos.Row, pPos.Col = getNewPos(game, pPos.PrevRow, pPos.PrevCol, input)
	if CheckCollision(game, pPos.Row, pPos.Col, "#") || CheckCollision(game, pPos.Row, pPos.Col, "-") {
		return game, gameStatus
	}
	if CheckCollision(game, pPos.Row, pPos.Col, "G") {
		if game.Ghosts[0].Status != config.GhostStatusNormal {
			index := whichGhost(game.Ghosts, pPos.Row, pPos.Col)
			moveGhostToStartingPos(&game.Ghosts[index])

		} else {
			MoveToStartingPos(pPos)
			score.MinusLive(gameStatus)
		}

	}
	if CheckCollision(game, pPos.Row, pPos.Col, "X") {
		score.ScoreCandy(gameStatus)
		go processPill(game)
	}
	if CheckCollision(game, pPos.Row, pPos.Col, ".") {
		gameStatus.DotsLeft--
		score.ScoreUp(gameStatus)
	}
	if CheckCollision(game, pPos.Row, pPos.Col, " ") {
	}
	game.Maze[pPos.PrevRow] = gameboard.ClearSpace(game.Maze[pPos.PrevRow], pPos.PrevCol)
	pPos.PrevRow, pPos.PrevCol = pPos.Row, pPos.Col
	game.Maze[pPos.Row] = changeCharacter(game.Maze[pPos.Row], pPos.Col, "P")

	return game, gameStatus
}
func ProcessGhostAction(game *config.Board, input string, ghostNum int, gameStatus *config.GameStatus) (*config.Board, *config.GameStatus) {
	var m sync.Mutex
	m.Lock()
	defer m.Unlock()
	stepOn := game.Ghosts[ghostNum].StepOn
	gPos := &game.Ghosts[ghostNum].Position
	gPos.Row, gPos.Col = getNewPos(game, gPos.PrevRow, gPos.PrevCol, input)
	if CheckCollision(game, gPos.Row, gPos.Col, "#") {
		return game, gameStatus
	}
	if CheckCollision(game, gPos.Row, gPos.Col, "G") {
		return game, gameStatus
	}

	if stepOn == "." || stepOn == "X" || stepOn == " " || stepOn == "-" {
		game.Maze[gPos.PrevRow] = changeCharacter(game.Maze[gPos.PrevRow], gPos.PrevCol, stepOn)
		game.Ghosts[ghostNum].StepOn = ""
	}
	if CheckCollision(game, gPos.Row, gPos.Col, "P") {
		if game.Ghosts[ghostNum].Status != config.GhostStatusNormal {
			moveGhostToStartingPos(&game.Ghosts[ghostNum])
			return game, gameStatus
		} else {
			MoveToStartingPos(&game.Packman.Position)
			game.Ghosts[ghostNum].StepOn = " "
			score.MinusLive(gameStatus)
		}

	}
	if CheckCollision(game, gPos.Row, gPos.Col, ".") {
		game.Ghosts[ghostNum].StepOn = "."
	}
	if CheckCollision(game, gPos.Row, gPos.Col, "X") {
		game.Ghosts[ghostNum].StepOn = "X"
	}
	if CheckCollision(game, gPos.Row, gPos.Col, " ") {
		game.Ghosts[ghostNum].StepOn = " "
	}
	if CheckCollision(game, gPos.Row, gPos.Col, "-") {
		game.Ghosts[ghostNum].StepOn = "-"
	}

	gPos.PrevRow, gPos.PrevCol = gPos.Row, gPos.Col
	game.Maze[gPos.Row] = changeCharacter(game.Maze[gPos.Row], gPos.Col, "G")

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

func MoveGhosts(game *config.Board, gameStatus *config.GameStatus, ghostNum int) {
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			ProcessGhostAction(game, RandomDirection(), ghostNum, gameStatus)
		}
	}

}

func MoveToStartingPos(pPos *config.CharacterPosition) {
	pPos.PrevRow, pPos.PrevCol = config.StartingRow, config.StartingCol
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

func processPill(game *config.Board) {
	var m sync.Mutex
	m.Lock()
	game.Ghosts = updateGhosts(game.Ghosts, config.GhostStatusBlue)
	if game.PillTimer != nil {
		game.PillTimer.Stop()
	}
	game.PillTimer = time.NewTimer(time.Second * config.PillDuration)
	m.Unlock()
	<-game.PillTimer.C
	m.Lock()
	game.PillTimer.Stop()
	game.Ghosts = updateGhosts(game.Ghosts, config.GhostStatusNormal)
	m.Unlock()
}

func moveGhostToStartingPos(g *config.Ghost) {
	g.Position.PrevRow, g.Position.PrevCol = config.StartingGRow, config.StartingGCol
}

func updateGhosts(g []config.Ghost, status string) []config.Ghost {
	var m sync.Mutex
	m.Lock()
	defer m.Unlock()
	for i, _ := range g {
		g[i].Status = status

	}
	return g
}

func whichGhost(ghosts []config.Ghost, row, col int) int {
	var index int
	for i, ghost := range ghosts {
		if ghost.Position.Col == col && ghost.Position.Row == row {
			return i
		}
	}
	return index
}
