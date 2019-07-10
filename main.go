package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/shhj1998/rss-search-api/rsserver"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
	}

	port, name, address, id, password := os.Getenv("PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_ID"), os.Getenv("DB_PW")
	rssDB := rsserver.DB{}
	err = rssDB.Open(name, address, id, password)
	if err != nil {
		panic(err)
	}
	defer rssDB.Close()

	itemInstance := rsserver.ItemController{Table: &rssDB}
	channelInstance := rsserver.ChannelController{Table: &rssDB}

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
