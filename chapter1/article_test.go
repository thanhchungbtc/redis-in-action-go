package chapter1

import (
	"github.com/go-redis/redis/v7"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Chapter 1", func() {
	var a *App
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	BeforeEach(func() {
		a = NewApp(client)
	})
	AfterEach(func() {
		a.conn.FlushAll()
	})

	Describe("ArticleVote", func() {

		It("Should work", func() {
			// prepare
			article := "article:10048"
			articleID := "10048"
			user := "btc"
			now := time.Now().Unix()
			a.conn.ZAdd("time:", &redis.Z{Member: article, Score: float64(now)})

			// execute
			a.ArticleVote(user, article)

			// assert
			Expect(a.conn.SCard("voted:"+articleID).Val() == 1).To(Equal(true))
		})
	})

})
