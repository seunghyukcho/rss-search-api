package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"net/http"
)

func (env *Env) getChannel(ctx *gin.Context) {
	var err error
	var Id int
	var RSSLink string
	var channels []gofeed.Feed

	rows, err := env.db.Query(`SELECT * FROM Channel`)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	} else {
		var channel gofeed.Feed
		for rows.Next() {
			err = rows.Scan(&Id, &channel.Title, &channel.Description, &channel.Link, &RSSLink)

			if err != nil {
				break
			}

			channels = append(channels, channel)
		}

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, channels)
		}
	}
}

func (env *Env) getChannelItems(ctx *gin.Context) {
	var err error
	var Id int
	var RSSLink string
	var channels []gofeed.Feed

	rows, err := env.db.Query(`SELECT * FROM Channel`)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	} else {
		var channel gofeed.Feed
		for rows.Next() {
			err = rows.Scan(&Id, &channel.Title, &channel.Description, &channel.Link, &RSSLink)

			if err != nil {
				break
			}

			channels = append(channels, channel)
		}

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, channels)
		}
	}
}
