package feature

import (
	"github.com/gin-gonic/gin"
)

func New(e *gin.Engine, startsWith string) *gin.RouterGroup {
	return e.Group(startsWith)
}
