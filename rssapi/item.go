package rssapi

import (
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver/item"
	"net/http"
)

type Item struct {
	Controller *item.Controller
}

func (I *Item) GetItems(ctx *gin.Context) {
	var items []*gofeed.Item

	if err := I.Controller.GetItems(&items); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, items)
	}
}
