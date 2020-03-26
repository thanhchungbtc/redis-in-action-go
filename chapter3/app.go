package chapter3

import "github.com/go-redis/redis/v7"

type App struct {
	conn *redis.Client
}

func NewApp(conn *redis.Client) *App {
	return &App{
		conn: conn,
	}
}
