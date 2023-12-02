package feature

import "github.com/gin-gonic/gin"

// ================================================================
//
// ================================================================
type PrivateAssets struct {
	*Feature
}

func (m *PrivateAssets) GET(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipPrivate, "GET", relativePath, scopes)
	return m.RouterGroup.GET(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PrivateAssets) POST(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipPrivate, "POST", relativePath, scopes)
	return m.RouterGroup.POST(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PrivateAssets) PUT(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipPrivate, "PUT", relativePath, scopes)
	return m.RouterGroup.PUT(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PrivateAssets) PATCH(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipPrivate, "PATCH", relativePath, scopes)
	return m.RouterGroup.PATCH(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PrivateAssets) DELETE(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipPrivate, "DELETE", relativePath, scopes)
	return m.RouterGroup.DELETE(relativePath, handlerFuncs(e, handlers)...)
}
