package chapter5

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"time"
)

type App struct {
	conn *redis.Client
}

func NewApp(conn *redis.Client) *App {
	return &App{
		conn: conn,
	}
}

var (
	Debug    = "debug"
	Info     = "info"
	Warning  = "warning"
	Error    = "error"
	Critical = "critical"
)

func (a *App) LogRecent(name, message, severity string, pipe redis.Pipeliner) error {
	if severity == "" {
		severity = Info
	}

	destination := fmt.Sprintf("recent:%s:%s", name, severity)
	message = time.Now().Local().String() + " " + message
	if pipe == nil {
		pipe = a.conn.Pipeline()
	}
	pipe.LPush(destination, message)
	pipe.LTrim(destination, 0, 99)
	_, err := pipe.Exec()
	return err
}

func (a *App) LogCommon(name, message, severity string, timeout int64) {
	if severity == "" {
		severity = Info
	}
	destination := fmt.Sprintf("common:%s:%s", name, severity)
	startKey := destination + ":start"
	end := time.Now().Unix() + timeout
	for time.Now().Unix() < end {
		if err := a.conn.Watch(func(tx *redis.Tx) error {
			pipe := a.conn.Pipeline()
			hourStart := time.Now().Hour()
			existing, _ := pipe.Get(startKey).Int()
			if existing != 0 && existing < hourStart {
				pipe.Rename(destination, destination+":last")
				pipe.Rename(startKey, destination+":pstart")
				pipe.Set(startKey, hourStart, 0)
			}
			pipe.ZIncrBy(destination, 1, message)
			if err := a.LogRecent(name, message, severity, pipe); err != nil {
				log.Printf("Error %+v", err)
			}
			return nil
		}, startKey); err != nil {
			log.Printf("Error: %+v", err)
		}
	}
}
