package chapter2

import (
	"math"
	"time"
)

func (a *App) AddToCart(session, item string, count int) {
	if count <= 0 {
		a.conn.HDel("cart:"+session, item)
	} else {
		a.conn.HSet("cart:"+session, item, count)
	}
}

func (a*App) CleanFullSessions() {
	for !Quit {
		size := a.conn.ZCard("recent:").Val()
		if size <= Limit {
			time.Sleep(time.Second)
			continue
		}

		endIndex := int64(math.Min(float64(size - Limit), float64(100)))
		sessions := a.conn.ZRange("recent:", 0, endIndex - 1).Val()
		var sessionKeys []string
		for _, sess := range sessions {
			sessionKeys = append(sessionKeys, "viewed:"+sess)
			sessionKeys = append(sessionKeys, "cart:"+sess)
		}
		a.conn.Del(sessionKeys...)
		a.conn.HDel("login:", sessionKeys...)
		a.conn.ZRem("recent:", sessions)
	}
}