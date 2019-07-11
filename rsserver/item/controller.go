package item

import (
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver"
	"net/http"
)

type Controller struct {
	Table *rsserver.DB
}

func (controller *Controller) GetItems(ctx *gin.Context) {
	var err error
	var items []*gofeed.Item

	itemRows, err := controller.Table.Connection.Query(`SELECT guid, title, link, description, pub_date, creator, url, length, type FROM Item LEFT JOIN Enclosure ON guid=item`)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err := Fetch(itemRows, &items); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, items)
		}
	}
}
