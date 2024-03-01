package gameUI

import (
	gameSrvs "battlship/internal/service/game"
	"battlship/internal/service/game/domain"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type gameServer interface {
	CreateGame(ctx context.Context) (user2 string, err error)
	DelGame() error
	JoinGame(creatorUserName string) error
	GetAvailableGames() (games []string, err error)
	SaveGameResult(winner, loser string) error
	GetUserStat(username string) (domain.Statistics, error)
	GetOpponentName() (string, error)
}

type gameBattle interface {
	StartBattle() error
	Ready() (iFirst bool, err error)
	Attack(x, y int) (msgToUser string, err error)
	Defend() (msgToUser string, err error)
}

type gameMap interface {
	StringMap(myMap bool) [10]string
	PlaceShip(x, y int, sType gameSrvs.ShipType, direction gameSrvs.ShipDirection) error
	CanAttack(x, y int) bool
	AllShipsDestroyed() bool
	AllShipsPlaced() bool
	GetAvailCntShipType(sType gameSrvs.ShipType) int
}

type gameService interface {
	gameBattle
	gameMap
	gameServer
}

type GameUI struct {
	game gameService
}

func New(gService gameService) *GameUI {
	return &GameUI{
		game: gService,
	}
}

func (g *GameUI) printMap(my bool) {
	myMap := g.game.StringMap(my)
	if my {
		fmt.Println("Your map: ")
	} else {
		fmt.Println("Opponent map: ")
	}
	for i := range myMap {
		fmt.Println(string(rune('A'+i)), "|", myMap[i], "|")
	}
	fmt.Println("  | 0 1 2 3 4 5 6 7 8 9 |")
}

func (g *GameUI) StartGame(userName string) {
	userStats, err := g.game.GetUserStat(userName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Your statistics: ", userStats)

	battleStarted := false
	for !battleStarted {
		fmt.Println("Select command: \n1. Create game\n2. Get available games\n3. Exit\nEnter number of command: ")
		var command int
		cntScan, err := fmt.Scan(&command)
		if err != nil || cntScan != 1 {
			fmt.Println("Invalid input")
		} else {
			switch command {
			case 1:
				ctx, _ := context.WithCancel(context.Background())
				var user2Name string
				opponentWait := make(chan error)

				go func() {
					user2Name, err = g.game.CreateGame(ctx)
					if err != nil {
						opponentWait <- err
					}
					close(opponentWait)
				}()

				fmt.Println("Waiting for the opponent to join...")

				//fmt.Println("Cancel waiting with Ctrl+C")
				//stop := make(chan os.Signal, 1)
				//signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

				select {
				case err = <-opponentWait:
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("Opponent found: ", user2Name)
					err = g.game.StartBattle()
					if err != nil {
						fmt.Println(err)
					} else {
						battleStarted = true
					}
					//case <-stop:
					//	cancel()
				}
			case 2:
				games, err := g.game.GetAvailableGames()
				if err != nil {
					fmt.Println(err)
					continue
				}
				if len(games) == 0 {
					fmt.Println("No available games")
					continue
				}
				fmt.Println("Available games: ")
				for i, game := range games {
					fmt.Println(i+1, ". ", game)
				}
				for { // while invalid input
					fmt.Println("Enter number of game: ")
					var gameNumber int
					cntScan, err := fmt.Scan(&gameNumber)
					if err != nil || cntScan != 1 || gameNumber < 1 || gameNumber > len(games) {
						fmt.Println("Invalid input")
					} else {
						err = g.game.JoinGame(games[gameNumber-1])
						if err != nil {
							fmt.Println(err)
						} else {
							err = g.game.StartBattle()
							if err != nil {
								fmt.Println(err)
								os.Exit(1)
							}
							battleStarted = true
							break
						}
					}
				}
			case 3:
				os.Exit(0)
			}
		}
	}
}

func (g *GameUI) ReadyToBattle() {
	fmt.Println("It's time to place the ships")
	curShipType := gameSrvs.FourDeck
	for !g.game.AllShipsPlaced() {
		g.printMap(true)
		if g.game.GetAvailCntShipType(curShipType) == 0 {
			curShipType--
		}
		for {
			switch curShipType {
			case gameSrvs.FourDeck:
				fmt.Println("Place four-deck ship")
			case gameSrvs.ThreeDeck:
				fmt.Println("Place three-deck ship")
			case gameSrvs.TwoDeck:
				fmt.Println("Place two-deck ship")
			case gameSrvs.SingleDeck:
				fmt.Println("Place one-deck ship")
			}
			var x, y int
			fmt.Println("Select coordinate: ")
			var coordinate string
			cntScan, err := fmt.Scan(&coordinate)
			if err != nil || cntScan != 1 {
				fmt.Println("Invalid input")
				continue
			}

			coordinate = strings.ToUpper(coordinate)
			if len(coordinate) != 2 {
				fmt.Println("Invalid coordinate")
				continue
			}
			y = int(coordinate[0] - 'A')
			x, err = strconv.Atoi(string(coordinate[1]))
			if err != nil || x < 0 || x > 9 || y < 0 || y > 9 {
				fmt.Println("Invalid coordinate")
				continue
			}

			fmt.Println("Select direction: \n1. Up\n2. Down\n3. Left\n4. Right\nEnter number of direction: ")
			var direction int
			cntScan, err = fmt.Scan(&direction)
			if err != nil || cntScan != 1 {
				fmt.Println("Invalid input")
				continue
			}
			err = g.game.PlaceShip(x, y, curShipType, gameSrvs.ShipDirection(direction-1))
			if err != nil {
				fmt.Println(err)
				continue
			}
			break
		}
	}
}

func (g *GameUI) StartBattle() (win bool) {
	fmt.Println("Ready!")
	IFirst, err := g.game.Ready()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Battle started")
	if IFirst {
		for !g.game.AllShipsDestroyed() {
			g.printMap(false)
			//Attack
			var x, y int
			for {
				fmt.Println("Select coordinate: ")
				cntScan, err := fmt.Scan(&x, &y)
				if err != nil || cntScan != 2 {
					fmt.Println("Invalid input")
					continue
				}
				if x < 0 || x > 9 || y < 0 || y > 9 {
					fmt.Println("Invalid coordinate")
					continue
				}
				if !g.game.CanAttack(y, x) {
					fmt.Println("Already attacked")
					continue
				}
				break
			}
			msg, err := g.game.Attack(y, x)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(msg)
			if msg == gameSrvs.Win {
				break
			}
			g.printMap(false)

			//Defend
			msg, err = g.game.Defend()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(msg)
			if msg == gameSrvs.Lose {
				break
			}
			g.printMap(true)
		}
	} else {
		for !g.game.AllShipsDestroyed() {
			//Defend
			msg, err := g.game.Defend()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(msg)
			if msg == gameSrvs.Lose {
				win = false
				break
			}

			//Attack
			var x, y int
			for {
				fmt.Println("Select coordinate: ")
				cntScan, err := fmt.Scan(&x, &y)
				if err != nil || cntScan != 2 {
					fmt.Println("Invalid input")
					continue
				}
				if x < 0 || x > 9 || y < 0 || y > 9 {
					fmt.Println("Invalid coordinate")
					continue
				}
				break
			}
			msg, err = g.game.Attack(x, y)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(msg)
			if msg == gameSrvs.Win {
				win = true
				break
			}
		}
	}
	return win
}

func (g *GameUI) GetOpponentName() string {
	user2Name, _ := g.game.GetOpponentName()
	return user2Name
}

func (g *GameUI) SendResult(user1Name, user2Name string) {
	err := g.game.SaveGameResult(user1Name, user2Name)
	if err != nil {
		fmt.Println(err)
	}
}
