package channel

import (
	"github.com/gin-gonic/gin"
	"github.com/shhj1998/rss-search-api/rsserver"
	"github.com/shhj1998/rss-search-api/rsserver/item"
	"net/http"
)

type Controller struct {
	Table *rsserver.DB
}

func (controller *Controller) GetChannels(ctx *gin.Context) {
	var channels []Channel
	channelRows, err := controller.Table.Connection.Query(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	} else {
		if err := Fetch(channelRows, &channels); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, channels)
		}
	}
}

func (controller *Controller) GetChannelItems(ctx *gin.Context) {
	var err error
	var channels []Channel

	searchWord := "%" + ctx.Param("word") + "%"
	channelRows, err := controller.Table.Connection.Query(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err = Fetch(channelRows, &channels); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			for idx, channel := range channels {
				itemRows, err := controller.Table.Connection.Query(`SELECT I.guid, I.title, I.link, I.description, I.pub_date FROM Channel JOIN Publish ON channel_id=channel JOIN Item I ON item=guid WHERE channel_id=? AND I.title LIKE ?`, channel.Id, searchWord)

				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				} else if err = item.Fetch(itemRows, &channels[idx].Items); err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				} else {
					continue
				}
				return
			}

			ctx.JSON(http.StatusOK, channels)
		}
	}
}
