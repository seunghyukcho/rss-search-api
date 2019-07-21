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
	"net/http"
	"os"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

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
				if err := db.Update(); err != nil {
					customLogger.Error(err.Error())
				} else {
					customLogger.Info("Update RSS Database successfully!")
				}

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
		channelRouter.GET("", apiServer.SelectChannels)
		channelRouter.GET("/items/", apiServer.SelectChannelsWithItems)
		channelRouter.GET("/items/count/:count", apiServer.SelectChannelsWithItems)
		channelRouter.POST("", apiServer.CreateChannel)
	}

	{
		itemRouter.GET("", apiServer.SelectItems)
	}

	_ = mainRouter.Run(":" + port)
}
