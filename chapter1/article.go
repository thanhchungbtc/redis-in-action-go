package chapter1

import (
	"github.com/go-redis/redis/v7"
	"strings"
	"time"
)

var (
	OneWeekInSeconds int64 = 7 * 86400
	VoteScore        int64 = 432
	ArticlesPerPage  int64 = 25
)

type App struct {
	conn *redis.Client
}

func NewApp(conn *redis.Client) *App {
	return &App{
		conn: conn,
	}
}

func (a *App) ArticleVote(user, article string) {
	// only recent articles can be voted
	cutOff := time.Now().Unix() - OneWeekInSeconds
	if a.conn.ZScore("time:", article).Val() < float64(cutOff) {
		return
	}

	// e.g. article:10040
	articleID := strings.Split(article, ":")[1]
	if a.conn.SAdd("voted:"+articleID, user).Val() != 0 {
		a.conn.ZIncrBy("score:", float64(VoteScore), article)
		a.conn.HIncrBy(article, "votes", 1)
	}
}

func (a *App) PostArticle(user, title, link string) string {
	articleID := string(a.conn.Incr("article:").Val())
	voted := "voted:" + articleID

	a.conn.SAdd(voted, user)
	now := time.Now().Unix()
	article := "article:" + articleID
	a.conn.HMSet(article, map[string]interface{}{
		"title":  title,
		"link":   link,
		"poster": user,
		"time":   now,
		"votes":  1,
	})
	a.conn.ZAdd("score:", &redis.Z{Member: article, Score: float64(now + VoteScore)})
	a.conn.ZAdd("time:", &redis.Z{Member: article, Score: float64(now)})

	return articleID
}

func (a *App) GetArticles(page int, order string) []map[string]string {
	start := int64(page-1) * ArticlesPerPage
	end := int64(start) + ArticlesPerPage - 1
	ids := a.conn.ZRevRange(order, start, end).Val()
	articles := []map[string]string{}
	for _, id := range ids {
		articleData := a.conn.HGetAll(id).Val()
		articleData["id"] = id
		articles = append(articles, articleData)
	}
	return articles

}

func (a *App) AddRemoveGroups(articleID string, toAdd, toRemove []string) {
	article := "article:" + articleID
	for _, group := range toAdd {
		a.conn.SAdd("group:"+group, article)
	}
	for _, group := range toRemove {
		a.conn.SRem("group:"+group, article)
	}
}

func (a *App) GetGroupArticles(group string, page int, order string) []map[string]string {
	key := order + group
	if a.conn.Exists(key).Val() == 0 {
		a.conn.ZInterStore(key, &redis.ZStore{
			Aggregate: "Max",
			Keys:      []string{"group:" + group, order},
		})
		a.conn.Expire(key, 60)
	}
	return a.GetArticles(page, key)

}
