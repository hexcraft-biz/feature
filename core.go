package feature

import (
	"github.com/gin-gonic/gin"
)

// ================================================================
//
// ================================================================
type Feature struct {
	*gin.RouterGroup
}

func New(e *gin.Engine, startsWith string) *Feature {
	return &Feature{
		RouterGroup: e.Group(startsWith),
	}
}

//func (f *Feature) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
//	return f.RouterGroup.GET(relativePath, handlers...)
//}
//
//func (f *Feature) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
//	return f.RouterGroup.POST(relativePath, handlers...)
//}
//
//func (f *Feature) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
//	return f.RouterGroup.PUT(relativePath, handlers...)
//}
//
//func (f *Feature) PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
//	return f.RouterGroup.PATCH(relativePath, handlers...)
//}
//
//func (f *Feature) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
//	return f.RouterGroup.DELETE(relativePath, handlers...)
//}
