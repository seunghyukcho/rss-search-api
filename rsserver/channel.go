package rsserver

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"net/http"
)

type Channel struct {
	Id      int    `json:"id,omitempty"`
	RSSLink string `json:"rss_link,omitempty"`
	gofeed.Feed
}

type ChannelController struct {
	Table *DB
}

func FetchChannels(channels *sql.Rows, ret *[]Channel) (err error) {
	for channels.Next() {
		var channel Channel

		if err = channels.Scan(&channel.Id, &channel.Title, &channel.Description, &channel.Link, &channel.RSSLink); err != nil {
			return err
		}

		*ret = append(*ret, channel)
	}

	return nil
}

func (controller *ChannelController) GetChannels(ctx *gin.Context) {
	var channels []Channel
	channelRows, err := controller.Table.conn.Query(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	} else {
		if err := FetchChannels(channelRows, &channels); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, channels)
		}
	}
}

func (controller *ChannelController) GetChannelItems(ctx *gin.Context) {
	var err error
	var channels []Channel

	searchWord := "%" + ctx.Param("word") + "%"
	channelRows, err := controller.Table.conn.Query(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	} else {
		if err = FetchChannels(channelRows, &channels); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			for idx, channel := range channels {
				itemRows, err := controller.Table.conn.Query(`SELECT I.guid, I.title, I.link, I.description, I.pub_date FROM Channel JOIN Publish ON channel_id=channel JOIN Item I ON item=guid WHERE channel_id=? AND I.title LIKE ?`, channel.Id, searchWord)

				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				} else if err = FetchItems(itemRows, &channels[idx].Items); err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				} else {
					continue
				}
				return
			}

			ctx.JSON(http.StatusOK, channels)
		}
	}
}
