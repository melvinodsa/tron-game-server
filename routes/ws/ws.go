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
		s.Emit("reply", player.Name)
		req := game.Request{Player: *player, Type: game.New, Out: make(chan game.Request)}
		go game.SendRequest(game.GameRequestChannel, req)
		req = <-req.Out
		room := req.Room
		room.SetBoard(10, 2)
		player2 := game.NewPlayer("new2")
		room.Participate(*player, 0, 0, "#acacac")
		room.Participate(*player2, 5, 5, "#acacac")
		go room.StartGame()
	})
	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
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
