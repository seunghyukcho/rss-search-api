package rssapi

import (
	"github.com/gin-gonic/gin"
	"github.com/shhj1998/rss-search-api/rsserver/rssitem"
	"net/http"
)

func (server *Server) GetItems(ctx *gin.Context) {
	var items []*rssitem.Schema

	if err := server.DB.ItemTable.GetItems(&items); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, items)
	}
}
