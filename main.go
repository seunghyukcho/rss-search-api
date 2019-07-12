package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/logger"
	"github.com/joho/godotenv"
	"github.com/shhj1998/rss-search-api/rsserver"
	"github.com/shhj1998/rss-search-api/rsserver/channel"
	"github.com/shhj1998/rss-search-api/rsserver/item"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	logger.SetFlags(log.LstdFlags)
	customLogger := logger.Init("LogUpdate", true, false, ioutil.Discard)
	if err := godotenv.Load(); err != nil {
		fmt.Println(err.Error())
	}

	port, name, address, id, password := os.Getenv("PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_ID"), os.Getenv("DB_PW")
	rssDB := rsserver.DB{}
	if err := rssDB.Open(name, address, id, password); err != nil {
		panic(err)
	}

	defer rssDB.Close()

	go func(db *rsserver.DB) {
		ticker := time.NewTicker(1 * time.Hour)

		for {
			select {
			case <-ticker.C:
				if err := db.Update(); err != nil {
					customLogger.Error(err.Error())
				} else {
					customLogger.Info("Update RSS Database successfully!")
				}
			}
		}
	}(&rssDB)

	itemInstance := item.Controller{Table: &rssDB}
	channelInstance := channel.Controller{Table: &rssDB}

	mainRouter := gin.Default()
	v1 := mainRouter.Group("/api/v1")

	channelRouter := v1.Group("/channel")
	itemRouter := v1.Group("/item")

	{
		channelRouter.GET("", channelInstance.GetChannels)
		channelRouter.GET("/items", channelInstance.GetChannelItems)
		channelRouter.GET("/items/searchWord/:word", channelInstance.GetChannelItems)
		channelRouter.POST("", channelInstance.CreateChannel)
	}

	{
		itemRouter.GET("", itemInstance.GetItems)
	}

	mainRouter.Run(":" + port)
}
