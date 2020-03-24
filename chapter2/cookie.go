package chapter2

import (
	"github.com/go-redis/redis/v7"
	"math"
	"time"
)


func (a *App) CheckToken(token string) string {
	return a.conn.HGet("login:", token).Val()
}

func (a *App) UpdateToken(token, user, item string) {
	timestamp := time.Now().Unix()
	a.conn.HSet("login:", token, user)
	a.conn.ZAdd("recent:", &redis.Z{
		Member: token,
		Score:  float64(timestamp),
	})
	if item != "" {
		a.conn.ZAdd("viewed:"+token, &redis.Z{
			Member: item,
			Score:  float64(timestamp),
		})
		// keeping the most 25 recent items
		a.conn.ZRemRangeByRank("viewed:"+token, 0, -26)
	}
}

func (a *App) CLeanSessions() {
	for !Quit {
		size := a.conn.ZCard("recent:").Val()
		if size <= Limit {
			time.Sleep(time.Second)
			continue
		}
		endIndex := int64(math.Min(float64(size-Limit), 1000))
		tokens := a.conn.ZRange("recent:", 0, endIndex-1).Val()
		var sessionKeys []string
		for _, token := range tokens {
			sessionKeys = append(sessionKeys, "viewed:"+token)
		}
		a.conn.Del(sessionKeys...)
		a.conn.HDel("login:", tokens...)
		a.conn.ZRem("recent:", tokens)
	}
}

