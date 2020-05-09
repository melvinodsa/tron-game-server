// Copyright 2020 Melvin Davis. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package web has the web request handler implementation for the server
package web

import (
	"context"
	"net/http"

	"github.com/melvinodsa/tron-game-server/routes"
)

//Home returns the home page of the game
func Home(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "static/home.html")
}

func init() {
	routes.AddRoutes(routes.Route{
		Version:     "v1",
		HandlerFunc: Home,
		Pattern:     "/",
	})
}
