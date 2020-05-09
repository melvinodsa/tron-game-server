// Copyright 2020 Melvin Davis. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package game has the definitions and implementations required by the game
package game

import "github.com/google/uuid"

//RequesType is the type of request for the game rooms
type RequesType uint

const (
	//New to create a game room
	New RequesType = 1
	//Get to fetch an existing game room
	Get RequesType = 2
)

//Request helps communicating with the game room
type Request struct {
	ID     uuid.UUID
	Type   RequesType
	Player Player
	Room   *Room
	Out    chan Request
	Valid  bool
}

//SendRequest sends the request to a channel
func SendRequest(ch chan Request, req Request) {
	ch <- req
}

//GameRequestChannel is the common request channel for game room
var GameRequestChannel = make(chan Request)

//Game routine controls the games
func Game(in chan Request) {
	rooms := map[uuid.UUID]*Room{}
	for {
		req := <-in
		switch req.Type {
		case New:
			room := NewRoom(req.Player)
			rooms[room.ID] = room
			req.Room = room
			go SendRequest(req.Out, req)
		case Get:
			req.Room, req.Valid = rooms[req.ID]
			go SendRequest(req.Out, req)
		}
	}
}

func init() {
	go Game(GameRequestChannel)
}
