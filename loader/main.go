package main

import (
	//  "bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	} else {
		fmt.Println("successful...")
	}
}

type User struct {
	State string `json:"state" `
	Count uint   `json:"count"`
}

type CacheEntry struct {
	Value []User `json:"value" `
	Delta int64  `json:"delta"`
}

const cacheDuration = 100 * time.Second

func main() {
	gin.SetMode(gin.ReleaseMode)
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName: "mymaster",
		SentinelAddrs: []string{
			"redis-sentinel1:5000",
			"redis-sentinel2:5000",
			"redis-sentinel3:5000",
		},
	})
	db, err := sql.Open("postgres", "host=localhost user=postgres password=1234 sslmode=disable database=mydb")
	db.SetMaxOpenConns(100000)
	checkErr(err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	err = db.Ping()
	checkErr(err)

	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/prob-cache", func(c *gin.Context) {
		data := rdb.Get(c, "comp_value")
		if data.Val() == "" {
			fetchDb(c, db, rdb)
		} else {
			var result CacheEntry
			err := json.Unmarshal([]byte(data.Val()), &result)
			if err != nil {
				panic(err)
			}
			ttl := rdb.TTL(c, "comp_value").Val()
			if int64(float64(result.Delta)*rand.ExpFloat64()) >= int64(ttl/time.Second) {
				fetchDb(c, db, rdb)
			} else {
				c.JSON(200, gin.H{
					"data": result.Value,
				})
			}
		}
	})
	router.GET("/cache", func(c *gin.Context) {
		data := rdb.Get(c, "comp_value")
		if data.Val() == "" {
			fetchDb(c, db, rdb)
		} else {
			var result CacheEntry
			err := json.Unmarshal([]byte(data.Val()), &result)
			if err != nil {
				panic(err)
			}
			c.JSON(200, gin.H{
				"data": result.Value,
			})
		}
	})
	err = router.Run(":9000")
	if err != nil {
		panic(err)
	}
}

func fetchDb(c *gin.Context, db *sql.DB, rdb *redis.Client) {
	println("fetch db")
	rows, delta := fetchFromDB(db)
	c.JSON(200, gin.H{
		"data": rows,
	})
	data, _ := json.Marshal(CacheEntry{Value: rows, Delta: delta})
	rdb.SetEx(c, "comp_value", interface{}(data), cacheDuration)
}

func fetchFromDB(db *sql.DB) ([]User, int64) {
	start := time.Now().Unix()
	rows, err := db.Query("select state, count(state) from \"user\" where birthdate between '2000-05-26' and '2004-06-26' group by state")
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)
	if err != nil {
		panic(err)
	}
	var data []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.State, &user.Count)
		if err != nil {
			panic(err)
			return nil, 0
		}
		data = append(data, user)
	}
	delta := time.Now().Unix() - start
	return data, delta * 2
}
