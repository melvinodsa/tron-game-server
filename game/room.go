// Copyright 2020 Melvin Davis. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

//Room has the game
type Room struct {
	ID           uuid.UUID
	players      []Player
	Status       Status
	ended        bool
	lastActivity time.Time
	board        *Board
	admin        Player
	in           chan DirectionUpdate
	out          chan Update
}

//NewRoom creates a new room
func NewRoom(p Player) *Room {
	return &Room{
		ID:           uuid.New(),
		players:      []Player{p},
		Status:       Waiting,
		ended:        false,
		lastActivity: time.Now(),
		admin:        p,
		in:           make(chan DirectionUpdate),
		out:          make(chan Update),
	}
}

//Join will add a user to the room
func (r *Room) Join(p Player) {
	r.players = append(r.players, p)
}

//Exit will remove a player from the room
func (r *Room) Exit(p Player) {
	id := p.ID.String()
	for i := len(r.players) - 1; i >= 0; i-- {
		if strings.Compare(r.players[i].ID.String(), id) == 0 {
			r.players = append(r.players[:i], r.players[i+1:]...)
			break
		}
	}
}

//Participate will enable the player to participate the game.
//It returns true if ther user was able to be added to a board position
func (r *Room) Participate(p Player, x, y int, color string) bool {
	if r.Status != Waiting {
		return false
	}
	return r.board.Add(p, x, y, color)
}

//SetBoard will set the board for the game and sets the status to waiting
func (r *Room) SetBoard(size, speed int) {
	if r.Status == Started {
		return
	}
	r.board = NewBoard(size, speed)
	r.Status = Waiting
}

//StartGame will start the game if one isn't in progression
func (r *Room) StartGame() {
	if r.Status != Waiting || r.board == nil {
		return
	}
	r.Status = Started
	ch := make(chan Update)
	end := time.NewTicker(r.board.Life()).C
	go r.board.Start(r.in, ch)
	for {
		select {
		case update := <-ch:
			if update.Status == Winner {
				r.Status = Ended
				go SendUpdate(r.out, update)
				return
			}
			go SendUpdate(r.out, update)
		case <-end:
			return
		}
	}
}
