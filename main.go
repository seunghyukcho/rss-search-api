package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/logger"
	"github.com/joho/godotenv"
	"github.com/shhj1998/rss-search-api/rssapi"
	"github.com/shhj1998/rss-search-api/rsserver"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

func update(db *rsserver.DB, logger *logger.Logger) {
	if err := db.Update(); err != nil {
		logger.Error(err.Error())
	} else {
		logger.Info("Update RSS Database successfully!")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.SetFlags(log.LstdFlags)
	customLogger := logger.Init("LogUpdate", true, false, ioutil.Discard)
	if err := godotenv.Load(); err != nil {
		fmt.Println(err.Error())
	}

	logger.Info("Total CPUs : ", runtime.NumCPU())
	port, name, address, id, password := os.Getenv("PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_ID"), os.Getenv("DB_PW")
	rssDB := rsserver.DB{}

	if err := rssDB.Open(name, address, id, password); err != nil {
		panic(err)
	}

	defer rssDB.Close()

	if err := rssDB.Create(); err != nil {
		panic(err)
	}

	go func(db *rsserver.DB) {
		ticker := time.NewTicker(25 * time.Minute)

		for {
			select {
			case <-ticker.C:
				update(db, customLogger)
				_, _ = http.Get("https://rss-search-api.herokuapp.com/api/v1/channel")
			}
		}
	}(&rssDB)

	apiServer := rssapi.Server{DB: &rssDB}

	mainRouter := gin.Default()
	v1 := mainRouter.Group("/api/v1")

	channelRouter := v1.Group("/channel")
	itemRouter := v1.Group("/item")

	{
		channelRouter.GET("", apiServer.GetChannels)
		channelRouter.GET("/items/", apiServer.GetChannelsWithItems)
		channelRouter.GET("/items/count/:count", apiServer.GetChannelsWithItems)
		channelRouter.POST("", apiServer.CreateChannel)
	}

	{
		itemRouter.GET("", apiServer.GetItems)
	}

	mainRouter.Run(":" + port)
}
