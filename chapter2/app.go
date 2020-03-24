package chapter2

import "github.com/go-redis/redis/v7"

var (
	Quit        = false
	Limit int64 = 1e7
)

type App struct {
	conn *redis.Client
}

func NewApp(conn *redis.Client) *App {
	return &App{
		conn: conn,
	}
}
