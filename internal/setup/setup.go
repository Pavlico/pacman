package setup

import (
	"fmt"
	"os"
	"os/signal"
	"packman/internal/gameboard"
	"packman/internal/movement"
	config "packman/internal/utils"
	"syscall"
	"time"
)

func StartGame() {
	defer gameboard.Cleanup()
	game := gameboard.Initialize()
	gameStatus := &config.GameStatus{
		On:       true,
		Lives:    3,
		Score:    0,
		DotsLeft: game.DotsNum,
	}
	input := make(chan string)
	watchForSignal(gameStatus)
	go watchForUserAction(input, game, gameStatus)
	for i := 0; i < len(game.Ghosts); i++ {
		go movement.MoveGhosts(game, gameStatus, i)
	}
	for {
		select {
		case inp := <-input:
			movement.ProcessMoveAction(game, inp, gameStatus)
			if inp == "ESC" {
				gameStatus.On = false
			}
		default:
		}
		select {
		case <-time.After(100 * time.Millisecond):
			gameboard.PrintScreen(game)
			fmt.Println("Score: ", gameStatus.Score, " Lives: ", gameStatus.Lives)
		}
		if gameStatus.On == false || gameStatus.Lives == 0 || gameStatus.DotsLeft == 0 {
			break
		}

	}
}

func watchForSignal(gameStatus *config.GameStatus) {
	abortSingalChannel := make(chan os.Signal, 1)
	signal.Notify(abortSingalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGPROF)
	go func() {
		<-abortSingalChannel
		gameStatus.On = false
	}()
}

func watchForUserAction(input chan string, game *config.Board, gameStatus *config.GameStatus) {
	for {
		direction, _ := movement.GetCharacterDirection()
		input <- direction
	}

}
