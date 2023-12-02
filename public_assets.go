package feature

import "github.com/gin-gonic/gin"

// ================================================================
//
// ================================================================
type PublicAssets struct {
	*Feature
}

func (m *PublicAssets) GET(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipPublic, "GET", relativePath, scopes)
	return m.RouterGroup.GET(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PublicAssets) POST(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipPublic, "POST", relativePath, scopes)
	return m.RouterGroup.POST(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PublicAssets) PUT(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipPublic, "PUT", relativePath, scopes)
	return m.RouterGroup.PUT(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PublicAssets) PATCH(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipPublic, "PATCH", relativePath, scopes)
	return m.RouterGroup.PATCH(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PublicAssets) DELETE(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipPublic, "DELETE", relativePath, scopes)
	return m.RouterGroup.DELETE(relativePath, handlerFuncs(e, handlers)...)
}
