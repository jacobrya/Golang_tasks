package v1

import (
	"assignment7/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, u usecase.UserInterface) {

	v1 := handler.Group("/v1")
	{
		newUserRoutes(v1, u)
	}
}
