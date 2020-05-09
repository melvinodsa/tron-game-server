// Copyright 2020 Melvin Davis. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

//Status of the game
type Status uint

const (
	//Waiting for the game to start
	Waiting Status = 1
	//Started indicates the game has started
	Started Status = 2
	//Ended indicates that a round of game has been ended
	Ended Status = 3
)

//Direction in which a player move
type Direction uint

const (
	//Up means the player is moving up
	Up Direction = 1
	//Down means the player is moving down
	Down Direction = 2
	//Right means the player is moving right
	Right Direction = 3
	//Left means the player is moving left
	Left Direction = 4
)

//DirectionUpdate will give the direction update for a player
type DirectionUpdate struct {
	Player    int
	Direction Direction
}

//PlayerStatus has status of a player in a game
type PlayerStatus uint

const (
	//Alive means the player is still in the game
	Alive PlayerStatus = 1
	//Crashed means the player crashed out of the game
	Crashed PlayerStatus = 2
	//Winner means the player won the game
	Winner PlayerStatus = 3
)

//Update will give game update events
type Update struct {
	Player int
	Status PlayerStatus
}

//SendUpdate sends an update to an channel
func SendUpdate(ch chan Update, update Update) {
	ch <- update
}
