package chapter4

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

// ListItem lists an item to the market
func (a *App) ListItem(itemID, sellerID string, price float64) bool {
	inventory := fmt.Sprintf("inventory: %s", sellerID)
	item := fmt.Sprintf("%s.%s", itemID, sellerID)
	end := time.Now().Unix() + 5

	for time.Now().Unix() < end {
		err := a.conn.Watch(func(tx *redis.Tx) error {

			if _, err := tx.Pipelined(func(pipeliner redis.Pipeliner) error {
				if !tx.SIsMember(inventory, itemID).Val() {
					tx.Unwatch(inventory)
					return nil
				}
				pipeliner.ZAdd("market:", &redis.Z{Member: item, Score: float64(price)})
				pipeliner.SRem(inventory, itemID)
				return nil
			}); err != nil {
				return err
			}

			return nil
		}, inventory)

		if err != nil {
			log.Printf("watch err %+v", err)
			return false
		}
		return true
	}
	return false
}

func (a *App) PurchaseItem(buyerID, itemID, sellerID string, lprice float64) bool {
	buyer := fmt.Sprintf("users:%s", buyerID)
	seller := fmt.Sprintf("users:%s", sellerID)
	item := fmt.Sprintf("%s.%s", itemID, sellerID)
	inventory := fmt.Sprintf("inventory:%s", buyerID)
	end := time.Now().Unix() + 10

	for time.Now().Unix() < end {
		if err := a.conn.Watch(func(tx *redis.Tx) error {
			price := tx.ZScore("market:", item).Val()
			funds, _ := tx.HGet(buyer, "funds").Float64()
			if price != lprice || price > funds {
				tx.Unwatch()
			}
			pipe := tx.Pipeline()

			pipe.HIncrBy(seller, "funds", int64(price))
			pipe.HIncrBy(buyer, "funds", int64(-price))
			pipe.SAdd(inventory, itemID)
			pipe.ZRem("market:", item)
			_, err := pipe.Exec()
			return err
		}, "market:", buyer); err != nil {
			log.Printf("%+v", err)
			return false
		}
		return true
	}
	return false

}
