package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed/rss"
	"net/http"
)

func (env *Env) getItem(ctx *gin.Context) {
	var Id int
	var RSSLink string
	var item rss.Item

	rows, err := env.db.Query(`SELECT * FROM Item`)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	} else {
		for rows.Next() {
			rows.Scan(&item.GUID.Value, &item.Title, &item.Link, &item.Enclosure.URL, &item.Description, &item.PubDate)

			fmt.Println(Id, RSSLink)
		}

		ctx.String(http.StatusOK, "item exists")
	}
}
