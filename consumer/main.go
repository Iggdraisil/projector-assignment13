package main

import (
	"context"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

func main() {
	c := context.Background()
	redisRDB := redis.NewClient(&redis.Options{
		Addr: "redis-rdb:6379",
	})
	redisAOF := redis.NewClient(&redis.Options{
		Addr: "redis-aof:6379",
	})

	redisNoPersist := redis.NewClient(&redis.Options{
		Addr: "redis-nopersist:6379",
	})

	client, err := beanstalk.Dial("tcp", "beanstalkd-persist:11300")
	clientNoPersist, err := beanstalk.Dial("tcp", "beanstalkd-nopersist:11300")
	consumed := make(chan int)
	if err != nil {
		panic(err)
	}
	consumers, err := strconv.Atoi(os.Getenv("NUM_CONSUMERS"))
	if err != nil || consumers <= 0 {
		consumers = 1
	}
	for i := 0; i < consumers; i++ {
		switch os.Getenv("CONSUMER_TYPE") {
		case "0":
			go consumeLPop(redisNoPersist, c, consumed)
			break
		case "1":
			go consumeLPop(redisAOF, c, consumed)
			break
		case "2":
			go consumeLPop(redisRDB, c, consumed)
			break
		case "3":
			go consumeBeanstalk(clientNoPersist, consumed)
			break
		case "4":
			go consumeBeanstalk(client, consumed)
			break
		}
	}
	go trackConsumption(consumed)
	<-make(chan int)

}

func trackConsumption(consumed chan int) {
	total := int64(0)
	for {
		for i := 0; i < 1000; i += <-consumed {
		}
		total += int64(1000)
		println(total)
	}
}

func consumeLPop(client *redis.Client, ctx context.Context, consumed chan int) {
	for {
		_, err := client.LPop(ctx, "job").Result()
		if err != nil {
			if err.Error() != "redis: nil" {
				println(err.Error())
			}
		} else {
			consumed <- 1
		}

	}
}

func consumeBeanstalk(client *beanstalk.Conn, consumed chan int) {
	for {
		_, _, err := client.Reserve(5 * time.Second)
		if err != nil {
			println(err.Error())
		} else {
			consumed <- 1
		}
	}
}
