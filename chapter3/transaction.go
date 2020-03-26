package chapter3

import (
	"log"
	"time"
)

func (a *App) NoTrans() {
	a.conn.IncrBy("notrans:", 1)
	time.Sleep(time.Millisecond * 100)
	a.conn.IncrBy("notrans:", -1)
}

func (a *App) Trans() {
	pipeline := a.conn.Pipeline()
	pipeline.IncrBy("trans:", 1)
	time.Sleep(time.Millisecond * 100)
	pipeline.IncrBy("trans:", -1)
	_, err := pipeline.Exec()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
