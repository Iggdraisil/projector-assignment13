package main

import (
	"github.com/beanstalkd/go-beanstalk"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
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
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/produce", func(c *gin.Context) {
		count, err := strconv.Atoi(c.DefaultQuery("count", "1"))
		if err != nil {
			panic(err)
		}
		for i := 0; i < count; i++ {
			switch os.Getenv("CONSUMER_TYPE") {
			case "0":
				redisNoPersist.LPush(c, "job", "foobar")
				break
			case "1":
				redisAOF.LPush(c, "job", "foobar")
				break
			case "2":
				redisRDB.LPush(c, "job", "foobar")
				break
			case "3":
				_, err := clientNoPersist.Put([]byte("hello"), 1, 0, 120*time.Second)
				if err != nil {
					println(err.Error())
				}
				break
			case "4":
				_, err := client.Put([]byte("hello"), 1, 0, 120*time.Second)
				if err != nil {
					println(err.Error())
				}
				break
			}
			consumed <- 1
		}
		c.AbortWithStatus(204)
	})
	go trackConsumption(consumed)
	err = router.Run(":9000")
	if err != nil {
		panic(err)
	}
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
