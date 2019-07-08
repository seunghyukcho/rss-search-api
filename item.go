package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (env *Env) getItem(ctx *gin.Context) {
	ctx.String(http.StatusOK, "item!")
}
