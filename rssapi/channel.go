package rssapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shhj1998/rss-search-api/rsserver/rsschannel"
	"net/http"
	"strconv"
)

type createParams struct {
	Link string `form:"rss" json:"rss" xml:"rss" binding:"required"`
}

func (server *Server) GetChannels(ctx *gin.Context) {
	var channels []*rsschannel.Schema
	if err := server.DB.ChannelTable.Get(&channels); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, channels)
	}
}

func (server *Server) GetChannelsWithItems(ctx *gin.Context) {
	var channels []*rsschannel.Schema
	var err error
	var count int
	var idParams *[]int

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

	fmt.Println(ctx.Request.URL, ctx.QueryArray("id"))
	ids := ctx.QueryArray("id")

	if len(ids) == 0 {
		idParams = nil
	} else {
		params := make([]int, len(ids))
		idParams = &params
		for idx, val := range ids {
			id, err := strconv.Atoi(val)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			} else {
				params[idx] = id
			}
		}
	}

	if err := server.DB.ChannelTable.GetWithItems(&channels, idParams, count); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, channels)
	}
}

func (server *Server) CreateChannel(ctx *gin.Context) {
	var rss createParams
	if err := ctx.ShouldBindJSON(&rss); err == nil {
		if err := server.DB.ChannelTable.Create(rss.Link); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"success": "successfully created a new Channel"})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request body doesn't match the api... Please read our api docs"})
	}
}
