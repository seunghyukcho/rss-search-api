package rssapi

import (
	"github.com/gin-gonic/gin"
	"github.com/shhj1998/rss-search-api/rsserver/channel"
	"net/http"
	"strconv"
)

type Channel struct {
	Controller *channel.Controller
}

type createParams struct {
	Link string `form:"rss" json:"rss" xml:"rss" binding:"required"`
}

func (C *Channel) GetChannels(ctx *gin.Context) {
	var channels []channel.Channel
	if err := C.Controller.GetChannels(&channels); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, channels)
	}
}

func (C *Channel) GetChannelsWithItems(ctx *gin.Context) {
	var channels []channel.Channel
	var err error
	var count int

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

	if err := C.Controller.GetChannelsWithItems(&channels, count); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, channels)
	}
}

func (C *Channel) CreateChannel(ctx *gin.Context) {
	var rss createParams
	if err := ctx.ShouldBind(&rss); err == nil {
		if err := C.Controller.CreateChannel(rss.Link); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": "successfully created a new Channel"})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request body doesn't match the api... Please read our api docs"})
	}
}
