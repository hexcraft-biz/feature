package feature

import (
	"github.com/gin-gonic/gin"
	dogmas "github.com/hexcraft-biz/envmod-dogmas"
)

type Feature struct {
	*gin.RouterGroup
	Scopes []*dogmas.Scope
}

func New(e *gin.Engine, startsWith string) *Feature {
	return &Feature{
		RouterGroup: e.Group(startsWith),
		Scopes:      []*dogmas.Scope{},
	}
}

func (f *Feature) GET(relativePath string, scope *dogmas.Scope, handlers ...gin.HandlerFunc) gin.IRoutes {
	f.Scopes = append(f.Scopes, scope)
	return f.RouterGroup.GET(relativePath, handlers...)
}

func (f *Feature) POST(relativePath string, scope *dogmas.Scope, handlers ...gin.HandlerFunc) gin.IRoutes {
	f.Scopes = append(f.Scopes, scope)
	return f.RouterGroup.POST(relativePath, handlers...)
}

func (f *Feature) PUT(relativePath string, scope *dogmas.Scope, handlers ...gin.HandlerFunc) gin.IRoutes {
	f.Scopes = append(f.Scopes, scope)
	return f.RouterGroup.PUT(relativePath, handlers...)
}

func (f *Feature) PATCH(relativePath string, scope *dogmas.Scope, handlers ...gin.HandlerFunc) gin.IRoutes {
	f.Scopes = append(f.Scopes, scope)
	return f.RouterGroup.PATCH(relativePath, handlers...)
}

func (f *Feature) DELETE(relativePath string, scope *dogmas.Scope, handlers ...gin.HandlerFunc) gin.IRoutes {
	f.Scopes = append(f.Scopes, scope)
	return f.RouterGroup.DELETE(relativePath, handlers...)
}
