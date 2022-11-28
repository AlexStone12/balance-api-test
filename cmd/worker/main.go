package main

import (
	"bank-api-test/internal/database"
	"bank-api-test/internal/worker"
	"github.com/adjust/rmq/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env файл не найден")
	}
}

/* for {
Main Worker получает список очередей, добавляет горутину в waitgroup, запоминает id этого воркера в map
		routine{
			берет пэйлоад, распарсивает json в структуру, обращается в базу по данным структуры и меняет значения, +
			добавляет сделанную операцию в историю транзакций

		}
	убираем очередь из мапы?
	НЕ?спим 5 секунд
}
*/

func main() {
	redis := database.RedisConnect()
	db := database.PostgresConnect()
	var workers sync.Map

	for {
		queues, err := redis.GetOpenQueues()
		if err != nil {
			return
		}
		stats, err := rmq.CollectStats(queues, redis)
		if err != nil {
			return
		}
		for _, q := range queues {
			if _, ok := workers.Load(q); !ok {
				ready := stats.QueueStats[q].ReadyCount
				if ready > 0 {
					workers.Store(q, struct{}{})
					go Worker(db, redis, &workers, q)
				}
			}
		}
	}
}

func Worker(db *sqlx.DB, redis rmq.Connection, workers *sync.Map, name string) {
	log.Printf("Worker started")
	queue, err := redis.OpenQueue(name)
	if err != nil {
		log.Printf("%e", err)
		return
	}
	err = queue.StartConsuming(10, time.Second)
	if err != nil {
		log.Printf("%e", err)
		return
	}
	if _, err := queue.AddConsumer(name, worker.NewConsumer(name, db)); err != nil {
		log.Printf("%e", err)
		return
	}

	stats, err := rmq.CollectStats([]string{name}, redis)
	if err != nil {
		return
	}

	ready := stats.QueueStats[name].ReadyCount
	for ready > 0 {
		stats, err = rmq.CollectStats([]string{name}, redis)
		if err != nil {
			return
		}
		ready = stats.QueueStats[name].ReadyCount
		time.Sleep(time.Second)
	}
	workers.Delete(name)
	log.Printf("Worker stoped: %s", name)
	queue.StopConsuming()
}
