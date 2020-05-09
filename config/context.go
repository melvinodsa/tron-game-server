// Copyright 2020 Melvin Davis. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

/* This file contains the definition of AppContext */

//AppContext contains the
type AppContext struct {
	//Log for logging purposes
	Log Logger
}

var rootAppContext *AppContext

func init() {
	/*
	 * We will initialize the context
	 * We will connect to the database
	 */
	rootAppContext = &AppContext{}
}

//NewAppContext returns an initlized app context
func NewAppContext(l Logger) *AppContext {
	return &AppContext{Log: l}
}
