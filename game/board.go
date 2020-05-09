// Copyright 2020 Melvin Davis. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

//Board is the game board on which will be played
type Board struct {
	positions        [][]int
	speed            int
	size             int
	playerMap        map[uuid.UUID]int
	colorMap         map[int]string
	currentPositions map[int][]int
	directionMap     map[int]Direction
	status           Status
	l                *sync.Mutex
	debug            bool
}

//NewBoard will initialize a board game
func NewBoard(size int, speed int) *Board {
	m := [][]int{}
	for i := 0; i < size; i++ {
		res := []int{}
		for j := 0; j < size; j++ {
			res = append(res, 0)
		}
		m = append(m, res)
	}
	if size < 10 {
		size = 10
	}
	if speed <= 0 {
		speed = 1
	}
	debug := false
	if os.Getenv("DEBUG") == "true" {
		debug = true
	}
	return &Board{
		positions:        m,
		speed:            speed,
		playerMap:        map[uuid.UUID]int{},
		size:             size,
		directionMap:     map[int]Direction{},
		colorMap:         map[int]string{},
		currentPositions: map[int][]int{},
		status:           Waiting,
		l:                &sync.Mutex{},
		debug:            debug,
	}
}

//String method has the stringer implementation of the board
func (b Board) String() string {
	for r, v := range b.positions {
		for range v {
			fmt.Print(" --- ")
		}
		fmt.Print("\n")
		for c, val := range v {
			fmt.Printf("| %d ", val)
			if c == len(v)-1 {
				fmt.Print("|")
			}
		}
		fmt.Print("\n")
		if r != len(b.positions)-1 {
			continue
		}
		for range v {
			fmt.Print(" --- ")
		}
		fmt.Print("\n")
	}
	return ""
}

//Size of the board
func (b Board) Size() int {
	return b.size
}

//Life is the maximum time a game can exist
func (b Board) Life() time.Duration {
	return time.Duration(b.size*b.size) * time.Second / time.Duration(b.speed)
}

//ColorMap returns the colormap of the board
func (b Board) ColorMap() map[int]string {
	return b.colorMap
}

//PlayerMap returns the player map of board
func (b Board) PlayerMap() map[uuid.UUID]int {
	return b.playerMap
}

//StartingPositions returns the starting positions of each player
func (b Board) StartingPositions() map[int][]int {
	return b.currentPositions
}

//DirectionMap returns the direction map of the players
func (b Board) DirectionMap() map[int]Direction {
	return b.directionMap
}

//Add will add a player to the board at a given position
func (b *Board) Add(p Player, x, y int, color string) bool {
	if b.status != Waiting {
		return false
	}
	if x < 0 || x >= b.size || y < 0 || y >= b.size {
		return false
	}
	b.l.Lock()
	if _, ok := b.playerMap[p.ID]; ok {
		return false
	}
	if b.positions[y][x] != 0 {
		return false
	}
	index := len(b.playerMap) + 1
	b.positions[y][x] = index
	b.playerMap[p.ID] = index
	b.colorMap[index] = color
	b.currentPositions[index] = []int{x, y}
	b.setDirection(index)
	b.l.Unlock()
	return true
}

//UpdatePosition will update the starting position of the player on the board
func (b *Board) UpdatePosition(p Player, x, y int) bool {
	if b.status != Waiting {
		return false
	}
	if x < 0 || x >= b.size || y < 0 || y >= b.size {
		return false
	}
	b.l.Lock()
	index, ok := b.playerMap[p.ID]
	if !ok {
		return false
	}
	if b.positions[y][x] != 0 {
		return false
	}
	sp := b.currentPositions[index]
	b.positions[sp[1]][sp[0]] = 0
	b.positions[y][x] = index
	b.currentPositions[index] = []int{x, y}
	b.setDirection(index)
	b.l.Unlock()
	return true
}

func (b *Board) setDirection(index int) {
	sp, ok := b.currentPositions[index]
	if !ok {
		return
	}
	y1, y2, x1, x2 := sp[1], b.size-sp[1], sp[0], b.size-sp[0]
	maxY, up := y1, true
	if y2 > y1 {
		maxY = y2
		up = false
	}
	maxX, left := x1, true
	if x2 > x1 {
		maxX = x2
		left = false
	}
	y := true
	dir := Up
	if maxX > maxY {
		y = false
		dir = Left
	}
	if y && !up {
		dir = Down
	} else if !y && !left {
		dir = Right
	}
	b.directionMap[index] = dir
}

//Start the board timers and game
func (b *Board) Start(in chan DirectionUpdate, out chan Update) {
	if b.status != Waiting {
		return
	}
	b.status = Started
	timers := time.NewTicker(time.Second / time.Duration(b.speed)).C
	for {
		select {
		case <-timers:
			//do  move
			toBeRemoved := []uuid.UUID{}
			lastPlayer := 0
			lastStanding := 0
			for k, v := range b.playerMap {
				ok := b.Move(v)
				if !ok {
					go SendUpdate(out, Update{Player: v, Status: Crashed})
					toBeRemoved = append(toBeRemoved, k)
					lastPlayer = v
				} else {
					lastStanding = v
				}
			}
			if b.debug {
				fmt.Println(b)
			}
			if len(toBeRemoved) == 0 {
				break
			}
			for _, v := range toBeRemoved {
				i := b.playerMap[v]
				delete(b.playerMap, v)
				delete(b.colorMap, i)
				delete(b.directionMap, i)
				delete(b.currentPositions, i)
			}
			if len(b.playerMap) == 1 {
				b.status = Ended
				go SendUpdate(out, Update{Player: lastStanding, Status: Winner})
				return
			}
			if len(b.playerMap) == 0 {
				b.status = Ended
				go SendUpdate(out, Update{Player: lastPlayer, Status: Winner})
				return
			}
		case dir := <-in:
			//update direction
			_, ok := b.directionMap[dir.Player]
			if !ok {
				break
			}
			b.directionMap[dir.Player] = dir.Direction
		}
	}
}

//Move will make the player at the given index move as per his current direction
func (b *Board) Move(index int) bool {
	pos, ok := b.currentPositions[index]
	dir, dOk := b.directionMap[index]
	if !ok || !dOk {
		return false
	}
	x, y := pos[0], pos[1]
	if dir == Up {
		y--
	} else if dir == Down {
		y++
	} else if dir == Left {
		x--
	} else if dir == Right {
		x++
	}
	b.currentPositions[index] = []int{x, y}
	if y < 0 || y >= b.size || x < 0 || x >= b.size {
		return false
	}
	v := b.positions[y][x]
	if v != 0 {
		return false
	}
	b.positions[y][x] = index
	return true
}
