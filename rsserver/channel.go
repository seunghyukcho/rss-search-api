package rsserver

import (
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"net/http"
)

type ChannelController struct {
	Table *DB
}

func (controller *ChannelController) GetChannel(ctx *gin.Context) {
	var err error
	var Id int
	var RSSLink string
	var channels []gofeed.Feed

	rows, err := controller.Table.conn.Query(`SELECT * FROM Channel`)

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

func (controller *ChannelController) GetChannelItems(ctx *gin.Context) {
	var err error
	var Id int
	var channels []gofeed.Feed

	searchWord := "%" + ctx.Param("word") + "%"

	rows, err := controller.Table.conn.Query(`SELECT channel_id, title, description, site_link FROM Channel`)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	} else {
		for rows.Next() {
			var channel gofeed.Feed
			_ = rows.Scan(&Id, &channel.Title, &channel.Description, &channel.Link)
			items, err := controller.Table.conn.Query(`SELECT I.guid, I.title, I.link, I.description, I.pub_date FROM Channel JOIN Publish ON channel_id=channel JOIN Item I ON item=guid WHERE channel_id=? AND I.title LIKE ?`, Id, searchWord)
			if err != nil {
				panic(err)
			}

			for items.Next() {
				var item gofeed.Item
				_ = items.Scan(&item.GUID, &item.Title, &item.Link, &item.Description, &item.Published)

				channel.Items = append(channel.Items, &item)
			}

			if len(channel.Items) > 0 {
				channels = append(channels, channel)
			}
		}

		if len(channels) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "no channel with items"})
		}
		ctx.JSON(http.StatusOK, channels)
	}
}
