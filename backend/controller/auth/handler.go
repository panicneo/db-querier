package auth

import (
	"db-querier/utils/token"
	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	token.Create("")
}
