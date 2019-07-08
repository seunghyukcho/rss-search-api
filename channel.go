package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed/rss"
	"net/http"
)

func (env *Env) getChannel(ctx *gin.Context) {
	var Id int
	var RSSLink string
	var channel rss.Feed

	rows, err := env.db.Query(`SELECT * FROM Channel`)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	} else {
		for rows.Next() {
			rows.Scan(&Id, &channel.Title, &channel.Description, &channel.Link, &RSSLink)

			fmt.Println(Id, RSSLink)
		}

		ctx.String(http.StatusOK, "channel exists")
	}

}
