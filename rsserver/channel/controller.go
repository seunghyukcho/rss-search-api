package channel

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver"
	"github.com/shhj1998/rss-search-api/rsserver/item"
	"net/http"
	"strconv"
)

type Controller struct {
	Table *rsserver.DB
}

type createParams struct {
	Link string `form:"rss" json:"rss" xml:"rss" binding:"required"`
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
	var count int

	searchWord := "%" + ctx.Param("word") + "%"
	countParam := ctx.Param("count")
	if countParam == "" {
		count = 20
	} else {
		count, err = strconv.Atoi(ctx.Param("count"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	channelRows, err := controller.Table.Connection.Query(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err = Fetch(channelRows, &channels); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			for idx, channel := range channels {
				itemRows, err := controller.Table.Connection.Query(`SELECT I.guid, I.title, I.link, I.description, I.pub_date, I.creator, e.url, e.length, e.type 
																			FROM Channel c JOIN Publish p ON c.channel_id=p.channel JOIN Item I ON p.item=I.guid LEFT JOIN Enclosure e ON I.guid=e.item 
																			WHERE channel_id=? AND I.title LIKE ? ORDER BY I.pub_date DESC LIMIT 0,?`, channel.Id, searchWord, count)

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

func (controller *Controller) CreateChannel(ctx *gin.Context) {
	var rss createParams
	if err := ctx.ShouldBind(&rss); err == nil {
		fmt.Println(rss.Link)
		fp := gofeed.NewParser()
		if _, err := fp.ParseURL(rss.Link); err == nil {
			if _, err = controller.Table.Connection.Exec(`INSERT INTO Channel(rss_link) VALUE(?)`, rss.Link); err == nil {
				ctx.JSON(http.StatusOK, gin.H{"message": "created a new channel!"})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid rss link..."})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request body doesn't match the api... Please read our api docs"})
	}
}
