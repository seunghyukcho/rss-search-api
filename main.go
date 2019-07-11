package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/shhj1998/rss-search-api/rsserver"
	"github.com/shhj1998/rss-search-api/rsserver/channel"
	"github.com/shhj1998/rss-search-api/rsserver/item"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err.Error())
	}

	port, name, address, id, password := os.Getenv("PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_ID"), os.Getenv("DB_PW")
	rssDB := rsserver.DB{}
	if err := rssDB.Open(name, address, id, password); err != nil {
		panic(err)
	}

	if err := rssDB.Update(); err != nil {
		panic(err)
	}

	defer rssDB.Close()

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
	}

	{
		itemRouter.GET("", itemInstance.GetItems)
	}

	mainRouter.Run(":" + port)
}
