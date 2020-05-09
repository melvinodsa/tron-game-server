// Copyright 2020 Melvin Davis. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import "github.com/google/uuid"

//Player is the induvidual player in the game
type Player struct {
	ID   uuid.UUID
	Name string
}

//NewPlayer returns a new player
func NewPlayer(name string) *Player {
	return &Player{ID: uuid.New(), Name: name}
}
