package gameboard

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	config "packman/internal/utils"

	"github.com/danicat/simpleansi"
)

type BoardStruct struct {
	Board      config.Board
	Characters CharactersPositions
}
type CharactersPositions struct {
	Packman config.Packman
	Ghosts  []config.Ghost
}

func Initialize() *config.Board {
	board := &config.Board{}
	cbTerm := exec.Command("stty", "cbreak", "-echo")
	cbTerm.Stdin = os.Stdin

	err := cbTerm.Run()
	if err != nil {
		log.Fatalln("unable to activate cbreak mode:", err)
	}
	LoadMaze(board, "maze01.txt")
	LoadCharacters(board)
	if err != nil {
		log.Println("failed to load maze:", err)
		return nil
	}
	return board
}

func Cleanup() {
	cookedTerm := exec.Command("stty", "-cbreak", "echo")
	cookedTerm.Stdin = os.Stdin

	err := cookedTerm.Run()
	if err != nil {
		log.Fatalln("unable to activate cooked mode:", err)
	}
}

func PrintScreen(bg *config.Board) {
	simpleansi.ClearScreen()
	for _, line := range bg.Maze {
		fmt.Println(line)
	}
}
func LoadMaze(bg *config.Board, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		bg.Maze = append(bg.Maze, line)
	}

	return nil
}

func LoadCharacters(bg *config.Board) {
	ghostNum := 0
	for row, line := range bg.Maze {
		for col, char := range line {
			switch char {
			case 'P':
				bg.Packman = config.Packman{
					Position: config.CharacterPosition{
						Row:     row,
						Col:     col,
						PrevRow: row,
						PrevCol: col,
					},
				}
			case 'G':
				ghostNum++
				bg.Ghosts = append(bg.Ghosts, config.Ghost{
					Position: config.CharacterPosition{
						Row: row, Col: col, PrevRow: row, PrevCol: col,
					},
					Status: config.GhostStatusNormal,
					StepOn: " ",
				})
			case '.':
				bg.DotsNum++
			}
		}
	}

}

func ClearSpace(rowStr string, col int) string {
	spaceBarByte := uint8(32)
	row := []byte(rowStr)
	row[col] = spaceBarByte
	return string(row)
}
