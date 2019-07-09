package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbId := os.Getenv("DB_ID")
	dbPassword := os.Getenv("DB_PW")

	dbInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbId, dbPassword, dbHost, dbName)
	conn, err := sql.Open("mysql", dbInfo)

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	dbInstance := &Env{db: conn}

	mainRouter := gin.Default()
	v1 := mainRouter.Group("/api/v1")

	channelRouter := v1.Group("/channel")
	itemRouter := v1.Group("/item")

	{
		channelRouter.GET("", dbInstance.getChannel)
		channelRouter.GET("/items", dbInstance.getChannelItems)
		channelRouter.GET("/items/searchWord/:word", dbInstance.getChannelItems)
	}

	{
		itemRouter.GET("", dbInstance.getItem)
	}

	mainRouter.Run(":8080")
}
