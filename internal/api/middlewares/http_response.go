package middlewares

import (
	"clean_architecture_go/internal/pkg/error"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Json[TResponse any](c *gin.Context, result TResponse, err *error.RequestError) {
	if err != nil {
		var reqErr *error.RequestError

		if errors.As(err, &reqErr) {
			c.JSON(reqErr.Status, gin.H{
				"message": reqErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
