package helpers

import "github.com/gin-gonic/gin"

func SendResponse(ctx *gin.Context, statusCode int, message string, body any) {
	ctx.JSON(statusCode, gin.H{
		"message": message,
		"body":    body,
	})
}
