package gocfm

import "github.com/ryouaki/koa/session"

var SessionStore session.MemStore = *session.NewMemStore()
