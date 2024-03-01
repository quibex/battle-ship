package gameSrvs

import (
	"battlship/internal/adapters/rabbitmq"
	"errors"
	"fmt"
	"strconv"
)

type gameMQ interface {
	serverMQ
	battleMQ
}

const seaSize = 10

type BattleShip struct {
	mySea        [seaSize][seaSize]SeaCell
	opponentSea  [seaSize][seaSize]SeaCell
	ships        [4]int
	mq           gameMQ
	opponentMsgs <-chan rabbitmq.Message
}

var maxShips = [4]int{4, 3, 2, 1}

func New(mq gameMQ) *BattleShip {
	emptySea := [seaSize][seaSize]SeaCell{}
	for i := 0; i < seaSize; i++ {
		for j := 0; j < seaSize; j++ {
			emptySea[i][j] = emptyCell
		}
	}
	unknownSea := [seaSize][seaSize]SeaCell{}
	for i := 0; i < seaSize; i++ {
		for j := 0; j < seaSize; j++ {
			unknownSea[i][j] = unknownCell
		}
	}
	return &BattleShip{
		mySea:       emptySea,
		opponentSea: unknownSea,
		mq:          mq,
	}
}

type SeaCell rune

const ( //Sea symbols
	emptyCell   SeaCell = '~'
	shipCell    SeaCell = '▢'
	hitCell     SeaCell = 'x'
	missCell    SeaCell = '•'
	unknownCell SeaCell = '□'
)

type ShipDirection uint8

const ( //direction
	up ShipDirection = iota
	down
	left
	right
)

type ShipType uint

const ( //ShipType
	SingleDeck ShipType = iota
	TwoDeck
	ThreeDeck
	FourDeck
)

func (b *BattleShip) StringMap(myMap bool) [10]string {
	var sea *[seaSize][seaSize]SeaCell
	if myMap {
		fmt.Println("Print my sea: ")
		sea = &b.mySea
	} else {
		fmt.Println("Print opponent sea: ")
		sea = &b.opponentSea
	}
	var seaMap [seaSize]string
	for i := 0; i < seaSize; i++ {
		for j := 0; j < seaSize; j++ {
			if rune(sea[i][j]) >= '0' && rune(sea[i][j]) <= '9' { // == shipCell
				seaMap[i] += string(shipCell)
			} else {
				seaMap[i] += fmt.Sprint(string(b.mySea[i][j]))
			}
			if j != seaSize-1 {
				seaMap[i] += " "
			}
		}
	}
	return seaMap
}

func (b *BattleShip) PlaceShip(x, y int, sType ShipType, direction ShipDirection) error {
	if b.ships[sType] == maxShips[sType] {
		return errors.New("count shipCell of this type over flow max")
	}
	if !b.canPlaceShip(x, y, sType, direction) {
		return errors.New("fields is already occupied")
	}

	switch direction {
	case up:
		for i := y; i > y-(int(sType)+1); i-- {
			b.mySea[i][x] = SeaCell(strconv.Itoa(int(sType))[0])
		}
	case down:
		for i := y; i < y+int(sType)+1; i++ {
			b.mySea[i][x] = SeaCell(strconv.Itoa(int(sType))[0])
		}
	case right:
		for j := x; j < x+int(sType)+1; j++ {
			b.mySea[y][j] = SeaCell(strconv.Itoa(int(sType))[0])
		}
	case left:
		for j := x; j > x-(int(sType)+1); j-- {
			b.mySea[y][j] = SeaCell(strconv.Itoa(int(sType))[0])
		}
	}
	b.ships[sType]++
	return nil
}

func (b *BattleShip) canPlaceShip(x, y int, sType ShipType, direction ShipDirection) bool {
	switch direction {
	case up:
		if y-(int(sType)+1) < 0 {
			return false
		}
		for i := y; i >= y-(int(sType)+1); i-- {
			if b.mySea[i][x] != emptyCell {
				return false
			}
			if x+1 < seaSize && b.mySea[i][x+1] != emptyCell {
				return false
			}
			if x-1 >= 0 && b.mySea[i][x-1] != emptyCell {
				return false
			}
			if i+1 < seaSize && b.mySea[i+1][x] != emptyCell {
				return false
			}
			if i-1 >= 0 && b.mySea[i-1][x] != emptyCell {
				return false
			}
		}
	case down:
		if y+(int(sType)+1) >= seaSize {
			return false
		}
		for i := y; i < y+(int(sType)+1); i++ {
			if b.mySea[i][x] != emptyCell {
				return false
			}
			if x+1 < seaSize && b.mySea[i][x+1] != emptyCell {
				return false
			}
			if x-1 >= 0 && b.mySea[i][x-1] != emptyCell {
				return false
			}
			if i+1 < seaSize && b.mySea[i+1][x] != emptyCell {
				return false
			}
			if i-1 >= 0 && b.mySea[i-1][x] != emptyCell {
				return false
			}
		}
	case right:
		if x+(int(sType)+1) >= seaSize {
			return false
		}
		for j := x; j < x+(int(sType)+1); j++ {
			if b.mySea[y][j] != emptyCell {
				return false
			}
			if y+1 < seaSize && b.mySea[y+1][j] != emptyCell {
				return false
			}
			if y-1 >= 0 && b.mySea[y-1][j] != emptyCell {
				return false
			}
			if j+1 < seaSize && b.mySea[y][j+1] != emptyCell {
				return false
			}
			if j-1 >= 0 && b.mySea[y][j-1] != emptyCell {
				return false
			}
		}
	case left:
		if x-(int(sType)+1) < 0 {
			return false
		}
		for j := x; j >= x-(int(sType)+1); j-- {
			if b.mySea[y][j] != emptyCell {
				return false
			}
			if y+1 < seaSize && b.mySea[y+1][j] != emptyCell {
				return false
			}
			if y-1 >= 0 && b.mySea[y-1][j] != emptyCell {
				return false
			}
			if j+1 < seaSize && b.mySea[y][j+1] != emptyCell {
				return false
			}
			if j-1 >= 0 && b.mySea[y][j-1] != emptyCell {
				return false
			}
		}
	}
	return true
}

func (b *BattleShip) GetAvailCntShipType(sType ShipType) int {
	return maxShips[sType] - b.ships[sType]
}

func (b *BattleShip) CanAttack(x, y int) bool {
	if b.opponentSea[y][x] != unknownCell {
		return false
	}
	return true
}

func (b *BattleShip) hit(x, y int) (hit, destroy bool) {
	if b.mySea[y][x] < '0' && b.mySea[y][x] > '9' { // != shipCell
		hit, destroy = false, false
		return
	}
	cntHitCell := 1 + b.cntHitCells(x+1, y, right) + b.cntHitCells(x-1, y, left) + b.cntHitCells(x, y+1, down) + b.cntHitCells(x, y-1, up)
	if cntHitCell == int(b.mySea[x][y]) { //the entire cell of the ship was hit
		hit, destroy = true, true
		return
	}
	hit, destroy = true, false
	return
}

// the direction to eliminate the cycle
func (b *BattleShip) cntHitCells(x, y int, direction ShipDirection) int {
	if x < 0 || x >= seaSize || y < 0 || y >= seaSize || b.mySea[y][x] != hitCell {
		return 0
	}
	switch direction {
	case up:
		return 1 + b.cntHitCells(x, y-1, up)
	case down:
		return 1 + b.cntHitCells(x, y+1, down)
	case right:
		return 1 + b.cntHitCells(x+1, y, right)
	case left:
		return 1 + b.cntHitCells(x-1, y, left)
	}
	return -1
}

func (b *BattleShip) markHitOrMiss(x, y int, hit, destroy bool, mySea bool) {
	var sea *[seaSize][seaSize]SeaCell
	if mySea {
		sea = &b.mySea
	} else {
		sea = &b.opponentSea
	}
	if hit {
		b.mySea[y][x] = hitCell
	}
	if destroy {
		if mySea {
			b.ships[int(b.mySea[y][x])]--
		}
		markDestroy(x+1, y, sea, right)
		markDestroy(x-1, y, sea, left)
		markDestroy(x, y-1, sea, up)
		markDestroy(x, y+1, sea, down)
	}
}

// the direction to eliminate the cycle
func markDestroy(x, y int, sea *[seaSize][seaSize]SeaCell, direction ShipDirection) {
	if x < 0 || x >= seaSize || y < 0 || y >= seaSize {
		return
	}
	if sea[y][x] == hitCell {
		if direction != left {
			markDestroy(x+1, y, sea, right)
		}
		if direction != right {
			markDestroy(x-1, y, sea, left)
		}
		if direction != up {
			markDestroy(x, y+1, sea, down)
		}
		if direction != down {
			markDestroy(x, y-1, sea, up)
		}
	} else {
		sea[y][x] = missCell
	}
}

func (b *BattleShip) AllShipsDestroyed() bool {
	for _, cnt := range b.ships {
		if cnt != 0 {
			return false
		}
	}
	return true
}

func (b *BattleShip) AllShipsPlaced() bool {
	for i := range b.ships {
		if b.ships[i] != maxShips[i] {
			return false
		}
	}
	return true
}
