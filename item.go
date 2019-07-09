package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"net/http"
)

func (env *Env) getItem(ctx *gin.Context) {
	var err error
	var items []gofeed.Item

	rows, err := env.db.Query(`SELECT guid, title, link, description, pub_date FROM Item`)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	} else {
		for rows.Next() {
			var item gofeed.Item

			err = rows.Scan(&item.GUID, &item.Title, &item.Link, &item.Description, &item.Published)
			if err != nil {
				break
			}

			items = append(items, item)
		}

		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
		} else {
			ctx.JSON(http.StatusOK, items)
		}
	}
}
