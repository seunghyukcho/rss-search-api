package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (env *Env) getChannel(ctx *gin.Context) {
	ctx.String(http.StatusOK, "channel")
}
