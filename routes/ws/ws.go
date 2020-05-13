// Copyright 2020 Melvin Davis. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package ws has the websockets handler implementation for the server
package ws

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/melvinodsa/tron-game-server/config"
	"github.com/melvinodsa/tron-game-server/game"
	"github.com/melvinodsa/tron-game-server/routes"

	socketio "github.com/googollee/go-socket.io"
)

//WebSockets is the websockets handler for the application
func WebSockets(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appCtx := ctx.Value(routes.AppContextKey).(*config.AppContext)
	appCtx.Log.Info("Got a websockets connection request")
	socketioServer.ServeHTTP(res, req)
}

type user struct {
	Player game.Player
	Room   game.Room
}

var socketioServer *socketio.Server

func init() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "name", func(s socketio.Conn, msg string) {
		player := game.NewPlayer(msg)
		u := user{Player: *player}
		s.SetContext(&u)
		s.Emit("player", player)
	})
	server.OnEvent("/", "join", func(s socketio.Conn, msg string) string {
		u := s.Context()
		if u == nil {
			return `{"error":true, "msg": "username not set"}`
		}
		id, err := uuid.Parse(msg)
		if err != nil {
			return `{"error": true, "msg": "invalid room id"}`
		}
		req := game.Request{ID: id, Type: game.Get, Out: make(chan game.Request)}
		go game.SendRequest(game.GameRequestChannel, req)
		res := <-req.Out
		if !res.Valid {
			return `{"error":true, "msg": "room not found"}`
		}
		return `{"error":false, "msg": "joined the room"}`
	})
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	socketioServer = server
	go socketioServer.Serve()
}

func init() {
	routes.AddRoutes(routes.Route{
		Version:     "v1",
		HandlerFunc: WebSockets,
		Pattern:     "/socket.io/",
	})
}
